package entity

import "go.mongodb.org/mongo-driver/v2/bson"

type BatchRequestPayload struct {
	Ids []bson.ObjectID `form:"ids" json:"ids"`
}
