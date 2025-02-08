package tenancy

// 租户配置信息
type TenantConfig struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	AdminUserId string `json:"adminUserId"`
}

// new TenantConfig
func NewTenantConfig(id, name, adminUserId string) *TenantConfig {
	return &TenantConfig{
		Id:          id,
		Name:        name,
		AdminUserId: adminUserId,
	}
}

var _ ITenantInfo = (*TenantConfig)(nil)

func (c *TenantConfig) GetTenantId() string {
	return c.Id
}

// get admin user id for tenant
func (c *TenantConfig) GetAdminUserId() string {
	return c.AdminUserId
}
