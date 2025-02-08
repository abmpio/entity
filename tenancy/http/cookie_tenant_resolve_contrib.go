package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/abmpio/entity/tenancy"
)

type CookieTenantResolveContrib struct {
	key     string
	request *http.Request
}

var _ tenancy.ITenantResolveContrib = (*CookieTenantResolveContrib)(nil)

func NewCookieTenantResolveContrib(key string, r *http.Request) *CookieTenantResolveContrib {
	return &CookieTenantResolveContrib{
		key:     key,
		request: r,
	}
}

// #region ITenantResolveContrib  Members

func (h *CookieTenantResolveContrib) Name() string {
	return "Cookie"
}

func (h *CookieTenantResolveContrib) Resolve(ctx context.Context) (string, error) {
	v, err := h.request.Cookie(h.key)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		// other error
		return "", err
	}
	return v.Value, nil
}

// #endregion
