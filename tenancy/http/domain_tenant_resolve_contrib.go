package http

import (
	"context"
	"net/http"
	"regexp"

	"github.com/abmpio/entity/tenancy"
)

type DomainTenantResolveContrib struct {
	request *http.Request
	format  string
}

var _ tenancy.ITenantResolveContrib = (*DomainTenantResolveContrib)(nil)

func NewDomainTenantResolveContrib(f string, r *http.Request) *DomainTenantResolveContrib {
	return &DomainTenantResolveContrib{
		request: r,
		format:  f,
	}
}

// #region ITenantResolveContrib  Members

func (h *DomainTenantResolveContrib) Name() string {
	return "Domain"
}

func (h *DomainTenantResolveContrib) Resolve(ctx context.Context) (string, error) {
	host := h.request.Host
	r := regexp.MustCompile(h.format)
	f := r.FindAllStringSubmatch(host, -1)
	if f == nil {
		//no match
		return "", nil
	}
	return f[0][1], nil
}

// #endregion
