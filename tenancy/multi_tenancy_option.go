package tenancy

type DatabaseStyleType int32

const (
	// 单一数据库，所有租户一个数据库模式
	Single DatabaseStyleType = 1 << 0
)

type MultiTenancyOption struct {
	IsEnabled     bool
	DatabaseStyle DatabaseStyleType
}

type option func(tenancyOption *MultiTenancyOption)

// WithEnabled enable status
func WithEnabled(isEnabled bool) option {
	return func(tenancyOption *MultiTenancyOption) {
		tenancyOption.IsEnabled = isEnabled
	}
}

// WithDatabaseStyle database style, support Single/PerTenant/Multi
func WithDatabaseStyle(databaseStyle DatabaseStyleType) option {
	return func(tenancyOption *MultiTenancyOption) {
		tenancyOption.DatabaseStyle = databaseStyle
	}
}

// 创建MultiTenancyOption对象实例
func NewMultiTenancyOption(opts ...option) *MultiTenancyOption {
	option := MultiTenancyOption{}
	for _, opt := range opts {
		opt(&option)
	}
	return &option
}

// 构建默认的多租户参数
func DefaultMultiTenancyOption() *MultiTenancyOption {
	return NewMultiTenancyOption(WithEnabled(true), WithDatabaseStyle(Single))
}
