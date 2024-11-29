package tenancy

import (
	"context"
)

func GetMultiTenantSide(ctx context.Context) MultiTenancySide {
	tenantId, _ := TenantIdFromContext(ctx)
	if tenantId == "" {
		return Host
	} else {
		return Tenant
	}
}

const DefaultKey = "Tenant"

var (
	_tenantIdKey string = DefaultKey
)

func SetTenantIdKey(v string) {
	_tenantIdKey = v
}

func KeyOrDefault(key string) string {
	if len(key) > 0 {
		return key
	}
	return DefaultKey
}
