package entity

type TenantInfo interface {
	GetId() string
	GetName() string
}

// #region TenantInfo Members

func (b *BasicTenantInfo) GetId() string {
	return b.Id
}

func (b *BasicTenantInfo) GetName() string {
	return b.Name
}

// #endregion

type BasicTenantInfo struct {
	Id   string
	Name string
}

func NewBasicTenantInfo(id string, name string) *BasicTenantInfo {
	return &BasicTenantInfo{Id: id, Name: name}
}
