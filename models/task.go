package models

import "gopkg.in/mgo.v2/bson"

// Represents a task, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type Task struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
}
