package main

import (
	"net/http"
	"log"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/pippps/sire-felcevi/internal/database"
)

type apiConfig struct {
	dbQueries      *database.Queries
	platform       string
	secret         string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	port := "8080"
	mux := http.NewServeMux()

	apiCfg := &apiConfig{
		dbQueries: database.New(db),
		platform:  os.Getenv("PLATFORM"),
		secret:    os.Getenv("SECRET"),
	}

	srv := &http.Server{
		Addr : ":" + port,
		Handler : mux,
	}

	mux.HandleFunc("POST /api/users", apiCfg.usersHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("DEL /admin/reset", apiCfg.resetHandler)

	
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}