package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/cmwylie19/lists/helper"
	"github.com/cmwylie19/lists/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Connection mongoDB with helper class
var collection = helper.ConnectDB()

func getLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var lists []models.List

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		helper.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var list models.List

		err := cur.Decode(&list)
		if err != nil {
			log.Fatal(err)
		}
		lists = append(lists, list)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(lists)
}

func getList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var list models.List
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&list)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func createList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var list models.List

	// decode request
	_ = json.NewDecoder(r.Body).Decode(&list)
	result, err := collection.InsertOne(context.TODO(), list)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}
func updateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var list models.List
	filter := bson.M{"_id", id}
	_ = json.NewDecoder(r.Body).Decode(&list)
	update := bson.D{
		{"$set", bson.D{
			{"Collaborators", list.Collaborators},
			{"Name": list.Name},
			{"Owner": list.Owner},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&list)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	list.ID = id
	json.NewEncoder(w).Encode(list)

}

func deleteList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}

type Remote struct {
	XFF string `json:"x-forwarded-for"`
}

func getRemote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	remote := Remote{XFF: r.Header.Get("x-forwarded-for")}
	result, err := json.Marshal(remote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func main() {
	//Init Router
	r := mux.NewRouter()

	r.HandleFunc("/api/lists", getLists).Methods("GET")
	r.HandleFunc("/api/lists/{id}", getList).Methods("GET")
	r.HandleFunc("/api/lists", createList).Methods("POST")
	r.HandleFunc("/api/lists/{id}", updateList).Methods("PUT")
	r.HandleFunc("/api/lists/{id}", deleteList).Methods("DELETE")
	r.HandleFunc("/remote", getRemote).Methods("GET")

	config := helper.GetConfiguration()

	log.Fatal(http.ListenAndServe(config.Port, r))

}
