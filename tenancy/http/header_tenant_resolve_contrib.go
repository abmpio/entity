package http

import (
	"context"
	"net/http"

	"github.com/abmpio/entity/tenancy"
)

type HeaderTenantResolveContrib struct {
	key     string
	request *http.Request
}

var _ tenancy.ITenantResolveContrib = (*HeaderTenantResolveContrib)(nil)

func NewHeaderTenantResolveContrib(key string, r *http.Request) *HeaderTenantResolveContrib {
	return &HeaderTenantResolveContrib{
		key:     key,
		request: r,
	}
}

// #region ITenantResolveContrib  Members

func (h *HeaderTenantResolveContrib) Name() string {
	return "Header"
}

func (h *HeaderTenantResolveContrib) Resolve(ctx context.Context) (string, error) {
	return h.request.Header.Get(h.key), nil
}

// #endregion
