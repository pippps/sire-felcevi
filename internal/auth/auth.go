package auth
import(
	"fmt"
	"net/http"
	"strings"
	"time"
	"crypto/rand"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	ok, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signMeth := jwt.SigningMethodHS256
	claims := jwt.RegisteredClaims{
		Issuer:    "felcevi-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(signMeth, claims)
	jwt, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	idString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	rawToken := headers.Get("Authorization")
	token, ok := strings.CutPrefix(rawToken, "Bearer ")
	if !ok {
		return rawToken, fmt.Errorf("unauthorized token")
	}
	return token, nil
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}