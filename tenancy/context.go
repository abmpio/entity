package tenancy

import (
	"context"
)

type (
	tenantResolveRes struct{}
	tenantConfigKey  string
)

// get tenant id from context
func TenantIdFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(_tenantIdKey).(string)
	if ok {
		return value, true
	}
	return value, false
}

// get tenant info from context
func TenantInfoFromContext(ctx context.Context) (ITenantInfo, bool) {
	value, ok := ctx.Value("tenantInfo").(ITenantInfo)
	if ok {
		return value, true
	}
	return value, false
}

func FromTenantResolveRes(ctx context.Context) (*TenantResolveResult, bool) {
	v, ok := ctx.Value(tenantResolveRes{}).(*TenantResolveResult)
	if ok {
		return v, ok
	}
	return nil, false
}

func FromTenantConfigContext(ctx context.Context, tenantId string) (*TenantConfig, bool) {
	v, ok := ctx.Value(tenantConfigKey(tenantId)).(*TenantConfig)
	if ok {
		return v, ok && v != nil
	}
	return nil, false
}
