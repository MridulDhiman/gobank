package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store Storage
}

// create ne pointer to api server instance 
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

type apiFunc  func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}




func (s* APIServer) Run() {
	router:= mux.NewRouter();
	// default handlerFunc does not returns error, here we returning errors.
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleAccountById))
	
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

	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s* APIServer) handleAccountById(w http.ResponseWriter, r* http.Request) error {
	// account := NewAccount("Mridul", "Dhiman")
	id := mux.Vars(r)["id"]
	fmt.Print(id)

	if(r.Method == "DELETE") {
		return s.handleDeleteAccount(w,r)
	}
	return nil
}

func (s* APIServer) handleCreateAccount(w http.ResponseWriter, r* http.Request) error {
	
	// decode the json 
	accountRequest := new(CreateAccountRequest)
	if err:= json.NewDecoder(r.Body).Decode(accountRequest); err != nil {
		return err
	}

	newAccount := NewAccount(accountRequest.FirstName, accountRequest.LastName)

	if err:= s.store.CreateAccount(newAccount); err != nil {
		return err
	}

	
	return WriteJSON(w, http.StatusOK, newAccount)
}
func (s* APIServer) handleDeleteAccount(w http.ResponseWriter, r* http.Request) error  {
	id := mux.Vars(r)["id"]

	actualId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	if err:= s.store.DeleteAccount(actualId); err != nil {
		return err
	}

	return nil
}

func (s* APIServer) handleTransfer(w http.ResponseWriter, r* http.Request) error {
	return nil
}

func WriteJSON (w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func makeHTTPHandlerFunc(f apiFunc)  http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if err := f(w,r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{ Error: err.Error()} )
		}
	}
}