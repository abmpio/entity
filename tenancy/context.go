package tenancy

import (
	"context"
)

func TenantIdFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(_tenantIdKey).(string)
	if ok {
		return value, true
	}
	return value, false
}
