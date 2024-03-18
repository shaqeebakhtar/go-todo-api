package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos = []Todo{
	{
		Id:    "1",
		Title: "First Todo",
		Done:  false,
	},
	{
		Id:    "2",
		Title: "Second Todo",
		Done:  true,
	},
	{
		Id:    "3",
		Title: "Third Todo",
		Done:  false,
	},
}

func main() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos", addTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", getTodoById).Methods("GET")
	router.HandleFunc("/todos/{id}", updateTodoById).Methods("PATCH")
	router.HandleFunc("/todos/{id}", deleteTodoById).Methods("DELETE")

	fmt.Println("Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatal(err)
	}
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

func getTodoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, todo := range todos {
		if todo.Id == id {
			json.NewEncoder(w).Encode(todo)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todos = append(todos, todo)

	json.NewEncoder(w).Encode(todos)
	w.WriteHeader(http.StatusCreated)
}

func deleteTodoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, todo := range todos {
		if todo.Id == id {
			fmt.Print(todo)
			todos = append(todos[:i], todos[i+1:]...)
		}
	}

	json.NewEncoder(w).Encode(todos)
}

func updateTodoById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	title := r.URL.Query().Get("title")
	done := r.URL.Query().Get("done")

	for i, todo := range todos {
		if todo.Id == id {
			if title != "" {
				todos[i].Title = title
			}
			if done != "" {
				todos[i].Done, _ = strconv.ParseBool(done)
			}
		}
	}

	json.NewEncoder(w).Encode(todos)
}
