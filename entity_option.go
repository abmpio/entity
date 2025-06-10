package entity

import (
	"github.com/abmpio/abmp/pkg/utils/lang"
	"github.com/abmpio/libx/lang/tuple"
	"github.com/abmpio/mongodbr"
	mongodbBuilder "github.com/abmpio/mongodbr/builder"
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

type BsonBuilderOption func(*mongodbBuilder.BsonBuilder)

// new mongodbBuilder.BsonBuilder with option
func NewBsonBuilderWithOption(opts ...BsonBuilderOption) *mongodbBuilder.BsonBuilder {
	bsonBuilder := &mongodbBuilder.BsonBuilder{}
	for _, eachOpt := range opts {
		eachOpt(bsonBuilder)
	}
	return bsonBuilder
}

// set lastModificationTime field value to now
func BsonBuilderOptionWithLastModificationTime() BsonBuilderOption {
	return func(bb *mongodbBuilder.BsonBuilder) {
		bb.AppendSetField(tuple.New2[string, interface{}]("lastModificationTime", lang.NowToPtr()))
	}
}

// set lastModifierId field value to now
func BsonBuilderOptionWithLastModifierId(userId string) BsonBuilderOption {
	return func(bb *mongodbBuilder.BsonBuilder) {
		bb.AppendSetField(tuple.New2[string, interface{}]("lastModifierId", userId))
	}
}
