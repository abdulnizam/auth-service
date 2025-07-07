package service

import (
	"errors"
	"time"

	"auth-service/internal/model"
	"auth-service/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey") // Replace with env-based config in production

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func HashPassword(password string) (string, error) {
	return utils.HashPassword(password)
}

func CheckPasswordHash(password, hash string) bool {
	return utils.CheckPasswordHash(password, hash)
}

func AuthenticateUser(email, password string) (*model.User, error) {
	user, err := model.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
