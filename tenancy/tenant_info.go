package tenancy

type ITenantInfo interface {
	GetTenantId() string
	// get admin user id for tenant
	GetAdminUserId() string
}
