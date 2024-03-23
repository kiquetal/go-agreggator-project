package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kiquetal/go-agreggator-project/internal/database"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type myApi struct {
	DB *database.Queries
}

//headers for cors
// Access-control-allow-headers = [

func (api *myApi) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next(w, r)
	}

}

func (api *myApi) respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	//write the body of w
	en := json.NewEncoder(w)
	if err := en.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (api *myApi) respondWithError(w http.ResponseWriter, status int, message string) {
	api.respondWithJSON(w, status, map[string]string{"error": message})
}

func (api *myApi) healthHandler(writer http.ResponseWriter, request *http.Request) {

	api.respondWithJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (api *myApi) simulateError(writer http.ResponseWriter, request *http.Request) {

	api.respondWithError(writer, http.StatusInternalServerError, "Internal Server Error")
}

func (api *myApi) createUsers(writer http.ResponseWriter, request *http.Request) {

	log.Println("Creating user")
	//check body
	type bodyPost struct {
		Name string `json:"name"`
	}

	var body bodyPost
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		api.respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if body.Name == "" {
		api.respondWithError(writer, http.StatusBadRequest, "Name is required")
		return

	}

	//create user
	user, err := api.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      body.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		api.respondWithError(writer, http.StatusInternalServerError, "Internal Server Error")
		log.Println("Error creating user: ", err)
		return

	}
	api.respondWithJSON(writer, http.StatusCreated, user)

}

func main() {

	err := godotenv.Load()
	if err != nil {
		return
	}
	portListener := os.Getenv("PORT")
	if portListener == "" {
		log.Fatalf("PORT is not set")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %v", err)

	}
	api := &myApi{
		DB: database.New(db),
	}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/api/v1/health", api.corsMiddleware(api.healthHandler))
	serverMux.HandleFunc("/api/v1/err", api.corsMiddleware(api.simulateError))
	serverMux.HandleFunc("/v1/users", api.corsMiddleware(api.createUsers))
	srv := &http.Server{
		Addr:    ":" + portListener,
		Handler: serverMux,
	}
	fmt.Println("Server is running on port: ", portListener)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Println("Error shutting down server: ", err)

	}
	log.Println("Server gracefully stopped")

}
