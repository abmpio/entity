package entity

import (
	"github.com/abmpio/mongodbr"
)

func EntityWithCreatorUserId(userId string) mongodbr.EntityOption {
	return func(e mongodbr.IEntity) {
		entityWithUser, ok := e.(IEntityWithUser)
		if !ok || entityWithUser == nil {
			return
		}
		entityWithUser.SetUserCreator(userId)
	}
}

func ApplyEntityOption(e mongodbr.IEntity, opts ...mongodbr.EntityOption) {
	if e == nil || len(opts) <= 0 {
		return
	}
	for _, eachOpt := range opts {
		eachOpt(e)
	}
}
