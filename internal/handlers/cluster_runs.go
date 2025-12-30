package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/dyaksa/archer"
	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListClusterRuns godoc
//
//	@ID				ListClusterRuns
//	@Summary		List cluster runs (org scoped)
//	@Description	Returns runs for a cluster within the organization in X-Org-ID.
//	@Tags			ClusterRuns
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Success		200			{array}		dto.ClusterRunResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/runs [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListClusterRuns(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		// Ensure cluster exists + org scoped
		if err := db.Select("id").
			Where("id = ? AND organization_id = ?", clusterID, orgID).
			First(&models.Cluster{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		var rows []models.ClusterRun
		if err := db.
			Where("organization_id = ? AND cluster_id = ?", orgID, clusterID).
			Order("created_at DESC").
			Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.ClusterRunResponse, 0, len(rows))
		for _, cr := range rows {
			out = append(out, clusterRunToDTO(cr))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetClusterRun godoc
//
//	@ID				GetClusterRun
//	@Summary		Get a cluster run (org scoped)
//	@Description	Returns a single run for a cluster within the organization in X-Org-ID.
//	@Tags			ClusterRuns
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Param			runID		path		string	true	"Run ID"
//	@Success		200			{object}	dto.ClusterRunResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/runs/{runID} [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetClusterRun(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		runID, err := uuid.Parse(chi.URLParam(r, "runID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_run_id", "invalid run id")
			return
		}

		var row models.ClusterRun
		if err := db.
			Where("id = ? AND organization_id = ? AND cluster_id = ?", runID, orgID, clusterID).
			First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "run not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, clusterRunToDTO(row))
	}
}

// RunClusterAction godoc
//
//	@ID				RunClusterAction
//	@Summary		Run an admin-configured action on a cluster (org scoped)
//	@Description	Creates a ClusterRun record for the cluster/action. Execution is handled asynchronously by workers.
//	@Tags			ClusterRuns
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Param			actionID	path		string	true	"Action ID"
//	@Success		201			{object}	dto.ClusterRunResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"cluster or action not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/actions/{actionID}/runs [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func RunClusterAction(db *gorm.DB, jobs *bg.Jobs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		clusterID, err := uuid.Parse(chi.URLParam(r, "clusterID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_cluster_id", "invalid cluster id")
			return
		}

		actionID, err := uuid.Parse(chi.URLParam(r, "actionID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_action_id", "invalid action id")
			return
		}

		// cluster must exist + org scoped
		var cluster models.Cluster
		if err := db.Select("id", "organization_id").
			Where("id = ? AND organization_id = ?", clusterID, orgID).
			First(&cluster).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "cluster not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		// action is global/admin-configured (not org scoped)
		var action models.Action
		if err := db.Where("id = ?", actionID).First(&action).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "action_not_found", "action not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		run := models.ClusterRun{
			OrganizationID: orgID,
			ClusterID:      clusterID,
			Action:         action.MakeTarget, // this is what you actually execute
			Status:         models.ClusterRunStatusQueued,
			Error:          "",
			FinishedAt:     time.Time{},
		}

		if err := db.Create(&run).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		args := bg.ClusterActionArgs{
			OrgID:      orgID,
			ClusterID:  clusterID,
			Action:     action.MakeTarget,
			MakeTarget: action.MakeTarget,
		}
		// Enqueue with run.ID as the job ID so the worker can look it up.
		_, enqueueErr := jobs.Enqueue(
			r.Context(),
			run.ID.String(),
			"cluster_action",
			args,
			archer.WithMaxRetries(0),
		)

		if enqueueErr != nil {
			_ = db.Model(&models.ClusterRun{}).
				Where("id = ?", run.ID).
				Updates(map[string]any{
					"status":      models.ClusterRunStatusFailed,
					"error":       "failed to enqueue job: " + enqueueErr.Error(),
					"finished_at": time.Now().UTC(),
				}).Error

			utils.WriteError(w, http.StatusInternalServerError, "job_error", "failed to enqueue cluster action")
			return
		}
		utils.WriteJSON(w, http.StatusCreated, clusterRunToDTO(run))

	}
}

func clusterRunToDTO(cr models.ClusterRun) dto.ClusterRunResponse {
	var finished *time.Time
	if !cr.FinishedAt.IsZero() {
		t := cr.FinishedAt
		finished = &t
	}
	return dto.ClusterRunResponse{
		ID:             cr.ID,
		OrganizationID: cr.OrganizationID,
		ClusterID:      cr.ClusterID,
		Action:         cr.Action,
		Status:         cr.Status,
		Error:          cr.Error,
		CreatedAt:      cr.CreatedAt,
		UpdatedAt:      cr.UpdatedAt,
		FinishedAt:     finished,
	}
}
