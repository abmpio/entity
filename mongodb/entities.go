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
	// 可以用于设置每个collection的参数，key为实体对象的fullname,value为一个func
	_registedEntityCollectionOptionsFunc map[string][]func(*options.CollectionOptions) = make(map[string][]func(*options.CollectionOptions))
)

// set collection name by type
func RegistCollectionNameByType(v interface{}, collectionName string) {
	_registedEntityCollection[reflector.GetFullName(v)] = collectionName
}

// set collection name by type
func RegistEntityCollectionName[T mongodbr.IEntity](collectionName string) {
	RegistCollectionNameByType(new(T), collectionName)
}

// set collection options by type
func RegistEntityCollectionOptionsByType(v interface{}, opts ...func(*options.CollectionOptions)) {
	if len(opts) <= 0 {
		return
	}
	key := reflector.GetFullName(v)
	list, ok := _registedEntityCollectionOptionsFunc[key]
	if !ok {
		list = make([]func(*options.CollectionOptions), 0)
	}
	list = append(list, opts...)
	_registedEntityCollectionOptionsFunc[key] = list
}

// set collection options by type
func RegistEntityCollectionOptions[T mongodbr.IEntity](opts ...func(collOptions *options.CollectionOptions)) {
	RegistEntityCollectionOptionsByType(new(T), opts...)
}

// apply entity collection options used registed collection options func
func applyCollectionOptions(entityV interface{}, collectionOptions *options.CollectionOptions) {
	key := reflector.GetFullName(entityV)
	optList, ok := _registedEntityCollectionOptionsFunc[key]
	if !ok || len(optList) <= 0 {
		// not registed collection option func for this entity type
		return
	}
	for _, eachOpt := range optList {
		eachOpt(collectionOptions)
	}
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
