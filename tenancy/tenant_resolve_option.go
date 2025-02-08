package tenancy

type TenantResolveOption struct {
	Resolvers []ITenantResolveContrib
}

type ResolveOption func(resolveOption *TenantResolveOption)

func AppendContribs(c ...ITenantResolveContrib) ResolveOption {
	return func(resolveOption *TenantResolveOption) {
		resolveOption.AppendContribs(c...)
	}
}

func RemoveContribs(c ...ITenantResolveContrib) ResolveOption {
	return func(resolveOption *TenantResolveOption) {
		resolveOption.RemoveContribs(c...)
	}
}

func NewTenantResolveOption(c ...ITenantResolveContrib) *TenantResolveOption {
	return &TenantResolveOption{
		Resolvers: c,
	}
}

func (opt *TenantResolveOption) AppendContribs(c ...ITenantResolveContrib) {
	opt.Resolvers = append(opt.Resolvers, c...)
}

func (opt *TenantResolveOption) RemoveContribs(c ...ITenantResolveContrib) {
	var r []ITenantResolveContrib
	for _, resolver := range opt.Resolvers {
		if !contains(c, resolver) {
			r = append(r, resolver)
		}
	}
	opt.Resolvers = r
}

func contains(a []ITenantResolveContrib, b ITenantResolveContrib) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == b {
			return true
		}
	}
	return false
}
