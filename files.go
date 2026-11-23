package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pippps/sire-felcevi/internal/auth"
	"github.com/pippps/sire-felcevi/internal/database"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, userID uuid.UUID)

type fileReq struct {
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type fileResp struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Etag        string    `json:"etag"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	UserID      uuid.UUID `json:"user_id"`
}

// nullString turns an empty string into SQL NULL.
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func toFileResp(m database.File) fileResp {
	return fileResp{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Name:        m.Name,
		Etag:        m.Etag,
		ContentType: m.ContentType,
		Size:        m.Size,
		UserID:      m.UserID,
	}
}

func (cfg *apiConfig) tokenValidationMiddleware(next authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
			return
		}
		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
			return
		}
		next(w, r, userID)
	})
}

func (cfg *apiConfig) createFileHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {

	params := fileReq{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name is required", nil)
		return
	}

	file, err := cfg.dbQueries.Createfile(r.Context(), database.CreatefileParams{
		ContentType: params.ContentType,
		Name:        params.Name,
		Size:        params.Size,
		UserID:      userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create file", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, toFileResp(file))
}

func (cfg *apiConfig) deleteFileHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	fileIDString := r.PathValue("fileID")
	fileID, err := uuid.Parse(fileIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	file, err := cfg.dbQueries.GetFile(r.Context(), fileID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get image", err)
		return
	}
	if file.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete image", err)
		return
	}

	err = cfg.dbQueries.DeleteFile(r.Context(), fileID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete  image", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) getFilesHandler(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {

	files, err := cfg.dbQueries.GetFiles(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve videos", err)
		return
	}

	if len(files) == 0 {
		respondWithJSON(w, http.StatusOK, files)
		return
	}

	filesResp := make([]fileResp, len(files))

	for i, _ := range files {
		filesResp[i] = toFileResp(files[i])
	}

	respondWithJSON(w, http.StatusOK, filesResp)
}
