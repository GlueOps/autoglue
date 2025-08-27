package servers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// ListServers godoc
// @Summary      List servers (org scoped)
// @Description  Returns servers for the organization in X-Org-ID. Optional filters: status, role.
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        status query string false "Filter by status (pending|provisioning|ready|failed)"
// @Param        role query string false "Filter by role"
// @Security     BearerAuth
// @Success      200 {array}  serverResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list servers"
// @Router       /api/v1/servers [get]
func ListServers(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	q := db.DB.Where("organization_id = ?", ac.OrganizationID)
	if s := strings.TrimSpace(r.URL.Query().Get("status")); s != "" {
		if !validStatus(s) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		q = q.Where("status = ?", strings.ToLower(s))
	}
	if role := strings.TrimSpace(r.URL.Query().Get("role")); role != "" {
		q = q.Where("role = ?", role)
	}

	var rows []models.Server
	if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list servers", http.StatusInternalServerError)
		return
	}

	out := make([]serverResponse, 0, len(rows))
	for _, s := range rows {
		out = append(out, toResponse(s))
	}
	writeJSON(w, http.StatusOK, out)
}

// GetServer godoc
// @Summary      Get server by ID (org scoped)
// @Description  Returns one server in the given organization.
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Server ID (UUID)"
// @Security     BearerAuth
// @Success      200 {object} serverResponse
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/servers/{id} [get]
func GetServer(w http.ResponseWriter, r *http.Request) {
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

	var s models.Server
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, toResponse(s))
}

// CreateServer godoc
// @Summary      Create server (org scoped)
// @Description  Creates a server bound to the org in X-Org-ID. Validates that ssh_key_id belongs to the org.
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createServerRequest true "Server payload"
// @Security     BearerAuth
// @Success      201 {object} serverResponse
// @Failure      400 {string} string "invalid json / missing fields / invalid status / invalid ssh_key_id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "create failed"
// @Router       /api/v1/servers [post]
func CreateServer(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.IPAddress == "" || req.SSHUser == "" || req.SshKeyID == "" || req.Role == "" {
		http.Error(w, "ip_address, ssh_user, ssh_key_id and role are required", http.StatusBadRequest)
		return
	}
	if req.Status != "" && !validStatus(req.Status) {
		http.Error(w, "invalid status", http.StatusBadRequest)
		return
	}

	keyID, err := uuid.Parse(req.SshKeyID)
	if err != nil {
		http.Error(w, "invalid ssh_key_id", http.StatusBadRequest)
		return
	}
	if err := ensureKeyBelongsToOrg(ac.OrganizationID, keyID); err != nil {
		http.Error(w, "invalid or unauthorized ssh_key_id", http.StatusBadRequest)
		return
	}

	s := models.Server{
		OrganizationID: ac.OrganizationID,
		Hostname:       req.Hostname,
		IPAddress:      req.IPAddress,
		SSHUser:        req.SSHUser,
		SshKeyID:       keyID,
		Role:           req.Role,
		Status:         "pending",
	}
	if req.Status != "" {
		s.Status = strings.ToLower(req.Status)
	}

	if err := db.DB.Create(&s).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(s))
}

// UpdateServer godoc
// @Summary      Update server (org scoped)
// @Description  Partially update fields; changing ssh_key_id validates ownership.
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Server ID (UUID)"
// @Param        body body updateServerRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} serverResponse
// @Failure      400 {string} string "invalid id / invalid json / invalid status / invalid ssh_key_id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "update failed"
// @Router       /api/v1/servers/{id} [patch]
func UpdateServer(w http.ResponseWriter, r *http.Request) {
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

	var s models.Server
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Hostname != nil {
		s.Hostname = *req.Hostname
	}
	if req.IPAddress != nil {
		s.IPAddress = *req.IPAddress
	}
	if req.SSHUser != nil {
		s.SSHUser = *req.SSHUser
	}
	if req.Role != nil {
		s.Role = *req.Role
	}
	if req.Status != nil {
		if !validStatus(*req.Status) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		s.Status = strings.ToLower(*req.Status)
	}
	if req.SshKeyID != nil {
		keyID, err := uuid.Parse(*req.SshKeyID)
		if err != nil {
			http.Error(w, "invalid ssh_key_id", http.StatusBadRequest)
			return
		}
		if err := ensureKeyBelongsToOrg(ac.OrganizationID, keyID); err != nil {
			http.Error(w, "invalid or unauthorized ssh_key_id", http.StatusBadRequest)
			return
		}
		s.SshKeyID = keyID
	}

	if err := db.DB.Save(&s).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, toResponse(s))
}

// DeleteServer godoc
// @Summary      Delete server (org scoped)
// @Description  Permanently deletes the server.
// @Tags         servers
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Server ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/servers/{id} [delete]
func DeleteServer(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Delete(&models.Server{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
