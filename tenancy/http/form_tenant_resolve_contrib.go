package http

import (
	"context"
	"net/http"

	"github.com/abmpio/entity/tenancy"
)

type FormTenantResolveContrib struct {
	key     string
	request *http.Request
}

var _ tenancy.ITenantResolveContrib = (*FormTenantResolveContrib)(nil)

func NewFormTenantResolveContrib(key string, r *http.Request) *FormTenantResolveContrib {
	return &FormTenantResolveContrib{
		key:     key,
		request: r,
	}
}

// #region ITenantResolveContrib  Members

func (h *FormTenantResolveContrib) Name() string {
	return "Form"
}

func (h *FormTenantResolveContrib) Resolve(ctx context.Context) (string, error) {
	return h.request.FormValue(h.key), nil
}

// #endregion
