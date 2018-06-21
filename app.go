package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// App struct stores the router and the database of the application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func respondWithError(res http.ResponseWriter, code int, message string) {
	respondWithJSON(res, code, map[string]string{"error": message})
}
func respondWithJSON(res http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(response)
}

func (a *App) getRecipe(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(res, http.StatusBadRequest, "Invalid user ID")
		return
	}
	r := recipe{ID: id}
	if err := r.getRecipe(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(res, http.StatusNotFound, "Recipe not found")
		default:
			respondWithError(res, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(res, http.StatusOK, r)
}

func (a *App) getRecipes(res http.ResponseWriter, req *http.Request) {
	recipes, err := getRecipes(a.DB)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, recipes)
}

func (a *App) createRecipe(res http.ResponseWriter, req *http.Request) {
	var r recipe
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		respondWithError(res, http.StatusBadRequest, "Bad request")
		return
	}
	defer req.Body.Close()
	if err := r.createRecipe(a.DB); err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusCreated, r)
}

func (a *App) updateRecipe(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(res, http.StatusBadRequest, "Invalid recipe ID")
		return
	}
	var r recipe
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&r); err != nil {
		respondWithError(res, http.StatusBadRequest, "Bad request")
		return
	}
	defer req.Body.Close()
	r.ID = id
	if err := r.updateRecipe(a.DB); err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, r)
}

func (a *App) deleteRecipe(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(res, http.StatusBadRequest, "Invalid Recipe ID")
		return
	}
	r := recipe{ID: id}
	if err := r.deleteRecipe(a.DB); err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/recipes", a.getRecipes).Methods("GET")
	a.Router.HandleFunc("/recipe", a.createRecipe).Methods("POST")
	a.Router.HandleFunc("/recipe/{id}", a.getRecipe).Methods("GET")
	a.Router.HandleFunc("/recipe/{id}", a.updateRecipe).Methods("PUT")
	a.Router.HandleFunc("/recipe/{id}", a.deleteRecipe).Methods("DELETE")
}

func (a *App) createDB() {
	var err error

	_, err = a.DB.Exec("CREATE DATABASE IF NOT EXISTS rest_api")
	if err != nil {
		log.Fatal(err)
	}

	_, err = a.DB.Exec("USE rest_api")
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) createTable() {
	statement := `
	CREATE TABLE IF NOT EXISTS recipes
	(
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(50) NOT NULL,
		ingredients VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL
	)`
	_, err := a.DB.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize connects with mysql, creating the database and table if necesary
// Also initialize the router with its handlers
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	a.createDB()
	a.createTable()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run starts the server
func (a *App) Run(addr string) {
	server := &http.Server{
		Addr:           addr,
		Handler:        a.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
}
