package tenancy

type IMultiTenancy interface {
	SetTenantId(tenantId string)
	GetTenantId() string

	MustHaveTenantId() bool
}

// MultiTenancy entity
type MultiTenancy struct {
	TenantId string `json:"tenantId" bson:"tenantId"`
}

type MustHaveMultiTenancy struct {
	MultiTenancy `bson:",inline"`
}

type MayHaveMultiTenancy struct {
	MultiTenancy `bson:",inline"`
}

var _ IMultiTenancy = (*MultiTenancy)(nil)
var _ IMultiTenancy = (*MustHaveMultiTenancy)(nil)
var _ IMultiTenancy = (*MayHaveMultiTenancy)(nil)

// #region IMultiTenancy Members

func (t *MultiTenancy) SetTenantId(tenantId string) {
	t.TenantId = tenantId
}

func (t *MultiTenancy) GetTenantId() string {
	return t.TenantId
}

func (t *MultiTenancy) MustHaveTenantId() bool {
	return false
}

// #endregion

func (t *MustHaveMultiTenancy) MustHaveTenantId() bool {
	return true
}

func CheckEntityIsIMultiTenancy(entityValue interface{}) IMultiTenancy {
	v, ok := entityValue.(IMultiTenancy)
	if !ok {
		return nil
	}
	return v
}
