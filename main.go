package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Api struct {
	addr string
}

var Users = []User{}

func (a *Api) getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	// encoding users to slices
	err := json.NewEncoder(w).Encode(Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func validateUser(u User) error {
	if u.Name == "" {
		return errors.New("name cannot be blank")
	}
	if u.Age < 1 {
		return errors.New("age must be greater than 1")
	}

	for _, user := range Users {
		if user.Name == u.Name {
			return errors.New("name already exists")
		}
	}
	return nil
}

func (a *Api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload User
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := User{
		Name: payload.Name,
		Age:  payload.Age,
	}

	if validateUser(u) == nil {
		Users = append(Users, u)
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	api_ := &Api{addr: ":8082"}
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    api_.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /users", api_.getUserHandler)
	mux.HandleFunc("POST /users", api_.createUserHandler)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
