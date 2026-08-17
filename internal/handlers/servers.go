package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ListServers godoc
//
//	@ID				ListServers
//	@Summary		List servers (org scoped)
//	@Description	Returns servers for the organization in X-Org-ID. Optional filters: status, role.
//	@Tags			Servers
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			status		query		string	false	"Filter by status (pending|provisioning|ready|failed)"
//	@Param			role		query		string	false	"Filter by role"
//	@Success		200			{array}		dto.ServerResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"failed to list servers"
//	@Router			/servers [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListServers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		q := db.Where("organization_id = ?", orgID)

		if s := strings.TrimSpace(r.URL.Query().Get("status")); s != "" {
			if !validStatus(s) {
				utils.WriteError(w, http.StatusBadRequest, "status_invalid", "invalid status")
				return
			}
			q = q.Where("status = ?", strings.ToLower(s))
		}

		if role := strings.TrimSpace(r.URL.Query().Get("role")); role != "" {
			q = q.Where("role = ?", role)
		}

		var rows []models.Server
		if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to list servers")
			return
		}

		out := make([]dto.ServerResponse, 0, len(rows))
		for _, row := range rows {
			out = append(out, dto.ServerResponse{
				ID:               row.ID,
				OrganizationID:   row.OrganizationID,
				Hostname:         row.Hostname,
				PublicIPAddress:  row.PublicIPAddress,
				PrivateIPAddress: row.PrivateIPAddress,
				SSHUser:          row.SSHUser,
				SshKeyID:         row.SshKeyID,
				Role:             row.Role,
				Status:           row.Status,
				CreatedAt:        row.CreatedAt.UTC().Format(time.RFC3339),
				UpdatedAt:        row.UpdatedAt.UTC().Format(time.RFC3339),
			})
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetServer godoc
//
//	@ID				GetServer
//	@Summary		Get server by ID (org scoped)
//	@Description	Returns one server in the given organization.
//	@Tags			Servers
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"Server ID (UUID)"
//	@Success		200			{object}	dto.ServerResponse
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"fetch failed"
//	@Router			/servers/{id} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_invalid", "invalid id")
			return
		}

		var row models.Server
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to get server")
			return
		}

		utils.WriteJSON(w, http.StatusOK, row)
	}
}

// CreateServer godoc
//
//	@ID				CreateServer
//	@Summary		Create server (org scoped)
//	@Description	Creates a server bound to the org in X-Org-ID. Validates that ssh_key_id belongs to the org.
//	@Tags			Servers
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string					false	"Organization UUID"
//	@Param			body		body		dto.CreateServerRequest	true	"Server payload"
//	@Success		201			{object}	dto.ServerResponse
//	@Failure		400			{string}	string	"invalid json / missing fields / invalid status / invalid ssh_key_id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"create failed"
//	@Router			/servers [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var req dto.CreateServerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		req.Role = strings.ToLower(strings.TrimSpace(req.Role))
		req.Status = strings.ToLower(strings.TrimSpace(req.Status))
		pub := strings.TrimSpace(req.PublicIPAddress)

		if req.PrivateIPAddress == "" || req.SSHUser == "" || req.SshKeyID == "" || req.Role == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "private_ip_address, ssh_user, ssh_key_id and role are required")
			return
		}

		if req.Status != "" && !validStatus(req.Status) {
			utils.WriteError(w, http.StatusBadRequest, "status_invalid", "invalid status")
			return
		}

		if req.Role == "bastion" && pub == "" {
			utils.WriteError(w, http.StatusBadRequest, "public_ip_required", "public_ip_address is required for role=bastion")
			return
		}

		keyID, err := uuid.Parse(req.SshKeyID)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid ssh_key_id")
			return
		}
		if err := ensureKeyBelongsToOrg(orgID, keyID, db); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid or unauthorized ssh_key_id")
			return
		}

		var publicPtr *string
		if pub != "" {
			publicPtr = &pub
		}

		s := models.Server{
			OrganizationID:   orgID,
			Hostname:         req.Hostname,
			PublicIPAddress:  publicPtr,
			PrivateIPAddress: req.PrivateIPAddress,
			SSHUser:          req.SSHUser,
			SshKeyID:         keyID,
			Role:             req.Role,
			Status:           "pending",
		}
		if req.Status != "" {
			s.Status = strings.ToLower(req.Status)
		}

		if err := db.Create(&s).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to create server")
			return
		}
		utils.WriteJSON(w, http.StatusCreated, s)
	}
}

// UpdateServer godoc
//
//	@ID				UpdateServer
//	@Summary		Update server (org scoped)
//	@Description	Partially update fields; changing ssh_key_id validates ownership.
//	@Tags			Servers
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string					false	"Organization UUID"
//	@Param			id			path		string					true	"Server ID (UUID)"
//	@Param			body		body		dto.UpdateServerRequest	true	"Fields to update"
//	@Success		200			{object}	dto.ServerResponse
//	@Failure		400			{string}	string	"invalid id / invalid json / invalid status / invalid ssh_key_id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"update failed"
//	@Router			/servers/{id} [patch]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func UpdateServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_invalid", "invalid id")
			return
		}

		var server models.Server
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&server).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to get server")
			return
		}

		var req dto.UpdateServerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "bad request")
			return
		}

		next := server

		if req.Hostname != nil {
			next.Hostname = *req.Hostname
		}
		if req.PrivateIPAddress != nil {
			next.PrivateIPAddress = *req.PrivateIPAddress
		}
		if req.PublicIPAddress != nil {
			next.PublicIPAddress = req.PublicIPAddress
		}
		if req.SSHUser != nil {
			next.SSHUser = *req.SSHUser
		}
		if req.Role != nil {
			next.Role = *req.Role
		}
		if req.Status != nil {
			st := strings.ToLower(strings.TrimSpace(*req.Status))
			if !validStatus(st) {
				utils.WriteError(w, http.StatusBadRequest, "status_invalid", "invalid status")
				return
			}
			next.Status = st
		}
		if req.SshKeyID != nil {
			keyID, err := uuid.Parse(*req.SshKeyID)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid ssh_key_id")
				return
			}
			if err := ensureKeyBelongsToOrg(orgID, keyID, db); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid or unauthorized ssh_key_id")
				return
			}
			next.SshKeyID = keyID
		}

		if strings.EqualFold(next.Role, "bastion") &&
			(next.PublicIPAddress == nil || strings.TrimSpace(*next.PublicIPAddress) == "") {
			utils.WriteError(w, http.StatusBadRequest, "public_ip_required", "public_ip_address is required for role=bastion")
			return
		}

		if err := db.Save(&next).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to update server")
			return
		}
		utils.WriteJSON(w, http.StatusOK, server)
	}
}

// DeleteServer godoc
//
//	@ID				DeleteServer
//	@Summary		Delete server (org scoped)
//	@Description	Permanently deletes the server.
//	@Tags			Servers
//	@Produce		json
//	@Param			X-Org-ID	header	string	false	"Organization UUID"
//	@Param			id			path	string	true	"Server ID (UUID)"
//	@Success		204			"No Content"
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"delete failed"
//	@Router			/servers/{id} [delete]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func DeleteServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_invalid", "invalid id")
			return
		}

		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&models.Server{}).Error; err != nil {
			utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found")
			return
		}

		if err := db.Where("id = ? AND organization_id = ?", id, orgID).Delete(&models.Server{}).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to delete server")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ResetServerHostKey godoc
//
//	@ID				ResetServerHostKey
//	@Summary		Reset SSH host key (org scoped)
//	@Description	Clears the stored SSH host key for this server. The next SSH connection will re-learn the host key (trust-on-first-use).
//	@Tags			Servers
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"Server ID (UUID)"
//	@Success		200			{object}	dto.ServerResponse
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"reset failed"
//	@Router			/servers/{id}/reset-hostkey [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ResetServerHostKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_invalid", "invalid id")
			return
		}

		var server models.Server
		if err := db.Where("id = ? AND organization_id = ?", id, orgID).First(&server).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to get server")
			return
		}

		// Clear stored host key so next SSH handshake will TOFU and persist a new one.
		server.SSHHostKey = ""
		server.SSHHostKeyAlgo = ""

		if err := db.Save(&server).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to reset host key")
			return
		}

		utils.WriteJSON(w, http.StatusOK, server)
	}
}

// ReprovisionServer godoc
//
//	@ID				ReprovisionServer
//	@Summary		Reprovision a rebuilt server (org scoped)
//	@Description	Marks the server pending and clears its stored SSH host key in one transaction, so the bastion sweep re-bootstraps it and the next handshake re-learns the host key. Use this when the host has genuinely been rebuilt: it deliberately drops the trust-on-first-use anchor.
//	@Tags			Servers
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"Server ID (UUID)"
//	@Success		200			{object}	dto.ServerResponse
//	@Failure		400			{string}	string	"invalid id"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"reprovision failed"
//	@Router			/servers/{id}/reprovision [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
//
// The two effects belong together. Clearing the host key without marking the
// server pending leaves it trusting whatever answers next with nothing queued
// to reconnect; marking it pending without clearing the key leaves every
// attempt failing on mismatch forever. Half of this is worse than neither, so
// it is one statement in one transaction rather than two calls a caller can
// interleave or abandon.
//
// It is also deliberately a human action. A host key mismatch means either the
// host was rebuilt or something is impersonating it, and nothing in the
// database distinguishes those — so this must never be wired into an automated
// recovery sweep, which would reset the trust anchor precisely whenever
// something went wrong. The log line below is the record of who decided the
// rebuild was legitimate.
//
// Restricted to bastions, because on any other role both halves are no-ops.
// autoglue only ever opens an SSH connection to a bastion — every other server
// is reached *from* the bastion by the make targets, using the keys shipped in
// the payload — so a worker's ssh_host_key is never populated and there is
// nothing to clear. The sweep likewise claims only role='bastion', so pending
// queues nothing. Succeeding at neither and returning 200 would tell a caller
// their rebuilt host was requeued when it was not.
func ReprovisionServer(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_invalid", "invalid id")
			return
		}

		var server models.Server
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id = ? AND organization_id = ?", id, orgID).
				First(&server).Error; err != nil {
				return err
			}
			if !strings.EqualFold(server.Role, "bastion") {
				return errNotABastion
			}
			return tx.Model(&models.Server{}).
				Where("id = ? AND organization_id = ?", id, orgID).
				Updates(map[string]any{
					"status":            "pending",
					"ssh_host_key":      "",
					"ssh_host_key_algo": "",
					"updated_at":        time.Now(),
				}).Error
		}); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found")
				return
			}
			if errors.Is(err, errNotABastion) {
				utils.WriteError(w, http.StatusBadRequest, "not_a_bastion",
					"only bastions can be reprovisioned; other servers are reached from the bastion and have no stored host key")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "failed to reprovision server")
			return
		}

		// Security-relevant and not otherwise recoverable: a raw UPDATE leaves
		// no trace of who decided the host was legitimately rebuilt.
		event := log.Warn().
			Str("server_id", server.ID.String()).
			Str("organization_id", orgID.String()).
			Str("hostname", server.Hostname).
			Str("previous_status", server.Status).
			Bool("host_key_cleared", server.SSHHostKey != "")
		if user, ok := httpmiddleware.UserFrom(r.Context()); ok && user != nil {
			event = event.Str("actor_user_id", user.ID.String())
		}
		event.Msg("[servers] reprovision requested; host key trust anchor cleared")

		server.Status = "pending"
		server.SSHHostKey = ""
		server.SSHHostKeyAlgo = ""
		utils.WriteJSON(w, http.StatusOK, server)
	}
}

// --- Helpers ---

// errNotABastion is a sentinel rather than a response written inside the
// transaction, so the rollback happens before anything is sent to the client.
var errNotABastion = errors.New("server is not a bastion")

func validStatus(status string) bool {
	switch strings.ToLower(status) {
	case "pending", "provisioning", "ready", "failed", "":
		return true
	default:
		return false
	}
}

func ensureKeyBelongsToOrg(orgID, keyID uuid.UUID, db *gorm.DB) error {
	var k models.SshKey
	if err := db.Where("id = ? AND organization_id = ?", keyID, orgID).First(&k).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ssh key not found for this organization")
		}
		return err
	}
	return nil
}
