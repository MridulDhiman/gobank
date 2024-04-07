package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}

// create ne pointer to api server instance 
func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{listenAddr: listenAddr}
}

type apiFunc  func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}


func WriteJSON (w http.ResponseWriter, status int, v any) error {
w.WriteHeader(status)
w.Header().Set("Content-Type", "application/json")
return json.NewEncoder(w).Encode(v)
}


func makeHTTPHandlerFunc(f apiFunc)  http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if err := f(w,r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{ Error: err.Error()} )
		}
 	}
}

func (s* APIServer) Run() {
	router:= mux.NewRouter();
// default handlerFunc does not returns error, here we returning errors.
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))

	log.Println("JSON API Server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
	
}


// golang by default, does not handles the methods, but we can access the methods using r.Method
func (s* APIServer) handleAccount(w http.ResponseWriter, r* http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w,r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w,r)
	}


	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w,r)
	}
return fmt.Errorf("method not allowed %s", r.Method)
}

func (s* APIServer) handleGetAccount(w http.ResponseWriter, r* http.Request) error {
	account := NewAccount("Mridul", "Dhiman")
return WriteJSON(w,http.StatusOK,account)
}

func (s* APIServer) handleCreateAccount(w http.ResponseWriter, r* http.Request) error {
return nil
}
func (s* APIServer) handleDeleteAccount(w http.ResponseWriter, r* http.Request) error  {
return nil
}

func (s* APIServer) handleTransfer(w http.ResponseWriter, r* http.Request) error {
return nil
}