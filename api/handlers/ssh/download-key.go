package ssh

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

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

	id, err := uuid.Parse(mux.Vars(r)["id"])
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
