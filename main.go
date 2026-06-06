package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rssag/internal/database"
	"time"

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

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
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

	go startScraping(db, 10, time.Minute)

	v1Router := chi.NewRouter()
	v1Router.Get("/test", handleReadiness)
	v1Router.Get("/err", handleError)

	v1Router.Post("/user", apiCfg.handleCreateUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handleGetUserByApiKey))

	v1Router.Get("/feed", apiCfg.middlewareAuth(apiCfg.handleGetFeedByApiKeyAndFeedId))
	v1Router.Get("/feedAll", apiCfg.middlewareAuth(apiCfg.handleGetAllFeeds))
	v1Router.Get("/feedUser", apiCfg.middlewareAuth(apiCfg.handleGetUserFeeds))
	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetUserFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteUserFollow))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetUserPosts))

	router.Mount("/v1", v1Router)

	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err)
	}
}
