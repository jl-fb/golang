package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jl-fb/crud/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

var pessoas []model.Person
var person model.Person

func getID(r *http.Request) primitive.ObjectID {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	return id
}

func createPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("testMongo").Collection("pessoas")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func getPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("testMongo").Collection("pessoas")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `}"`))
		return
	}
	defer cursor.Close(ctx)
	cursor.All(ctx, &pessoas)
	// for cursor.Next(ctx) {
	// 	var person Person
	// 	cursor.Decode(&person)
	// 	people = append(people, person)
	// }
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(w).Encode(pessoas)
}

func getPersonEndpoint(resp http.ResponseWriter, req *http.Request) {
	//resp.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(req)
	// id, _ := primitive.ObjectIDFromHex(params["id"])
	id := getID(req)
	collection := client.Database("testMongo").Collection("pessoas")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-time.After(6 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

	err := collection.FindOne(ctx, model.Person{ID: id}).Decode(&person)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(resp).Encode(person)
}

func deletePersonEndPoint(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(r)
	// id, _ := primitive.ObjectIDFromHex(params["id"])
	id := getID(r)
	collection := client.Database("testMongo").Collection("pessoas")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	result, err := collection.DeleteOne(ctx, model.Person{ID: id})
	if err != nil {
		fmt.Printf("Erro ao deletar usuário: %v", err)
	}
	json.NewEncoder(w).Encode(result)
}

func updatePersonEndPoint(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	//id, _ := primitive.ObjectIDFromHex(params["id"])
	id := getID(r)
	collection := client.Database("testMongo").Collection("pessoas")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	filter := bson.M{"_id": id}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("[crud] Erro ao ler conteudo do servidor. Erro: ", err.Error())
	}
	defer r.Body.Close()

	// Pedgando dados do body da requisão que serão usados para atualizar o DB
	err = json.Unmarshal(body, &person)
	if err != nil {
		fmt.Println("[crud] Erro ao converter o retorno json do servidor. Erro: ", err.Error())
	}

	// Passando os dados do body para poder ser tratado pelo GO e mandando para o mongo
	update := bson.M{
		"$set": bson.M{
			"firstname": person.Firstname,
			"lastname":  person.Lastname,
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("Erro ao atualizar usuário: %v", err)
	}
	json.NewEncoder(w).Encode(result)

}
