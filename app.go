package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	. "mongo-to-do/config"
	. "mongo-to-do/dao"
	. "mongo-to-do/models"

	"github.com/gorilla/mux"
)

var config = Config{}
var dao = TasksDAO{}

// GET list of tasks
func AllTasksEndPoint(w http.ResponseWriter, r *http.Request) {
	tasks, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tasks)
}

// GET a task by its name
func FindTaskByNameEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	task, err := dao.FindByName(params["name"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task name")
	}
	respondWithJson(w, http.StatusOK, task)
}

// GET a task by its ID
func FindTaskByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	task, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}
	respondWithJson(w, http.StatusOK, task)
}

// POST a new task
func CreateTaskEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	task.ID = bson.NewObjectId()
	if err := dao.Insert(task); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, task)
}

// PUT update an existing task
func UpdateTaskEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(task); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing task
func DeleteTaskEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(task); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tasks", AllTasksEndPoint).Methods("GET")
	r.HandleFunc("/tasks", CreateTaskEndPoint).Methods("POST")
	r.HandleFunc("/tasks", UpdateTaskEndPoint).Methods("PUT")
	r.HandleFunc("/tasks", DeleteTaskEndPoint).Methods("DELETE")
	r.HandleFunc("/tasks/id/{id}", FindTaskByIDEndpoint).Methods("GET")
	r.HandleFunc("/tasks/name/{name}", FindTaskByNameEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
