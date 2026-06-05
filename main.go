package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rssag/internal/database"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("$PORT not set in environment")
	}

	conn, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	fmt.Println("Port:", portString)
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Starting server on port %s", portString)

	v1Router := chi.NewRouter()
	v1Router.Get("/test", handleReadiness)
	v1Router.Get("/err", handleError)
	v1Router.Post("/user", apiCfg.handleCreateUser)
	router.Mount("/v1", v1Router)

	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err)
	}
}
