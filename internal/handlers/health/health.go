package health

import (
	"net/http"

	"github.com/glueops/autoglue/internal/response"
)

// HealthStatus is the JSON shape returned by the health check.
type HealthStatus struct {
	Status string `json:"status" example:"ok"`
}

// Check godoc
// @Summary      Basic health check
// @Description  Returns 200 OK when the service is up
// @Tags         health
// @ID           HealthCheck               // operationId
// @Accept       json
// @Produce      json
// @Success      200 {object} health.HealthStatus
// @Router       /api/healthz [get]
func Check(w http.ResponseWriter, r *http.Request) {
	_ = response.JSON(w, http.StatusOK, HealthStatus{Status: "ok"})
}
