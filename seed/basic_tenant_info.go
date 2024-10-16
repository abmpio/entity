package seed

type tenantInfo interface {
	GetId() string
	GetName() string
}

type basicTenantInfo struct {
	Id   string
	Name string
}

func (b *basicTenantInfo) GetId() string {
	return b.Id
}

func (b *basicTenantInfo) GetName() string {
	return b.Name
}

func NewBasicTenantInfo(id string, name string) *basicTenantInfo {
	return &basicTenantInfo{Id: id, Name: name}
}
