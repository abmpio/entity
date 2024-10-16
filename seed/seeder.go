package seed

import (
	"context"
)

type Seeder interface {
	Seed(ctx context.Context, option ...Option) error
}

var _ Seeder = (*DefaultSeeder)(nil)

type DefaultSeeder struct {
	contrib []Contrib
}

func NewDefaultSeeder(contrib ...Contrib) *DefaultSeeder {
	return &DefaultSeeder{
		contrib: contrib,
	}
}

func (d *DefaultSeeder) Seed(ctx context.Context, options ...Option) error {
	opt := NewOption()
	for _, option := range options {
		option(opt)
	}
	for _, tenantId := range opt.TenantIds {
		// change to next tenant
		ctxWithTenant := newCurrentTenant(ctx, tenantId, "")
		if err := d._seed(ctxWithTenant, tenantId, opt); err != nil {
			return err
		}
	}
	return nil
}

func (d *DefaultSeeder) _seed(ctx context.Context, tenantId string, opt *SeedOption) error {
	sCtx := NewSeedContext(tenantId, opt.Extra)
	// create seeder
	for _, contributor := range d.contrib {
		if err := contributor.Seed(ctx, sCtx); err != nil {
			return err
		}
	}
	return nil
}
