package tenancy

import (
	"context"
	"errors"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
)

// tenant 存储接口
type ITenantStore interface {
	// GetById return nil and ErrTenantNotFound if tenant not found
	GetById(ctx context.Context, id string) (*TenantConfig, error)
}

// ITenantStore的内存存储实现
type MemoryTenantStore struct {
	TenantConfig []TenantConfig
}

var _ ITenantStore = (*MemoryTenantStore)(nil)

func NewMemoryTenantStore(t []TenantConfig) *MemoryTenantStore {
	return &MemoryTenantStore{
		TenantConfig: t,
	}
}

// #region ITenantStore Members

func (m *MemoryTenantStore) GetById(_ context.Context, id string) (*TenantConfig, error) {
	for _, config := range m.TenantConfig {
		if config.Id == id || config.Name == id {
			return &config, nil
		}
	}
	return nil, ErrTenantNotFound
}

// #endregion
