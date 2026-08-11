package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	jobLogDefaultLimit = 200
	jobLogMaxLimit     = 1000
)

// readJobLogs returns one cursor page of output for a subject, newest last.
func readJobLogs(db *gorm.DB, orgID uuid.UUID, subjectType string, subjectID uuid.UUID, after int64, limit int) ([]dto.JobLogChunk, int64, error) {
	var rows []models.JobLog
	if err := db.
		Where("organization_id = ? AND subject_type = ? AND subject_id = ? AND id > ?",
			orgID, subjectType, subjectID, after).
		Order("id ASC").
		Limit(limit).
		Find(&rows).Error; err != nil {
		return nil, after, err
	}

	items := make([]dto.JobLogChunk, 0, len(rows))
	cursor := after
	for _, r := range rows {
		items = append(items, dto.JobLogChunk{
			ID:        r.ID,
			Stream:    r.Stream,
			Chunk:     r.Chunk,
			CreatedAt: r.CreatedAt,
		})
		cursor = r.ID
	}
	return items, cursor, nil
}

// atoiDefault and clamp were generic helpers that happened to live in the
// archer admin handlers, and went with them when those were removed.
func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}

func clamp(n, lo, hi int) int {
	if n < lo {
		return lo
	}
	if n > hi {
		return hi
	}
	return n
}

func jobLogQuery(r *http.Request) (after int64, limit int) {
	after, _ = strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("after")), 10, 64)
	if after < 0 {
		after = 0
	}
	limit = atoiDefault(r.URL.Query().Get("limit"), jobLogDefaultLimit)
	return after, clamp(limit, 1, jobLogMaxLimit)
}

// GetClusterRunLogs godoc
//
//	@ID				GetClusterRunLogs
//	@Summary		Tail the output of a cluster run
//	@Description	Returns output produced by the run's background job, in order. Poll by passing the previous `next_cursor` as `after`; stop when `done` is true.
//	@Tags			ClusterRuns
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			clusterID	path		string	true	"Cluster ID"
//	@Param			runID		path		string	true	"Run ID"
//	@Param			after		query		int		false	"Return chunks with an id greater than this"	default(0)
//	@Param			limit		query		int		false	"Maximum chunks to return"						minimum(1)	maximum(1000)	default(200)
//	@Success		200			{object}	dto.JobLogPage
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/clusters/{clusterID}/runs/{runID}/logs [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetClusterRunLogs(db *gorm.DB) http.HandlerFunc {
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

		// The run must exist and belong to this org and cluster. Without this
		// the org filter on job_logs alone would let a caller read any run in
		// their own org by guessing ids.
		var run models.ClusterRun
		if err := db.
			Where("id = ? AND organization_id = ? AND cluster_id = ?", runID, orgID, clusterID).
			First(&run).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "run not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		after, limit := jobLogQuery(r)
		items, cursor, err := readJobLogs(db, orgID, models.JobLogSubjectClusterRun, runID, after, limit)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		done := run.Status != models.ClusterRunStatusRunning && run.Status != models.ClusterRunStatusQueued

		utils.WriteJSON(w, http.StatusOK, dto.JobLogPage{
			Items:      items,
			NextCursor: cursor,
			// Only report done once the reader has caught up, otherwise a
			// client that stops on `done` can miss the final chunks.
			Done: done && len(items) < limit,
		})
	}
}

// GetServerLogs godoc
//
//	@ID				GetServerLogs
//	@Summary		Tail the bootstrap output of a server
//	@Description	Returns output produced by background jobs acting on this server, most usefully the bastion bootstrap. Poll by passing the previous `next_cursor` as `after`; stop when `done` is true.
//	@Tags			Servers
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			id			path		string	true	"Server ID"
//	@Param			after		query		int		false	"Return chunks with an id greater than this"	default(0)
//	@Param			limit		query		int		false	"Maximum chunks to return"						minimum(1)	maximum(1000)	default(200)
//	@Success		200			{object}	dto.JobLogPage
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/servers/{id}/logs [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func GetServerLogs(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		serverID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "id_invalid", "invalid id")
			return
		}

		var server models.Server
		if err := db.
			Where("id = ? AND organization_id = ?", serverID, orgID).
			First(&server).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "server_not_found", "server not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		after, limit := jobLogQuery(r)
		items, cursor, err := readJobLogs(db, orgID, models.JobLogSubjectServer, serverID, after, limit)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		// "pending" means it is waiting for a tick to claim it, so more output
		// is still coming; only a terminal status ends the stream.
		status := strings.ToLower(server.Status)
		done := status == "ready" || status == "failed"

		utils.WriteJSON(w, http.StatusOK, dto.JobLogPage{
			Items:      items,
			NextCursor: cursor,
			Done:       done && len(items) < limit,
		})
	}
}
