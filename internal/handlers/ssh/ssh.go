package ssh

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// ListPublicKeys godoc
// @Summary      List ssh keys (org scoped)
// @Description  Returns ssh keys for the organization in X-Org-ID.
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Security     BearerAuth
// @Success      200 {array}  sshResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list keys"
// @Router       /api/v1/ssh [get]
func ListPublicKeys(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var rows []models.SshKey
	if err := db.DB.Where("organization_id = ?", ac.OrganizationID).
		Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list ssh keys", http.StatusInternalServerError)
		return
	}

	out := make([]sshResponse, 0, len(rows))
	for _, s := range rows {
		out = append(out, sshResponse{
			ID:             s.ID,
			OrganizationID: s.OrganizationID,
			Name:           s.Name,
			PublicKey:      s.PublicKey,
			Fingerprint:    s.Fingerprint,
			CreatedAt:      s.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      s.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	_ = response.JSON(w, http.StatusOK, out)
}

// CreateSSHKey godoc
// @Summary      Create ssh keypair (org scoped)
// @Description  Generates an RSA keypair, saves it, and returns metadata. Optionally set `download` to "public", "private", or "both" to download files immediately.
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createSSHRequest true "Key generation options"
// @Security     BearerAuth
// @Success      201 {object} sshResponse
// @Header       201 {string} Content-Disposition "When download is requested"
// @Failure      400 {string} string "invalid json / invalid bits"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "generation/create failed"
// @Router       /api/v1/ssh [post]
func CreateSSHKey(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createSSHRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	bits := 4096
	if req.Bits != nil {
		if !allowedBits(*req.Bits) {
			http.Error(w, "invalid bits (allowed: 2048, 3072, 4096)", http.StatusBadRequest)
			return
		}
		bits = *req.Bits
	}

	privPEM, pubAuth, err := GenerateRSAPEMAndAuthorized(bits, strings.TrimSpace(req.Comment))
	if err != nil {
		http.Error(w, "key generation failed", http.StatusInternalServerError)
		return
	}

	cipher, iv, tag, err := utils.EncryptForOrg(ac.OrganizationID, []byte(privPEM))
	if err != nil {
		http.Error(w, "encryption failed", http.StatusInternalServerError)
		return
	}

	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubAuth))
	if err != nil {
		http.Error(w, "failed to parse public key", http.StatusInternalServerError)
		return
	}
	fp := ssh.FingerprintSHA256(parsed)

	key := models.SshKey{
		OrganizationID:      ac.OrganizationID,
		Name:                req.Name,
		PublicKey:           pubAuth,
		EncryptedPrivateKey: cipher,
		PrivateIV:           iv,
		PrivateTag:          tag,
		Fingerprint:         fp,
	}
	if err := db.DB.Create(&key).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	// Immediate download if requested
	switch strings.ToLower(strings.TrimSpace(req.Download)) {
	case "public":
		filename := fmt.Sprintf("id_rsa_%s.pub", key.ID.String())
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(pubAuth))
		return
	case "private":
		filename := fmt.Sprintf("id_rsa_%s.pem", key.ID.String())
		w.Header().Set("Content-Type", "application/x-pem-file")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(privPEM))
		return
	case "both":
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		_ = toZipFile(fmt.Sprintf("id_rsa_%s.pem", key.ID.String()), []byte(privPEM), zw)
		_ = toZipFile(fmt.Sprintf("id_rsa_%s.pub", key.ID.String()), []byte(pubAuth), zw)
		_ = zw.Close()

		filename := fmt.Sprintf("ssh_key_%s.zip", key.ID.String())
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(buf.Bytes())
		return
	}

	_ = response.JSON(w, http.StatusCreated, sshResponse{
		ID:             key.ID,
		OrganizationID: key.OrganizationID,
		PublicKey:      key.PublicKey,
		Fingerprint:    key.Fingerprint,
		CreatedAt:      key.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      key.UpdatedAt.UTC().Format(time.RFC3339),
	})
}

// GetSSHKey godoc
// @Summary      Get ssh key by ID (org scoped)
// @Description  Returns public key fields. Append `?reveal=true` to include the private key PEM.
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "SSH Key ID (UUID)"
// @Param        reveal query bool false "Reveal private key PEM"
// @Security     BearerAuth
// @Success      200 {object} sshResponse
// @Success      200 {object} sshRevealResponse "When reveal=true"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/ssh/{id} [get]
func GetSSHKey(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var key models.SshKey
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("reveal") != "true" {
		_ = response.JSON(w, http.StatusOK, sshResponse{
			ID:             key.ID,
			OrganizationID: key.OrganizationID,
			PublicKey:      key.PublicKey,
			CreatedAt:      key.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      key.UpdatedAt.UTC().Format(time.RFC3339),
		})
		return
	}

	plain, err := utils.DecryptForOrg(ac.OrganizationID, key.EncryptedPrivateKey, key.PrivateIV, key.PrivateTag)
	if err != nil {
		http.Error(w, "failed to decrypt", http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusOK, sshRevealResponse{
		sshResponse: sshResponse{
			ID:             key.ID,
			OrganizationID: key.OrganizationID,
			PublicKey:      key.PublicKey,
			CreatedAt:      key.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      key.UpdatedAt.UTC().Format(time.RFC3339),
		},
		PrivateKey: plain,
	})
}

// DeleteSSHKey godoc
// @Summary      Delete ssh keypair (org scoped)
// @Description  Permanently deletes a keypair.
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "SSH Key ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/ssh/{id} [delete]
func DeleteSSHKey(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/ssh/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Delete(&models.SshKey{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	response.NoContent(w)
}

// DownloadSSHKey godoc
// @Summary      Download ssh key files by ID (org scoped)
// @Description  Download `part=public|private|both` of the keypair. `both` returns a zip file.
// @Tags         ssh
// @Produce      text/plain
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "SSH Key ID (UUID)"
// @Param        part query string true "Which part to download" Enums(public,private,both)
// @Security     BearerAuth
// @Success      200 {string} string "file content"
// @Failure      400 {string} string "invalid id / invalid part"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "download failed"
// @Router       /api/v1/ssh/{id}/download [get]
func DownloadSSHKey(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var key models.SshKey
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	switch strings.ToLower(r.URL.Query().Get("part")) {
	case "public":
		filename := fmt.Sprintf("id_rsa_%s.pub", key.ID.String())
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		_, _ = w.Write([]byte(key.PublicKey))
	case "private":
		plain, err := utils.DecryptForOrg(ac.OrganizationID, key.EncryptedPrivateKey, key.PrivateIV, key.PrivateTag)
		if err != nil {
			http.Error(w, "decrypt failed", http.StatusInternalServerError)
			return
		}
		filename := fmt.Sprintf("id_rsa_%s.pem", key.ID.String())
		w.Header().Set("Content-Type", "application/x-pem-file")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		_, _ = w.Write([]byte(plain))
	case "both":
		plain, err := utils.DecryptForOrg(ac.OrganizationID, key.EncryptedPrivateKey, key.PrivateIV, key.PrivateTag)
		if err != nil {
			http.Error(w, "decrypt failed", http.StatusInternalServerError)
			return
		}
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		_ = toZipFile(fmt.Sprintf("id_rsa_%s.pem", key.ID.String()), []byte(plain), zw)
		_ = toZipFile(fmt.Sprintf("id_rsa_%s.pub", key.ID.String()), []byte(key.PublicKey), zw)
		_ = zw.Close()

		filename := fmt.Sprintf("ssh_key_%s.zip", key.ID.String())
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		_, _ = w.Write(buf.Bytes())
	default:
		http.Error(w, "invalid part (public|private|both)", http.StatusBadRequest)
	}
}
