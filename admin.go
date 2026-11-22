package main

import(
	"net/http"
	"fmt"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "forbidden",  fmt.Errorf("access denied"))
		return
	}
	if err := cfg.dbQueries.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, 400, "error deleting users", err)
		
	}
	w.Write([]byte("OK"))
}