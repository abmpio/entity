package http

import "github.com/abmpio/entity/tenancy"

type WebMultiTenancyOption struct {
	TenantKey    string
	DomainFormat string
}

func NewWebMultiTenancyOption(key string, domainFormat string) *WebMultiTenancyOption {
	key = tenancy.KeyOrDefault(key)
	return &WebMultiTenancyOption{
		TenantKey:    key,
		DomainFormat: domainFormat,
	}
}

func NewDefaultWebMultiTenancyOption() *WebMultiTenancyOption {
	return NewWebMultiTenancyOption("", "")
}
