package dao

import (
	"log"

	. "mongo-to-do/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TasksDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "tasks"
)

// Establish a connection to database
func (t *TasksDAO) Connect() {
	session, err := mgo.Dial(t.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(t.Database)
}

// Find list of tasks
func (t *TasksDAO) FindAll() ([]Task, error) {
	var tasks []Task
	err := db.C(COLLECTION).Find(bson.M{}).All(&tasks)
	return tasks, err
}

// Find a task bt its name
func (t *TasksDAO) FindByName(name string) (Task, error) {
	var task Task
	err := db.C(COLLECTION).Find(bson.M{"name": name}).One(&task)
	return task, err
}

// Find a task by its id
func (t *TasksDAO) FindById(id string) (Task, error) {
	var task Task
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&task)
	return task, err
}

// Insert a task into database
func (t *TasksDAO) Insert(task Task) error {
	err := db.C(COLLECTION).Insert(&task)
	return err
}

// Delete an existing task
func (t *TasksDAO) Delete(task Task) error {
	err := db.C(COLLECTION).Remove(&task)
	return err
}

// Update an existing task
func (t *TasksDAO) Update(task Task) error {
	err := db.C(COLLECTION).UpdateId(task.ID, &task)
	return err
}
