package handlers

import (
	"net/http"

	"github.com/glueops/autoglue/internal/utils"
)

type HealthStatus struct {
	Status string `json:"status" example:"ok"`
}

// HealthCheck godoc
// @Summary      Basic health check
// @Description  Returns 200 OK when the service is up
// @Tags         Health
// @ID           HealthCheck               // operationId
// @Accept       json
// @Produce      json
// @Success      200 {object} HealthStatus
// @Router       /healthz [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, HealthStatus{Status: "ok"})
}
