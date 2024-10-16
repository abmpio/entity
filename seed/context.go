package seed

type Context struct {
	TenantId string
	Extra    map[string]interface{}
}

func NewSeedContext(tenantId string, extra map[string]interface{}) *Context {
	return &Context{
		TenantId: tenantId,
		Extra:    extra,
	}
}

func (s *Context) WithExtra(k string, v interface{}) *Context {
	if s.Extra == nil {
		s.Extra = make(map[string]interface{})
	}
	s.Extra[k] = v
	return s
}
