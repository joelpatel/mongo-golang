package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joelpatel/mongo-golang/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017/")
	if err != nil {
		panic(err)
	}

	return session
}

func NewUserController() *UserController {
	return &UserController{getSession()}
}

func (uc *UserController) GetUser(res_w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	if !bson.IsObjectIdHex(id) {
		res_w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	oid := bson.ObjectIdHex(id)
	user := models.User{}

	if err := uc.session.DB("mongo-golang").C("users").FindId(oid).One(&user); err != nil { // mongo-golang > database, users > collection
		res_w.WriteHeader((http.StatusNotFound))
		return
	}

	user_json, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	res_w.Header().Set("Content-Type", "application/json")
	res_w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(res_w, "%s\n", user_json)
}

func (uc *UserController) CreateUser(res_w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := models.User{}

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res_w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	user.ID = bson.NewObjectId()

	uc.session.DB("mongo-golang").C("users").Insert(user)

	user_json, err := json.Marshal(user)

	if err != nil {
		res_w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	res_w.Header().Set("Content-Type", "application/json")
	res_w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(res_w, "%s\n", user_json)
}

func (uc *UserController) DeleteUser(res_w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	if !bson.IsObjectIdHex(id) {
		res_w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	oid := bson.ObjectIdHex(id)

	if err := uc.session.DB("mongo-golang").C("users").RemoveId(oid); err != nil {
		res_w.WriteHeader(http.StatusNotFound) // 404
		return
	}

	res_w.WriteHeader(http.StatusOK)
	fmt.Fprintf(res_w, "Deleted user { %v }.\n", oid)
}

// func UpdateUser() {

// }
