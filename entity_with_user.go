package entity

type IEntityWithUser interface {
	SetUserCreator(userId string)
	GetCreatorId() string
}

// type EntityWithUser struct {
// 	mongodbr.AuditedEntity `bson:",inline"`
// }

// // #region IEntityWithUser Members

// func (p *EntityWithUser) SetUserCreator(userId string) {
// 	p.CreatorId = userId
// }

// func (p *EntityWithUser) GetCreatorId() string {
// 	return p.CreatorId
// }

// #endregion

func CheckEntityIsIEntityWithUser(entityValue interface{}) IEntityWithUser {
	v, ok := entityValue.(IEntityWithUser)
	if !ok {
		return nil
	}
	return v
}
