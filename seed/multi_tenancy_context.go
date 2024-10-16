package seed

import (
	"context"
)

type (
	currentTenantCtx struct{}
)

func newCurrentTenant(ctx context.Context, id, name string) context.Context {
	return newCurrentTenantInfo(ctx, NewBasicTenantInfo(id, name))
}

func newCurrentTenantInfo(ctx context.Context, info tenantInfo) context.Context {
	return context.WithValue(ctx, currentTenantCtx{}, info)
}
