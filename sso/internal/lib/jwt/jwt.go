package jwt

import (
	"log/slog"
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {

	var jwtSecretKey = []byte("GOIDA")

	payload := jwt.MapClaims{
		"uid":   user.Id,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenSrting, err := token.SignedString(jwtSecretKey)
	if err != nil {
		slog.Info("signin token error ")
		return "", err
	}

	return tokenSrting, nil
}
