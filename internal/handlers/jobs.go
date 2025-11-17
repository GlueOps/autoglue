package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dyaksa/archer"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AdminListArcherJobs godoc
//
//	@ID				AdminListArcherJobs
//	@Summary		List Archer jobs (admin)
//	@Description	Paginated background jobs with optional filters. Search `q` may match id, type, error, payload (implementation-dependent).
//	@Tags			ArcherAdmin
//	@Produce		json
//	@Param			status		query		string	false	"Filter by status"	Enums(queued,running,succeeded,failed,canceled,retrying,scheduled)
//	@Param			queue		query		string	false	"Filter by queue name / worker name"
//	@Param			q			query		string	false	"Free-text search"
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Items per page"	minimum(1)	maximum(100)	default(25)
//	@Success		200			{object}	dto.PageJob
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"forbidden"
//	@Failure		500			{string}	string	"internal error"
//	@Router			/admin/archer/jobs [get]
//	@Security		BearerAuth
func AdminListArcherJobs(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := strings.TrimSpace(r.URL.Query().Get("status"))
		queue := strings.TrimSpace(r.URL.Query().Get("queue"))
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		page := atoiDefault(r.URL.Query().Get("page"), 1)
		size := clamp(atoiDefault(r.URL.Query().Get("page_size"), 25), 1, 100)

		base := db.Model(&models.Job{})
		if status != "" {
			base = base.Where("status = ?", status)
		}
		if queue != "" {
			base = base.Where("queue_name = ?", queue)
		}
		if q != "" {
			like := "%" + q + "%"
			base = base.Where(
				db.Where("id ILIKE ?", like).
					Or("queue_name ILIKE ?", like).
					Or("COALESCE(last_error,'') ILIKE ?", like).
					Or("CAST(arguments AS TEXT) ILIKE ?", like),
			)
		}

		var total int64
		if err := base.Count(&total).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		var rows []models.Job
		offset := (page - 1) * size
		if err := base.Order("created_at DESC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		items := make([]dto.Job, 0, len(rows))
		for _, m := range rows {
			items = append(items, mapModelJobToDTO(m))
		}

		utils.WriteJSON(w, http.StatusOK, dto.PageJob{
			Items:    items,
			Total:    int(total),
			Page:     page,
			PageSize: size,
		})
	}
}

// AdminEnqueueArcherJob godoc
//
//	@ID				AdminEnqueueArcherJob
//	@Summary		Enqueue a new Archer job (admin)
//	@Description	Create a job immediately or schedule it for the future via `run_at`.
//	@Tags			ArcherAdmin
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.EnqueueRequest	true	"Job parameters"
//	@Success		200		{object}	dto.Job
//	@Failure		400		{string}	string	"invalid json or missing fields"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		403		{string}	string	"forbidden"
//	@Failure		500		{string}	string	"internal error"
//	@Router			/admin/archer/jobs [post]
//	@Security		BearerAuth
func AdminEnqueueArcherJob(db *gorm.DB, jobs *bg.Jobs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in dto.EnqueueRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "invalid json")
			return
		}
		in.Queue = strings.TrimSpace(in.Queue)
		in.Type = strings.TrimSpace(in.Type)
		if in.Queue == "" || in.Type == "" {
			utils.WriteError(w, http.StatusBadRequest, "bad_request", "queue and type are required")
			return
		}

		// Parse payload into generic 'args' for Archer.
		var args any
		if len(in.Payload) > 0 && string(in.Payload) != "null" {
			if err := json.Unmarshal(in.Payload, &args); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "bad_request", "payload must be valid JSON")
				return
			}
		}

		id := uuid.NewString()

		opts := []archer.FnOptions{
			archer.WithMaxRetries(0), // adjust or expose in request if needed
		}
		if in.RunAt != nil {
			opts = append(opts, archer.WithScheduleTime(*in.RunAt))
		}

		// Schedule with Archer (queue == worker name).
		if _, err := jobs.Enqueue(context.Background(), id, in.Queue, args, opts...); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "enqueue_failed", err.Error())
			return
		}

		// Read back the just-created row.
		var m models.Job
		if err := db.First(&m, "id = ?", id).Error; err != nil {
			// Fallback: return a synthesized job if row not visible yet.
			now := time.Now()
			out := dto.Job{
				ID:          id,
				Type:        in.Type,
				Queue:       in.Queue,
				Status:      dto.StatusQueued,
				Attempts:    0,
				MaxAttempts: 0,
				CreatedAt:   now,
				UpdatedAt:   &now,
				RunAt:       in.RunAt,
				Payload:     args,
			}
			utils.WriteJSON(w, http.StatusOK, out)
			return
		}

		utils.WriteJSON(w, http.StatusOK, mapModelJobToDTO(m))
	}
}

// AdminRetryArcherJob godoc
//
//	@ID				AdminRetryArcherJob
//	@Summary		Retry a failed/canceled Archer job (admin)
//	@Description	Marks the job retriable (DB flip). Swap this for an Archer admin call if you expose one.
//	@Tags			ArcherAdmin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Job ID"
//	@Success		200	{object}	dto.Job
//	@Failure		400	{string}	string	"invalid job or not eligible"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"forbidden"
//	@Failure		404	{string}	string	"not found"
//	@Router			/admin/archer/jobs/{id}/retry [post]
//	@Security		BearerAuth
func AdminRetryArcherJob(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var m models.Job
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&m, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.WriteError(w, http.StatusNotFound, "not_found", "job not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		// Only allow retry from failed/canceled (adjust as you see fit).
		if m.Status != string(dto.StatusFailed) && m.Status != string(dto.StatusCanceled) {
			utils.WriteError(w, http.StatusBadRequest, "not_eligible", "job is not failed/canceled")
			return
		}

		// Reset to queued; clear started_at; bump updated_at.
		now := time.Now()
		if err := db.Model(&m).Updates(map[string]any{
			"status":     string(dto.StatusQueued),
			"started_at": nil,
			"updated_at": now,
		}).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		// Re-read and return.
		if err := db.First(&m, "id = ?", id).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, mapModelJobToDTO(m))
	}
}

// AdminCancelArcherJob godoc
//
//	@ID				AdminCancelArcherJob
//	@Summary		Cancel an Archer job (admin)
//	@Description	Set job status to canceled if cancellable. For running jobs, this only affects future picks; wire to Archer if you need active kill.
//	@Tags			ArcherAdmin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Job ID"
//	@Success		200	{object}	dto.Job
//	@Failure		400	{string}	string	"invalid job or not cancellable"
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"forbidden"
//	@Failure		404	{string}	string	"not found"
//	@Router			/admin/archer/jobs/{id}/cancel [post]
//	@Security		BearerAuth
func AdminCancelArcherJob(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var m models.Job
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&m, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.WriteError(w, http.StatusNotFound, "not_found", "job not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		// If already finished, bail.
		switch m.Status {
		case string(dto.StatusSucceeded), string(dto.StatusCanceled):
			utils.WriteError(w, http.StatusBadRequest, "not_cancellable", "job already finished")
			return
		}

		now := time.Now()
		if err := db.Model(&m).Updates(map[string]any{
			"status":     string(dto.StatusCanceled),
			"updated_at": now,
		}).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		if err := db.First(&m, "id = ?", id).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, mapModelJobToDTO(m))
	}
}

// AdminListArcherQueues godoc
//
//	@ID				AdminListArcherQueues
//	@Summary		List Archer queues (admin)
//	@Description	Summary metrics per queue (pending, running, failed, scheduled).
//	@Tags			ArcherAdmin
//	@Produce		json
//	@Success		200	{array}		dto.QueueInfo
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		403	{string}	string	"forbidden"
//	@Failure		500	{string}	string	"internal error"
//	@Router			/admin/archer/queues [get]
//	@Security		BearerAuth
func AdminListArcherQueues(db *gorm.DB) http.HandlerFunc {
	type row struct {
		QueueName string
		Pending   int
		Running   int
		Failed    int
		Scheduled int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var rows []row
		// Use filtered aggregate; adjust status values if your Archer differs.
		if err := db.
			Raw(`
				SELECT
					queue_name,
					COUNT(*) FILTER (WHERE status = 'queued')     AS pending,
					COUNT(*) FILTER (WHERE status = 'running')    AS running,
					COUNT(*) FILTER (WHERE status = 'failed')     AS failed,
					COUNT(*) FILTER (WHERE status = 'scheduled')  AS scheduled
				FROM jobs
				GROUP BY queue_name
				ORDER BY queue_name ASC
			`).Scan(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		out := make([]dto.QueueInfo, 0, len(rows))
		for _, r := range rows {
			out = append(out, dto.QueueInfo{
				Name:      r.QueueName,
				Pending:   r.Pending,
				Running:   r.Running,
				Failed:    r.Failed,
				Scheduled: r.Scheduled,
			})
		}

		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// Helpers
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

func mapModelJobToDTO(m models.Job) dto.Job {
	var payload any
	if len(m.Arguments) > 0 {
		_ = json.Unmarshal([]byte(m.Arguments), &payload)
	}

	var updated *time.Time
	if !m.UpdatedAt.IsZero() {
		updated = &m.UpdatedAt
	}

	var runAt *time.Time
	if !m.ScheduledAt.IsZero() {
		rt := m.ScheduledAt
		runAt = &rt
	}

	return dto.Job{
		ID: m.ID,
		// If you distinguish between queue and type elsewhere, set Type accordingly.
		Type:        m.QueueName,
		Queue:       m.QueueName,
		Status:      dto.JobStatus(m.Status),
		Attempts:    m.RetryCount,
		MaxAttempts: m.MaxRetry,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   updated,
		LastError:   m.LastError,
		RunAt:       runAt,
		Payload:     payload,
	}
}
