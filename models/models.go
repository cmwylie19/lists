package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type List struct {
	ID            primitive.ObjectID `json:"_id,onitempty" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name,omitempty"`
	Owner         string             `json:"owner" bson:"owner,omitempty"`
	Collaborators []string           `json:"collaborators,omitempty" bson:"collaborators,omitempty"`
}
