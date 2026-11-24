package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pippps/sire-felcevi/internal/database"
)

type apiConfig struct {
	dbQueries *database.Queries
	platform  string
	jwtSecret string
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
		jwtSecret: os.Getenv("SECRET"),
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("POST /api/users", apiCfg.usersHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	mux.HandleFunc("POST /api/files", apiCfg.tokenValidationMiddleware(apiCfg.createFileHandler))
	mux.HandleFunc("DELETE /api/files/{fileID}", apiCfg.tokenValidationMiddleware(apiCfg.deleteFileHandler))
	mux.HandleFunc("GET    /api/files", apiCfg.tokenValidationMiddleware(apiCfg.getFilesHandler))
	mux.HandleFunc("GET    /api/files/{fileID}", apiCfg.tokenValidationMiddleware(apiCfg.getFileHandler))

	mux.HandleFunc("DELETE /admin/reset", apiCfg.resetHandler)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
