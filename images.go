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

type imageMetaReq struct {
	Title       string `json:"title"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}

type imageMetaResp struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
}

// nullString turns an empty string into SQL NULL.
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func toImageMetaResp(m database.ImagesMd) imageMetaResp {
	return imageMetaResp{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Title:       m.Title.String,
		ImageURL:    m.ImageUrl,
		Description: m.Description.String,
		UserID:      m.UserID,
	}
}

func (cfg *apiConfig) imageMetaCreateHandler(w http.ResponseWriter, r *http.Request) {
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

	params := imageMetaReq{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	if params.ImageURL == "" {
		respondWithError(w, http.StatusBadRequest, "image_url is required", nil)
		return
	}

	imageMeta, err := cfg.dbQueries.CreateImageMeta(r.Context(), database.CreateImageMetaParams{
		Title:       nullString(params.Title),
		ImageUrl:    params.ImageURL,
		Description: nullString(params.Description),
		UserID:      userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create image meta data", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, toImageMetaResp(imageMeta))
}

func (cfg *apiConfig) imageMetaDeleteHandler(w http.ResponseWriter, r *http.Request) {
	imageIDString := r.PathValue("imageID")
	imageID, err := uuid.Parse(imageIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

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

	imageMeta, err := cfg.dbQueries.GetImageMeta(r.Context(), imageID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get image", err)
		return
	}
	if imageMeta.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete image", err)
		return
	}

	err = cfg.dbQueries.DeleteImageMeta(r.Context(), imageID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete  image", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
