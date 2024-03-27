package mongodb

import (
	"fmt"

	"github.com/abmpio/mongodbr"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	_registedDatabase map[string]*Database = make(map[string]*Database)
)

type Database struct {
	MongoClientKey string
	DatabaseName   string
	_db            *mongo.Database
	// key: CollectionName
	// value: RepositoryOption
	_entityRepositoryOptionMap map[string][]mongodbr.RepositoryOption
	_repositoryMapping         map[string]mongodbr.IRepository
}

func NewDatabase(databaseName string) *Database {
	d := &Database{
		DatabaseName: databaseName,
	}
	d._db = mongodbr.DefaultClient().Database(databaseName)
	d._entityRepositoryOptionMap = make(map[string][]mongodbr.RepositoryOption)
	d._repositoryMapping = make(map[string]mongodbr.IRepository)

	return d
}

func NewDatabaseWithClientKey(clientKey, databaseName string) *Database {
	if len(clientKey) <= 0 {
		return NewDatabase(databaseName)
	}
	d := &Database{
		MongoClientKey: clientKey,
		DatabaseName:   databaseName,
	}
	d._db = mongodbr.GetDatabaseByKey(clientKey, databaseName)
	if d._db == nil {
		panic(fmt.Errorf("没有配置好key为 %s 的mongodb的连接字符串", clientKey))
	}
	d._entityRepositoryOptionMap = make(map[string][]mongodbr.RepositoryOption)
	d._repositoryMapping = make(map[string]mongodbr.IRepository)

	return d
}

func ensureDatabaseRegisted(clientKey, databaseName string) *Database {
	key := fmt.Sprintf("%s_%s", clientKey, databaseName)
	d, ok := _registedDatabase[key]
	if ok {
		return d
	}
	d = NewDatabaseWithClientKey(clientKey, databaseName)
	_registedDatabase[key] = d
	return d
}

// 获取指定db的IRepository接口
func GetDatabase(clientKey, databaseName string) *Database {
	return ensureDatabaseRegisted(clientKey, databaseName)
}

func (d *Database) GetDatabase() *mongo.Database {
	return d._db
}

func (d *Database) GetEntityRepositoryOptionMap() map[string][]mongodbr.RepositoryOption {
	return d._entityRepositoryOptionMap
}

func (d *Database) GetRepository(modelInstance interface{}) mongodbr.IRepository {
	collectionName := GetCollectionName(modelInstance)
	instance, ok := d._repositoryMapping[collectionName]
	if ok {
		return instance
	}
	return nil
}

func (d *Database) ensureCreateRepository(collectionName string, opts ...mongodbr.RepositoryOption) mongodbr.IRepository {
	if len(d._entityRepositoryOptionMap) > 0 {
		registedOpts, ok := d._entityRepositoryOptionMap[collectionName]
		if ok && len(registedOpts) > 0 {
			opts = append(opts, registedOpts...)
		}
	}
	repository, err := mongodbr.NewRepositoryBase(func() *mongo.Collection {
		return d._db.Collection(collectionName)
	}, opts...)
	if err != nil {
		panic(err)
	}
	if repository == nil {
		panic("cannot create repository for collection" + collectionName)
	}
	return repository
}
