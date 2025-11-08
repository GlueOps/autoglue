package handlers

import (
	"archive/zip"
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// ListPublicSshKeys godoc
//
//	@ID				ListPublicSshKeys
//	@Summary		List ssh keys (org scoped)
//	@Description	Returns ssh keys for the organization in X-Org-ID.
//	@Tags			Ssh
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Success		200			{array}		dto.SshResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"failed to list keys"
//	@Router			/ssh [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListPublicSshKeys(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var out []dto.SshResponse
		if err := db.
			Model(&models.SshKey{}).
			Where("organization_id = ?", orgID).
			// avoid selecting encrypted columns here
			Select("id", "organization_id", "name", "public_key", "fingerprint", "created_at", "updated_at").
			Order("created_at DESC").
			Scan(&out).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to list ssh keys")
			return
		}

		if out == nil {
			out = []dto.SshResponse{}
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// CreateSSHKey
//
//	@ID				CreateSSHKey
//	@Summary		Create ssh keypair (org scoped)
//	@Description	Generates an RSA or ED25519 keypair, saves it, and returns metadata. For RSA you may set bits (2048/3072/4096). Default is 4096. ED25519 ignores bits.
//	@Tags			Ssh
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string					false	"Organization UUID"
//	@Param			body		body		dto.CreateSSHRequest	true	"Key generation options"
//	@Success		201			{object}	dto.SshResponse
//	@Failure		400			{string}	string	"invalid json / invalid bits"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"generation/create failed"
//	@Router			/ssh [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateSSHKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var req dto.CreateSSHRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid_payload", "invalid JSON payload")
			return
		}

		keyType := "rsa"
		if req.Type != nil && strings.TrimSpace(*req.Type) != "" {
			keyType = strings.ToLower(strings.TrimSpace(*req.Type))
		}

		if keyType != "rsa" && keyType != "ed25519" {
			utils.WriteError(w, http.StatusBadRequest, "invalid_type", "invalid type (rsa|ed25519)")
			return
		}

		var (
			privPEM string
			pubAuth string
			err     error
		)

		switch keyType {
		case "rsa":
			bits := 4096
			if req.Bits != nil {
				if !allowedBits(*req.Bits) {
					utils.WriteError(w, http.StatusBadRequest, "invalid_bits", "invalid bits (allowed: 2048, 3072, 4096)")
					return
				}
				bits = *req.Bits
			}
			privPEM, pubAuth, err = GenerateRSAPEMAndAuthorized(bits, strings.TrimSpace(req.Comment))

		case "ed25519":
			if req.Bits != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid_bits_for_type", "bits is only valid for RSA")
				return
			}
			privPEM, pubAuth, err = GenerateEd25519PEMAndAuthorized(strings.TrimSpace(req.Comment))
		}

		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "keygen_failure", "key generation failed")
			return
		}

		cipher, iv, tag, err := utils.EncryptForOrg(orgID, []byte(privPEM), db)
		if err != nil {
			http.Error(w, "encryption failed", http.StatusInternalServerError)
			return
		}

		parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubAuth))
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "ssh_failure", "ssh public key parsing failed")
			return
		}

		fp := ssh.FingerprintSHA256(parsed)

		key := models.SshKey{
			AuditFields: common.AuditFields{
				OrganizationID: orgID,
			},
			Name:                req.Name,
			PublicKey:           pubAuth,
			EncryptedPrivateKey: cipher,
			PrivateIV:           iv,
			PrivateTag:          tag,
			Fingerprint:         fp,
		}

		if err := db.Create(&key).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to create ssh key")
			return
		}

		utils.WriteJSON(w, http.StatusCreated, dto.SshResponse{
			AuditFields: key.AuditFields,
			Name:        key.Name,
			PublicKey:   key.PublicKey,
			Fingerprint: key.Fingerprint,
		})
	}
}

// GetSSHKey godoc
//
//	@ID				GetSSHKey
//	@Summary		Get ssh key by ID (org scoped)
//	@Description	Returns public key fields. Append `?reveal=true` to include the private key PEM.
//	@Tags			Ssh
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"SSH Key ID (UUID)"
//	@Param			reveal		query		bool	false	"Reveal private key PEM"
//	@Success		200			{object}	dto.SshResponse
//	@Success		200			{object}	dto.SshRevealResponse	"When reveal=true"
//	@Failure		400			{string}	string					"invalid id"
//	@Failure		401			{string}	string					"Unauthorized"
//	@Failure		403			{string}	string					"organization required"
//	@Failure		404			{string}	string					"not found"
//	@Failure		500			{string}	string					"fetch failed"
//	@Router			/ssh/{id} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetSSHKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid_ssh_key_id", "invalid SSH Key ID")
			return
		}

		reveal := strings.EqualFold(r.URL.Query().Get("reveal"), "true")

		if !reveal {
			var out dto.SshResponse
			if err := db.
				Model(&models.SshKey{}).
				Where("id = ? AND organization_id = ?", id, orgID).
				Select("id", "organization_id", "name", "public_key", "fingerprint", "created_at", "updated_at").
				Limit(1).
				Scan(&out).Error; err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to get ssh key")
				return
			}
			if out.ID == uuid.Nil {
				utils.WriteError(w, http.StatusNotFound, "ssh_key_not_found", "ssh key not found")
				return
			}
			utils.WriteJSON(w, http.StatusOK, out)
			return
		}

		var secret dto.SshResponse
		if err := db.
			Model(&models.SshKey{}).
			Where("id = ? AND organization_id = ?", id, orgID).
			// include the encrypted bits too
			Select("id", "organization_id", "name", "public_key", "fingerprint",
				"encrypted_private_key", "private_iv", "private_tag",
				"created_at", "updated_at").
			Limit(1).
			Scan(&secret).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to get ssh key")
			return
		}

		if secret.ID == uuid.Nil {
			utils.WriteError(w, http.StatusNotFound, "ssh_key_not_found", "ssh key not found")
			return
		}

		plain, err := utils.DecryptForOrg(orgID, secret.EncryptedPrivateKey, secret.PrivateIV, secret.PrivateTag, db)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to decrypt ssh key")
			return
		}

		utils.WriteJSON(w, http.StatusOK, dto.SshRevealResponse{
			SshResponse: dto.SshResponse{
				AuditFields: secret.AuditFields,
				Name:        secret.Name,
				PublicKey:   secret.PublicKey,
				Fingerprint: secret.Fingerprint,
			},
			PrivateKey: plain,
		})
	}
}

// DeleteSSHKey godoc
//
//	@ID				DeleteSSHKey
//	@Summary		Delete ssh keypair (org scoped)
//	@Description	Permanently deletes a keypair.
//	@Tags			Ssh
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"SSH Key ID (UUID)"
//	@Success		204			{string}	string	"No Content"
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"delete failed"
//	@Router			/ssh/{id} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DeleteSSHKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid_ssh_key_id", "invalid SSH Key ID")
			return
		}

		res := db.Where("id = ? AND organization_id = ?", id, orgID).
			Delete(&models.SshKey{})
		if res.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to delete ssh key")
			return
		}
		if res.RowsAffected == 0 {
			utils.WriteError(w, http.StatusNotFound, "ssh_key_not_found", "ssh key not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// DownloadSSHKey godoc
//
//	@ID				DownloadSSHKey
//	@Summary		Download ssh key files by ID (org scoped)
//	@Description	Download `part=public|private|both` of the keypair. `both` returns a zip file.
//	@Tags			Ssh
//	@Produce		json
//	@Param			X-Org-ID	header		string	true	"Organization UUID"
//	@Param			id			path		string	true	"SSH Key ID (UUID)"
//	@Param			part		query		string	true	"Which part to download"	Enums(public,private,both)
//	@Success		200			{string}	string	"file content"
//	@Failure		400			{string}	string	"invalid id / invalid part"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"download failed"
//	@Router			/ssh/{id}/download [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DownloadSSHKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid_ssh_key_id", "invalid SSH Key ID")
			return
		}

		var key models.SshKey
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).
			First(&key).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "ssh_key_not_found", "ssh key not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to get ssh key")
			return
		}

		part := strings.ToLower(r.URL.Query().Get("part"))
		if part == "" {
			utils.WriteError(w, http.StatusBadRequest, "invalid_ssh_part", "invalid part (public|private|both)")
			return
		}

		mode := strings.ToLower(r.URL.Query().Get("mode"))
		if mode != "" && mode != "json" {
			utils.WriteError(w, http.StatusBadRequest, "invalid_mode", "invalid mode (json|attachment[default])")
			return
		}

		if mode == "json" {
			resp := dto.SshMaterialJSON{
				ID:          key.ID.String(),
				Name:        key.Name,
				Fingerprint: key.Fingerprint,
			}
			switch part {
			case "public":
				pub := key.PublicKey
				resp.PublicKey = &pub
				resp.Filenames = []string{fmt.Sprintf("%s.pub", key.ID.String())}
				utils.WriteJSON(w, http.StatusOK, resp)
				return

			case "private":
				plain, err := utils.DecryptForOrg(orgID, key.EncryptedPrivateKey, key.PrivateIV, key.PrivateTag, db)
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to decrypt ssh key")
					return
				}
				resp.PrivatePEM = &plain
				resp.Filenames = []string{fmt.Sprintf("%s.pem", key.ID.String())}
				utils.WriteJSON(w, http.StatusOK, resp)
				return

			case "both":
				plain, err := utils.DecryptForOrg(orgID, key.EncryptedPrivateKey, key.PrivateIV, key.PrivateTag, db)
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to decrypt ssh key")
					return
				}

				var buf bytes.Buffer
				zw := zip.NewWriter(&buf)
				_ = toZipFile(fmt.Sprintf("%s.pem", key.ID.String()), []byte(plain), zw)
				_ = toZipFile(fmt.Sprintf("%s.pub", key.ID.String()), []byte(key.PublicKey), zw)
				_ = zw.Close()

				b64 := utils.EncodeB64(buf.Bytes())
				resp.ZipBase64 = &b64
				resp.Filenames = []string{
					fmt.Sprintf("%s.zip", key.ID.String()),
					fmt.Sprintf("%s.pem", key.ID.String()),
					fmt.Sprintf("%s.pub", key.ID.String()),
				}
				utils.WriteJSON(w, http.StatusOK, resp)
				return

			default:
				utils.WriteError(w, http.StatusBadRequest, "invalid_ssh_part", "invalid part (public|private|both)")
				return
			}
		}

		switch part {
		case "public":
			filename := fmt.Sprintf("%s.pub", key.ID.String())
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
			_, _ = w.Write([]byte(key.PublicKey))
			return

		case "private":
			plain, err := utils.DecryptForOrg(orgID, key.EncryptedPrivateKey, key.PrivateIV, key.PrivateTag, db)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to decrypt ssh key")
				return
			}
			filename := fmt.Sprintf("%s.pem", key.ID.String())
			w.Header().Set("Content-Type", "application/x-pem-file")
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
			_, _ = w.Write([]byte(plain))
			return

		case "both":
			plain, err := utils.DecryptForOrg(orgID, key.EncryptedPrivateKey, key.PrivateIV, key.PrivateTag, db)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to decrypt ssh key")
				return
			}

			var buf bytes.Buffer
			zw := zip.NewWriter(&buf)
			_ = toZipFile(fmt.Sprintf("%s.pem", key.ID.String()), []byte(plain), zw)
			_ = toZipFile(fmt.Sprintf("%s.pub", key.ID.String()), []byte(key.PublicKey), zw)
			_ = zw.Close()

			filename := fmt.Sprintf("ssh_key_%s.zip", key.ID.String())
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
			_, _ = w.Write(buf.Bytes())
			return

		default:
			utils.WriteError(w, http.StatusBadRequest, "invalid_ssh_part", "invalid part (public|private|both)")
			return
		}
	}
}

// --- Helpers ---

func allowedBits(b int) bool {
	return b == 2048 || b == 3072 || b == 4096
}

func GenerateRSA(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func RSAPrivateToPEMAndAuthorized(priv *rsa.PrivateKey, comment string) (privPEM string, authorized string, err error) {
	der := x509.MarshalPKCS1PrivateKey(priv)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}
	var buf bytes.Buffer
	if err = pem.Encode(&buf, block); err != nil {
		return "", "", err
	}

	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return "", "", err
	}
	auth := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pub)))
	comment = strings.TrimSpace(comment)
	if comment != "" {
		auth += " " + comment
	}
	return buf.String(), auth, nil
}

func GenerateRSAPEMAndAuthorized(bits int, comment string) (string, string, error) {
	priv, err := GenerateRSA(bits)
	if err != nil {
		return "", "", err
	}
	return RSAPrivateToPEMAndAuthorized(priv, comment)
}

func toZipFile(filename string, content []byte, zw *zip.Writer) error {
	f, err := zw.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	return err
}

func keyFilenamePrefix(pubAuth string) string {
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubAuth))
	if err != nil {
		return "id_key"
	}
	switch pk.Type() {
	case "ssh-ed25519":
		return "id_ed25519"
	case "ssh-rsa":
		return "id_rsa"
	default:
		return "id_key"
	}
}

func GenerateEd25519PEMAndAuthorized(comment string) (privPEM string, authorized string, err error) {
	// Generate ed25519 keypair
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	// Private: PKCS#8 PEM
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", "", err
	}
	block := &pem.Block{Type: "PRIVATE KEY", Bytes: der}
	var buf bytes.Buffer
	if err := pem.Encode(&buf, block); err != nil {
		return "", "", err
	}

	// Public: OpenSSH authorized_key
	sshPub, err := ssh.NewPublicKey(ed25519.PublicKey(pub))
	if err != nil {
		return "", "", err
	}
	auth := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(sshPub)))
	comment = strings.TrimSpace(comment)
	if comment != "" {
		auth += " " + comment
	}

	return buf.String(), auth, nil
}
