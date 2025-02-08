package http

import (
	"context"
	"net/http"

	"github.com/abmpio/entity/tenancy"
)

type QueryTenantResolveContrib struct {
	key     string
	request *http.Request
}

var _ tenancy.ITenantResolveContrib = (*QueryTenantResolveContrib)(nil)

func NewQueryTenantResolveContrib(key string, r *http.Request) *QueryTenantResolveContrib {
	return &QueryTenantResolveContrib{
		key:     key,
		request: r,
	}
}

// #region ITenantResolveContrib  Members

func (h *QueryTenantResolveContrib) Name() string {
	return "Query"
}

func (h *QueryTenantResolveContrib) Resolve(ctx context.Context) (string, error) {
	return h.request.URL.Query().Get(h.key), nil
}

// #endregion
