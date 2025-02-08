package tenancy

import "context"

type ITenantResolver interface {
	Resolve(ctx context.Context) (TenantResolveResult, error)
}

type DefaultTenantResolver struct {
	//options
	o *TenantResolveOption
}

func NewDefaultTenantResolver(opt ...ResolveOption) ITenantResolver {
	o := NewTenantResolveOption(&ContextContrib{})
	for _, resolveOption := range opt {
		resolveOption(o)
	}
	return &DefaultTenantResolver{
		o: o,
	}
}

// #region ITenantResolver Members

func (d *DefaultTenantResolver) Resolve(ctx context.Context) (TenantResolveResult, error) {
	res := TenantResolveResult{}
	for _, resolver := range d.o.Resolvers {
		tenantId, err := resolver.Resolve(ctx)
		if err != nil {
			return res, err
		}
		res.TenantId = tenantId
		res.AppliedResolvers = append(res.AppliedResolvers, resolver.Name())
		if tenantId != "" {
			break
		}
	}
	return res, nil
}

// #endregion
