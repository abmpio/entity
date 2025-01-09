package mongodb

import (
	"fmt"

	"github.com/abmpio/libx/reflector"
	"github.com/abmpio/libx/str"
	"github.com/abmpio/mongodbr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// key为实体对象的fullname,value为collection name
	_registedEntityCollection map[string]string = make(map[string]string)
)

func RegistCollectionNameByType(v interface{}, collectionName string) {
	_registedEntityCollection[reflector.GetFullName(v)] = collectionName
}

func RegistEntityCollectionName[T mongodbr.IEntity](collectionName string) {
	RegistCollectionNameByType(new(T), collectionName)
}

func GetCollectionName(v interface{}) string {
	key := reflector.GetFullName(v)
	collectionName, ok := _registedEntityCollection[key]
	if !ok {
		return str.ToSnake(reflector.GetName(v))
	}
	return collectionName
}

// regist Repository create option
func RegistEntityRepositoryOption[T mongodbr.IEntity](clientKey string, databaseName string, opts ...mongodbr.RepositoryOption) {
	d := GetDatabase(clientKey, databaseName)
	collectionName := GetCollectionName(new(T))
	d._entityRepositoryOptionMap[collectionName] = createEntityRepositoryOption[T](opts...)
	d._repositoryMapping[collectionName] = d.ensureCreateRepository(collectionName, opts...)
}

// regist Repository create option with collection name prefiex
func RegistEntityRepositoryOptionWithPrefix[T mongodbr.IEntity](clientKey string, databaseName string, collectionPrefix string, opts ...mongodbr.RepositoryOption) {
	d := GetDatabase(clientKey, databaseName)
	collectionName := GetCollectionName(new(T))
	if len(collectionPrefix) > 0 {
		collectionName = fmt.Sprintf("%s%s", collectionPrefix, collectionName)
		RegistEntityCollectionName[T](collectionName)
	}
	d._entityRepositoryOptionMap[collectionName] = createEntityRepositoryOption[T](opts...)
	d._repositoryMapping[collectionName] = d.ensureCreateRepository(collectionName, opts...)
}

func createEntityRepositoryOption[T mongodbr.IEntity](opts ...mongodbr.RepositoryOption) []mongodbr.RepositoryOption {
	if len(opts) <= 0 {
		opts = append(opts, mongodbr.WithCreateItemFunc(func() interface{} {
			return new(T)
		}))
		opts = append(opts, mongodbr.WithDefaultSort(func(fo *options.FindOptions) *options.FindOptions {
			return fo.SetSort(bson.D{{Key: "_id", Value: -1}})
		}))
	}
	return opts
}
