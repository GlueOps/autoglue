package jobs

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dyaksa/archer"
	"github.com/glueops/autoglue/internal/bg"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type JobListItem struct {
	ID           string     `json:"id"`
	QueueName    string     `json:"queue_name"`
	Status       string     `json:"status"`
	RetryCount   int        `json:"retry_count"`
	MaxRetry     int        `json:"max_retry"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastError    *string    `json:"last_error,omitempty"`
	ResultStatus string     `json:"result_status"`
	Processed    int        `json:"processed"`
	Ready        int        `json:"ready"`
	Failed       int        `json:"failed"`
	ElapsedMs    int        `json:"elapsed_ms"`
}

type EnqueueReq struct {
	Queue      string          `json:"queue"`
	Args       json.RawMessage `json:"args"` // keep raw and pass through to Archer
	MaxRetries *int            `json:"max_retries,omitempty"`
	ScheduleAt *time.Time      `json:"schedule_at,omitempty"`
}

type EnqueueResp struct {
	ID string `json:"id"`
}

func parseLimit(r *http.Request, def int) int {
	if s := r.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 && n <= 1000 {
			return n
		}
	}
	return def
}

func isNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return msg == "job not found" || msg == "no rows in result set"
}

// ---------------------- READ ENDPOINTS ----------------------

// GetKPI godoc
// @Summary      Jobs KPI
// @Description  Aggregated counters across all queues
// @Tags         jobs
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} jobs.KPI
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Router       /api/v1/jobs/kpi [get]
func GetKPI(w http.ResponseWriter, r *http.Request) {
	if middleware.GetAuthContext(r) == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	k, err := LoadKPI(db.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, k)
}

// GetQueues godoc
// @Summary      Per-queue rollups
// @Description  Counts and avg duration per queue (last 24h)
// @Tags         jobs
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} jobs.QueueRollup
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Router       /api/v1/jobs/queues [get]
func GetQueues(w http.ResponseWriter, r *http.Request) {
	if middleware.GetAuthContext(r) == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := LoadPerQueue(db.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, rows)
}

// GetActive godoc
// @Summary      Active jobs
// @Description  Currently running jobs (limit default 100)
// @Tags         jobs
// @Security     BearerAuth
// @Produce      json
// @Param        limit query int false "Max rows" default(100)
// @Success      200 {array} jobs.JobListItem
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Router       /api/v1/jobs/active [get]
func GetActive(w http.ResponseWriter, r *http.Request) {
	if middleware.GetAuthContext(r) == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit := parseLimit(r, 100)

	var rows []JobListItem
	err := db.DB.Model(&models.Job{}).
		Select(`
			id, queue_name, status, retry_count, max_retry, scheduled_at, started_at, updated_at, last_error,
			COALESCE(result->>'status','')                  AS result_status,
			COALESCE((result->>'processed')::int, 0)        AS processed,
			COALESCE((result->>'ready')::int, 0)            AS ready,
			COALESCE((result->>'failed')::int, 0)           AS failed,
			COALESCE((result->>'elapsed_ms')::int, 0)       AS elapsed_ms
		`).
		Where("status = ?", "running").
		Order("started_at DESC NULLS LAST, updated_at DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, rows)
}

// GetFailures godoc
// @Summary      Recent failures
// @Description  Failed jobs ordered by most recent (limit default 100)
// @Tags         jobs
// @Security     BearerAuth
// @Produce      json
// @Param        limit query int false "Max rows" default(100)
// @Success      200 {array} jobs.JobListItem
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Router       /api/v1/jobs/failures [get]
func GetFailures(w http.ResponseWriter, r *http.Request) {
	if middleware.GetAuthContext(r) == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit := parseLimit(r, 100)

	var rows []JobListItem
	err := db.DB.Model(&models.Job{}).
		Select(`
			id, queue_name, status, retry_count, max_retry, scheduled_at, started_at, updated_at, last_error,
			COALESCE(result->>'status','')                  AS result_status,
			COALESCE((result->>'processed')::int, 0)        AS processed,
			COALESCE((result->>'ready')::int, 0)            AS ready,
			COALESCE((result->>'failed')::int, 0)           AS failed,
			COALESCE((result->>'elapsed_ms')::int, 0)       AS elapsed_ms
		`).
		Where("status = ?", "failed").
		Order("updated_at DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, rows)
}

// ---------------------- MUTATION ENDPOINTS ----------------------

// RetryNow godoc
// @Summary      Retry a job immediately
// @Description  Calls Archer ScheduleNow on the job id
// @Tags         jobs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      204 {string} string "no content"
// @Failure      400 {string} string "bad request"
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "internal error"
// @Router       /api/v1/jobs/{id}/retry [post]
func RetryNow(w http.ResponseWriter, r *http.Request) {
	if middleware.GetAuthContext(r) == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// archer.ScheduleNow returns (any, error); if the id is unknown, expect an error you can surface as 404
	if _, err := bg.BgJobs.Client.ScheduleNow(r.Context(), id); err != nil {
		status := http.StatusInternalServerError
		// (Optional) map error text if Archer returns a recognizable "not found"
		if isNotFoundErr(err) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}
	response.NoContent(w)
}

// Cancel godoc
// @Summary      Cancel a job
// @Description  Cancels running or scheduled jobs
// @Tags         jobs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Job ID"
// @Success      204 {string} string "no content"
// @Failure      400 {string} string "bad request"
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "internal error"
// @Router       /api/v1/jobs/{id}/cancel [post]
func Cancel(w http.ResponseWriter, r *http.Request) {
	if middleware.GetAuthContext(r) == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	if _, err := bg.BgJobs.Client.Cancel(r.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if isNotFoundErr(err) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}
	response.NoContent(w)
}

// Enqueue godoc
// @Summary      Manually enqueue a job
// @Description  Schedules a job on a queue with optional args/schedule
// @Tags         jobs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload body jobs.EnqueueReq true "Enqueue request"
// @Success      202 {object} jobs.EnqueueResp
// @Failure      400 {string} string "bad request"
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Router       /api/v1/jobs/enqueue [post]
func Enqueue(w http.ResponseWriter, r *http.Request) {
	if middleware.GetAuthContext(r) == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req EnqueueReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Queue == "" {
		http.Error(w, "queue is required", http.StatusBadRequest)
		return
	}

	id := uuid.NewString()
	opts := []archer.FnOptions{}
	if req.MaxRetries != nil {
		opts = append(opts, archer.WithMaxRetries(*req.MaxRetries))
	}
	if req.ScheduleAt != nil {
		opts = append(opts, archer.WithScheduleTime(*req.ScheduleAt))
	}

	if _, err := bg.BgJobs.Client.Schedule(r.Context(), id, req.Queue, req.Args, opts...); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusAccepted, EnqueueResp{ID: id})
}
