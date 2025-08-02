package health

import "net/http"

// HealthCheck godoc
// @Summary      Basic health check
// @Description  Returns a 200 if the service is up
// @Tags         Health
// @Accept       json
// @Produce      plain
// @Success      200 {string} string "ok"
// @Router       /api/healthz [get]
func Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
