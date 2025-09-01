package ctxutil

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const (
	keyUserID ctxKey = "user_id"
	keyOrgID  ctxKey = "org_id"
)

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(keyUserID).(uuid.UUID)
	return v, ok
}

func WithOrgID(ctx context.Context, orgID uuid.UUID) context.Context {
	return context.WithValue(ctx, keyOrgID, orgID)
}

func OrgID(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(keyOrgID).(uuid.UUID)
	return v, ok
}
