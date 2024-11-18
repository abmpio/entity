package tenancy

import (
	"context"
)

func GetMultiTenantSide(ctx context.Context) MultiTenancySide {
	tenantInfo, _ := FromCurrentTenant(ctx)
	if tenantInfo.GetId() == "" {
		return Host
	} else {
		return Tenant
	}
}

const DefaultKey = "__tenant"

func KeyOrDefault(key string) string {
	if len(key) > 0 {
		return key
	}
	return DefaultKey
}
