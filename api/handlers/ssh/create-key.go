package ssh

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

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

	privPEM, pubAuth, err := generateRSA(bits, req.Comment)
	if err != nil {
		http.Error(w, "key generation failed", http.StatusInternalServerError)
		return
	}

	key := models.SshKey{
		OrganizationID: ac.OrganizationID,
		PublicKey:      pubAuth,
		PrivateKey:     privPEM,
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
	default:
		// fall through to JSON
	}

	writeJSON(w, http.StatusCreated, sshResponse{
		ID:             key.ID,
		OrganizationID: key.OrganizationID,
		PublicKey:      key.PublicKey,
		CreatedAt:      key.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      key.UpdatedAt.UTC().Format(time.RFC3339),
	})
}
