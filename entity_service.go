package entity

import (
	"github.com/abmpio/mongodbr"
	mongodbBuilder "github.com/abmpio/mongodbr/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IEntityService[T mongodbr.IEntity] interface {
	GetRepository() mongodbr.IRepository

	FindAll(opts ...mongodbr.MongodbrFindOption) ([]*T, error)
	FindList(filter interface{}, opts ...mongodbr.MongodbrFindOption) (list []*T, err error)
	FindListByIdList(idList []primitive.ObjectID, opts ...mongodbr.MongodbrFindOption) (list []*T, err error)
	Count(filter interface{}, opts ...mongodbr.MongodbrCountOption) (count int64, err error)
	FindById(id primitive.ObjectID, opts ...mongodbr.MongodbrFindOneOption) (*T, error)
	FindOne(filter interface{}, opts ...mongodbr.MongodbrFindOneOption) (*T, error)

	Create(interface{}, ...mongodbr.MongodbrInsertOneOption) (*T, error)
	Delete(primitive.ObjectID, ...mongodbr.MongodbrDeleteOption) error
	DeleteMany(interface{}, ...mongodbr.MongodbrDeleteOption) (*mongo.DeleteResult, error)
	DeleteManyByIdList(idList []primitive.ObjectID, opts ...mongodbr.MongodbrDeleteOption) (*mongo.DeleteResult, error)
	UpdateFields(id primitive.ObjectID, update map[string]interface{}, opts ...mongodbr.MongodbrFindOneAndUpdateOption) error
}

type EntityService[T mongodbr.IEntity] struct {
	repository mongodbr.IRepository
}

func NewEntityService[T mongodbr.IEntity](repository mongodbr.IRepository) IEntityService[T] {
	return &EntityService[T]{
		repository: repository,
	}
}

// #region IEntityService[mongodbr.IEntity]

func (s *EntityService[T]) GetRepository() mongodbr.IRepository {
	return s.repository
}

func (s *EntityService[T]) FindAll(opts ...mongodbr.MongodbrFindOption) (list []*T, err error) {
	res, err := mongodbr.FindAllT[T](s.repository, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *EntityService[T]) FindList(filter interface{}, opts ...mongodbr.MongodbrFindOption) (list []*T, err error) {
	return mongodbr.FindTByFilter[T](s.repository, filter, opts...)
}

func (s *EntityService[T]) FindListByIdList(idList []primitive.ObjectID, opts ...mongodbr.MongodbrFindOption) (list []*T, err error) {
	return mongodbr.FindTListByObjectIdList[T](s.repository, idList, opts...)
}

func (s *EntityService[T]) Count(filter interface{}, opts ...mongodbr.MongodbrCountOption) (count int64, err error) {
	return s.repository.CountByFilter(filter, opts...)
}

func (s *EntityService[T]) FindById(id primitive.ObjectID, opts ...mongodbr.MongodbrFindOneOption) (*T, error) {
	return mongodbr.FindTByObjectId[T](s.repository, id, opts...)
}

func (s *EntityService[T]) FindOne(filter interface{}, opts ...mongodbr.MongodbrFindOneOption) (*T, error) {
	return mongodbr.FindOneTByFilter[T](s.repository, filter, opts...)
}

func (s *EntityService[T]) Create(item interface{}, opts ...mongodbr.MongodbrInsertOneOption) (*T, error) {
	oid, err := s.repository.Create(item, opts...)

	if err != nil {
		return nil, err
	}

	// extract ctx
	insertOneOptions := &mongodbr.MongodbrInsertOneOptions{}
	for _, eachOpt := range opts {
		eachOpt(insertOneOptions)
	}
	dbItem, err := s.FindById(oid, func(mfoo *mongodbr.MongodbrFindOneOptions) {
		mfoo.WithCtx = insertOneOptions.WithCtx
	})
	if err != nil {
		return nil, err
	}
	return dbItem, nil
}

func (s *EntityService[T]) Delete(id primitive.ObjectID, opts ...mongodbr.MongodbrDeleteOption) error {
	_, err := s.repository.DeleteOne(id, opts...)
	if err != nil {
		return err
	}
	return nil
}

func (s *EntityService[T]) DeleteMany(filter interface{}, opts ...mongodbr.MongodbrDeleteOption) (*mongo.DeleteResult, error) {
	return s.repository.DeleteMany(filter, opts...)
}

func (service *EntityService[T]) DeleteManyByIdList(idList []primitive.ObjectID, opts ...mongodbr.MongodbrDeleteOption) (*mongo.DeleteResult, error) {
	if len(idList) <= 0 {
		return &mongo.DeleteResult{
			DeletedCount: 0,
		}, nil
	}
	filter := bson.M{
		"_id": bson.M{"$in": idList},
	}
	return service.DeleteMany(filter, opts...)
}

// update fields value
func (s *EntityService[T]) UpdateFields(id primitive.ObjectID, update map[string]interface{}, opts ...mongodbr.MongodbrFindOneAndUpdateOption) error {
	value := mongodbBuilder.NewBsonBuilder().NewOrUpdateSet(update).ToValue()
	return s.repository.FindOneAndUpdateWithId(id, value, opts...)
}

// #endregion
