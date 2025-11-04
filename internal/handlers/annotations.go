package handlers

// ListAnnotations godoc
// @ID           ListAnnotations
// @Summary      List annotations (org scoped)
// @Description  Returns annotations for the organization in X-Org-ID. Filters: `name`, `value`, and `q` (name contains). Add `include=node_pools` to include linked node pools.
// @Tags         Annotations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Organization UUID"
// @Param        name query string false "Exact name"
// @Param        value query string false "Exact value"
// @Param        q query string false "name contains (case-insensitive)"
// @Success      200 {array}  annotationResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list annotations"
// @Router       /api/v1/annotations [get]
// @Security     BearerAuth
// @Security     OrgKeyAuth
// @Security     OrgSecretAuth
