package main

import(
	"net/http"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/pippps/sire-felcevi/internal/database"
	"github.com/pippps/sire-felcevi/internal/auth"
)

type UserReq struct {
	Username    string `json:"username"`
	Password string `json:"password"`
}

type UserResp struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter,r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userJson := &UserReq{}
	err := decoder.Decode(userJson)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	hashedPassword, err := auth.HashPassword(userJson.Password)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Username:          userJson.Username,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 200, UserResp{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
	})
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userJson := &UserReq{}

	err := decoder.Decode(userJson)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	dbUser, err := cfg.dbQueries.GetUser(r.Context(), userJson.Username)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	ok, err := auth.CheckPasswordHash(userJson.Password, dbUser.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	userResp := UserResp{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Username:        dbUser.Username,
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, 200, userResp)
}