package health

import (
	"net/http"

	"github.com/glueops/autoglue/internal/response"
)

// Check HealthCheck godoc
// @Summary      Basic health check
// @Description  Returns a 200 if the service is up
// @Tags         health
// @Accept       json
// @Produce      plain
// @Success      200 {string} string "ok"
// @Router       /api/healthz [get]
func Check(w http.ResponseWriter, r *http.Request) {
	_ = response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
