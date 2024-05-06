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

	FindAll() ([]*T, error)
	FindList(filter interface{}, opts ...mongodbr.FindOption) (list []*T, err error)
	Count(filter interface{}) (count int64, err error)
	FindById(id primitive.ObjectID) (*T, error)
	FindOne(filter interface{}) (*T, error)

	Create(interface{}) (*T, error)
	Delete(primitive.ObjectID) error
	DeleteMany(interface{}) (*mongo.DeleteResult, error)
	DeleteManyByIdList(idList []primitive.ObjectID) (*mongo.DeleteResult, error)
	UpdateFields(id primitive.ObjectID, update map[string]interface{}) error
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

func (s *EntityService[T]) FindAll() (list []*T, err error) {
	res := s.repository.FindAll()
	err = res.All(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *EntityService[T]) FindList(filter interface{}, opts ...mongodbr.FindOption) (list []*T, err error) {
	return mongodbr.FindTByFilter[T](s.repository, filter, opts...)
}

func (s *EntityService[T]) Count(filter interface{}) (count int64, err error) {
	return s.repository.CountByFilter(filter)
}

func (s *EntityService[T]) FindById(id primitive.ObjectID) (*T, error) {
	return mongodbr.FindTByObjectId[T](s.repository, id)
}

func (s *EntityService[T]) FindOne(filter interface{}) (*T, error) {
	return mongodbr.FindOneTByFilter[T](s.repository, filter)
}

func (s *EntityService[T]) Create(item interface{}) (*T, error) {
	oid, err := s.repository.Create(item)

	if err != nil {
		return nil, err
	}
	dbItem, err := s.FindById(oid)
	if err != nil {
		return nil, err
	}
	return dbItem, nil
}

func (s *EntityService[T]) Delete(id primitive.ObjectID) error {
	_, err := s.repository.DeleteOne(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *EntityService[T]) DeleteMany(filter interface{}) (*mongo.DeleteResult, error) {
	return s.repository.DeleteMany(filter)
}

func (service *EntityService[T]) DeleteManyByIdList(idList []primitive.ObjectID) (*mongo.DeleteResult, error) {
	if len(idList) <= 0 {
		return &mongo.DeleteResult{
			DeletedCount: 0,
		}, nil
	}
	filter := bson.M{
		"_id": bson.M{"$in": idList},
	}
	return service.DeleteMany(filter)
}

// update fields value
func (s *EntityService[T]) UpdateFields(id primitive.ObjectID, update map[string]interface{}) error {
	value := mongodbBuilder.NewBsonBuilder().NewOrUpdateSet(update).ToValue()
	return s.repository.FindOneAndUpdateWithId(id, value)
}

// #endregion
