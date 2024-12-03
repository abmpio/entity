package tenancy

import (
	"context"
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
