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
