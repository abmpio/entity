package iris

import (
	"errors"

	"github.com/abmpio/entity/tenancy"
	"github.com/abmpio/entity/tenancy/http"
	"github.com/abmpio/webserver/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type ErrorFormatter func(context iris.Context, err error)

var (
	DefaultErrorFormatter ErrorFormatter = func(context iris.Context, err error) {
		if errors.Is(err, tenancy.ErrTenantNotFound) {
			context.StopWithError(404, err)
		} else {
			context.StopWithError(500, err)
		}
	}
)

type option struct {
	hmtOpt  *http.WebMultiTenancyOption
	ef      ErrorFormatter
	resolve []tenancy.ResolveOption
}

type Option func(*option)

func WithMultiTenancyOption(opt *http.WebMultiTenancyOption) Option {
	return func(o *option) {
		o.hmtOpt = opt
	}
}

func WithErrorFormatter(e ErrorFormatter) Option {
	return func(o *option) {
		o.ef = e
	}
}

func WithResolveOption(opt ...tenancy.ResolveOption) Option {
	return func(o *option) {
		o.resolve = opt
	}
}

func MultiTenancy(ts tenancy.ITenantStore, options ...Option) iris.Handler {
	opt := &option{
		hmtOpt:  http.NewDefaultWebMultiTenancyOption(),
		ef:      DefaultErrorFormatter,
		resolve: nil,
	}
	for _, o := range options {
		o(opt)
	}
	return func(context *context.Context) {
		var trOpt []tenancy.ResolveOption
		df := []tenancy.ITenantResolveContrib{
			http.NewHeaderTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request()),
			http.NewCookieTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request()),
			http.NewQueryTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request()),
			http.NewFormTenantResolveContrib(opt.hmtOpt.TenantKey, context.Request()),
		}
		if opt.hmtOpt.DomainFormat != "" {
			df = append(df, http.NewDomainTenantResolveContrib(opt.hmtOpt.DomainFormat, context.Request()))
		}
		trOpt = append(trOpt, tenancy.AppendContribs(df...))
		trOpt = append(trOpt, opt.resolve...)

		// get tenant config
		tenantConfigProvider := tenancy.NewDefaultTenantConfigProvider(tenancy.NewDefaultTenantResolver(trOpt...), ts)
		tenantConfig, err := tenantConfigProvider.Get(context)
		if err != nil {
			opt.ef(context, err)
			controller.HandleErrorForbidden(context)
			return
		}
		if tenantConfig != nil && len(tenantConfig.Id) > 0 {
			context.Values().Set(opt.hmtOpt.TenantKey, tenantConfig.Id)
			context.Values().Set("tenant", tenantConfig)
		}
		// next
		context.Next()
	}
}
