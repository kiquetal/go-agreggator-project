package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/kiquetal/go-agreggator-project/internal/database"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func (api *myApi) healthHandler(writer http.ResponseWriter, _ *http.Request) {

	api.respondWithJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (api *myApi) simulateError(writer http.ResponseWriter, _ *http.Request) {

	api.respondWithError(writer, http.StatusInternalServerError, "Internal Server Error")
}

func (api *myApi) verifyHeaderMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			api.respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		apikey := strings.Split(r.Header.Get("Authorization"), "ApiKey ")[1]
		ctx := context.WithValue(r.Context(), "apikey", apikey)
		next(w, r.WithContext(ctx))
	}
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

func (api *myApi) obtainUser(writer http.ResponseWriter, request *http.Request) {

	apikey := request.Context().Value("apikey").(string)
	user, err := api.DB.GetUserByApiKey(context.Background(), sql.NullString{String: apikey, Valid: true})
	if err != nil {
		// check if user is not found
		if err == sql.ErrNoRows {
			api.respondWithError(writer, http.StatusNotFound, "User not found")
			return
		}

		api.respondWithError(writer, http.StatusInternalServerError, "Internal Server Error")
		log.Println("Error getting user: ", err)
		return
	}
	api.respondWithJSON(writer, http.StatusOK, struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		ApiKey    string    `json:"api_key"`
	}{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ApiKey:    user.ApiKey.String,
	})
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
	router := chi.NewRouter()
	router.HandleFunc("/api/v1/health", api.corsMiddleware(api.healthHandler))
	router.HandleFunc("/api/v1/err", api.corsMiddleware(api.simulateError))
	router.Post("/v1/users", api.corsMiddleware(api.createUsers))
	router.Get("/v1/users", api.corsMiddleware(api.verifyHeaderMiddleware(api.obtainUser)))
	srv := &http.Server{
		Addr:    ":" + portListener,
		Handler: router,
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
