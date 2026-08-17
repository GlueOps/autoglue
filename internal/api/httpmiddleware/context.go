package httpmiddleware

import (
	"context"

	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
)

type ctxKey string

const (
	ctxUserKey  ctxKey = "ctx_user"
	ctxOrgKey   ctxKey = "ctx_org"
	ctxRolesKey ctxKey = "ctx_roles" // []string, user roles in current org
	ctxAgentKey ctxKey = "ctx_agent" // *models.Agent, cluster-scoped machine principal
)

func WithUser(ctx context.Context, u *models.User) context.Context {
	return context.WithValue(ctx, ctxUserKey, u)
}
func WithOrg(ctx context.Context, o *models.Organization) context.Context {
	return context.WithValue(ctx, ctxOrgKey, o)
}
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, ctxRolesKey, roles)
}

// WithAgent carries the authenticated bastion agent. An agent is never also a
// user, so nothing sets this alongside WithUser.
func WithAgent(ctx context.Context, a *models.Agent) context.Context {
	return context.WithValue(ctx, ctxAgentKey, a)
}

func UserFrom(ctx context.Context) (*models.User, bool) {
	u, ok := ctx.Value(ctxUserKey).(*models.User)
	return u, ok && u != nil
}
func OrgFrom(ctx context.Context) (*models.Organization, bool) {
	o, ok := ctx.Value(ctxOrgKey).(*models.Organization)
	return o, ok && o != nil
}
func OrgIDFrom(ctx context.Context) (uuid.UUID, bool) {
	if o, ok := OrgFrom(ctx); ok {
		return o.ID, true
	}
	return uuid.Nil, false
}
func RolesFrom(ctx context.Context) ([]string, bool) {
	r, ok := ctx.Value(ctxRolesKey).([]string)
	return r, ok && r != nil
}
func AgentFrom(ctx context.Context) (*models.Agent, bool) {
	a, ok := ctx.Value(ctxAgentKey).(*models.Agent)
	return a, ok && a != nil
}

// AgentClusterIDFrom is the authorization scope of an agent request. Every
// agent-plane query filters on this rather than on the org, so it is worth
// having as one call that cannot be got wrong.
func AgentClusterIDFrom(ctx context.Context) (uuid.UUID, bool) {
	if a, ok := AgentFrom(ctx); ok {
		return a.ClusterID, true
	}
	return uuid.Nil, false
}
