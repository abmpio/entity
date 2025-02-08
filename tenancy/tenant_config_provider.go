package tenancy

import (
	"context"
	"errors"
)

// TenantConfigProvider resolve tenant config from current context
type ITenantConfigProvider interface {
	// Get tenant config
	Get(ctx context.Context) (*TenantConfig, error)
}

type DefaultTenantConfigProvider struct {
	tr ITenantResolver
	ts ITenantStore
}

func NewDefaultTenantConfigProvider(tr ITenantResolver, ts ITenantStore) ITenantConfigProvider {
	return &DefaultTenantConfigProvider{
		tr: tr,
		ts: ts,
	}
}

// #region ITenantConfigProvider members

// Get read from context FromTenantConfigContext first, fallback with TenantStore and return new context with cached value
func (d *DefaultTenantConfigProvider) Get(ctx context.Context) (*TenantConfig, error) {
	rr, err := d.tr.Resolve(ctx)
	if err != nil {
		return &TenantConfig{}, err
	}
	if rr.TenantId != "" {
		//tenant side

		//read from cache
		if cfg, ok := FromTenantConfigContext(ctx, rr.TenantId); ok {
			return cfg, nil
		}
		//get config from tenant store
		cfg, err := d.ts.GetById(ctx, rr.TenantId)
		if err != nil {
			if errors.Is(err, ErrTenantNotFound) {
				return &TenantConfig{}, nil
			}
			return &TenantConfig{}, err
		}
		return cfg, nil
	}
	// host side
	return &TenantConfig{}, nil

}

// #endregion
