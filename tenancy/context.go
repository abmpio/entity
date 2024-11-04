package tenancy

import (
	"context"
)

type (
	currentTenantCtx struct{}
)

func NewCurrentTenant(ctx context.Context, id, name string) context.Context {
	return NewCurrentTenantInfo(ctx, NewBasicTenantInfo(id, name))
}

func NewCurrentTenantInfo(ctx context.Context, info TenantInfo) context.Context {
	return context.WithValue(ctx, currentTenantCtx{}, info)
}

func FromCurrentTenant(ctx context.Context) (TenantInfo, bool) {
	value, ok := ctx.Value(currentTenantCtx{}).(TenantInfo)
	if ok {
		return value, true
	}
	return NewBasicTenantInfo("", ""), false
}

func TenantIdFromContext(ctx context.Context) string {
	tenantInfo, ok := FromCurrentTenant(ctx)
	if !ok {
		return ""
	}
	return tenantInfo.GetId()
}
