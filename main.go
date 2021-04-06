package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func strToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

type Fruit struct {
	Apples  int `json:"apples,omitempty"`
	Oranges int `json:"oranges,omitempty"`
	Total   int `json:"total,omitempty"`
}

var fruits []Fruit

func createFruits(w http.ResponseWriter, r *http.Request) {
	var fruit Fruit
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &fruit); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	fruit.Total = fruit.Apples + fruit.Oranges
	fruits = append(fruits, fruit)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fruit)
	if err != nil {
		panic(err)
	}
}

func getFruit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	idInt, _ := strToInt(idStr)
	if idInt + 1 > len(fruits) {
		w.WriteHeader(http.StatusNotFound)
		err := errors.New("Index " + idStr + " not found!")
		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	fruit := fruits[idInt]
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(fruit)
	if err != nil {
		panic(err)
	}
}

func listFruits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(fruits)
	if err != nil {
		panic(err)
	}
}

func modifyFruits(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	var fruit Fruit
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &fruit); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	idInt, _ := strToInt(idStr)
	if idInt + 1 > len(fruits) {
		w.WriteHeader(http.StatusNotFound)
		err := errors.New("Index " + idStr + " not found!")
		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	var newFruit Fruit
	curFruit := fruits[idInt]
	newFruit.Apples = curFruit.Apples - fruit.Apples
	newFruit.Oranges = curFruit.Oranges - fruit.Oranges
	newFruit.Total = newFruit.Apples + newFruit.Oranges
	fruits[idInt] = newFruit
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newFruit)
	if err != nil {
		panic(err)
	}
}

func gapFruits(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStrA := params["idA"]
	idIntA, _ := strToInt(idStrA)
	if idIntA + 1 > len(fruits) {
		w.WriteHeader(http.StatusNotFound)
		err := errors.New("Index " + idStrA + " not found!")
		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	fruitA := fruits[idIntA]
	idStrB := params["idB"]
	idIntB, _ := strToInt(idStrB)
	if idIntB + 1 > len(fruits) {
		w.WriteHeader(http.StatusNotFound)
		err := errors.New("Index " + idStrB + " not found!")
		if err = json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	fruitB := fruits[idIntB]
	var newFruit Fruit
	newFruit.Apples = fruitA.Apples - fruitB.Apples
	newFruit.Oranges = fruitA.Oranges - fruitB.Oranges
	newFruit.Total = newFruit.Apples + newFruit.Oranges
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(newFruit)
	if err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/fruits", createFruits).Methods("POST")
	router.HandleFunc("/fruits/{id}", getFruit).Methods("GET")
	router.HandleFunc("/fruits", listFruits).Methods("GET")
	router.HandleFunc("/fruits/{id}", modifyFruits).Methods("PATCH")
	router.HandleFunc("/fruits/{idA}/{idB}", gapFruits).Methods("GET")
	http.ListenAndServe(":8080", router)
}

/*
5 PostMan Requests :
POST - create
				Headers 	-	Content-Type : Application/JSON
 				Body 		-	{ "apples": 4, "oranges": 8 }
GET - get
				Headers 	- 	Content-Type :  Application/JSON
GET - list
				Headers 	-	Content-Type :  Application/JSON
PATCH - modify
				Headers 	-	Content-Type :  Application/JSON
GET - gap
				Headers 	-	Content-Type :  Application/JSON
*/