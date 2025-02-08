package tenancy

import (
	"context"
)

type ITenantResolveContrib interface {
	// Name of resolver
	Name() string
	// Resolve tenant
	Resolve(ctx context.Context) (string, error)
}

// ContextContrib resolve from current context
type ContextContrib struct {
}

var _ ITenantResolveContrib = (*ContextContrib)(nil)

func (c *ContextContrib) Name() string {
	return "ContextContrib"
}

func (c *ContextContrib) Resolve(ctx context.Context) (string, error) {
	tenantId, ok := TenantIdFromContext(ctx)
	if ok {
		return tenantId, nil
	}
	return "", nil
}
