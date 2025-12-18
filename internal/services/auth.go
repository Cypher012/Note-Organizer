package services

import (
	"errors"
	"os"
	"time"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateJWT(user *models.User) (accessTokenString string, refreshTokenString string, err error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", "", errors.New("jwt secret missing")
	}

	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err = accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err = refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func RegisterUser(db *gorm.DB, body *models.RegisterRequest) (user *models.User, accessToken string, refreshToken string, err error) {
	if body.Email == "" || body.Password == "" || body.Username == "" {
		return nil, "", "", errors.New("missing required fields")
	}

	existing := new(models.User)

	if err := db.Where("email = ?", body.Email).First(existing).Error; err == nil {
		return nil, "", "", errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	user = &models.User{
		ID:       uuid.NewString(),
		Username: body.Username,
		Email:    body.Email,
		Password: string(hashed),
	}

	if err := db.Create(user).Error; err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err = GenerateJWT(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func LoginUser(db *gorm.DB, body *models.LoginRequest) (dbUser *models.User, accessToken string, refreshToken string, err error) {
	dbUser = new(models.User)
	if err := db.Where("email = ?", body.Email).First(dbUser).Error; err != nil {
		return nil, "", "", errors.New("Invalid credential")
	}
	password := body.Password
	hashedPassword := dbUser.Password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, "", "", errors.New("Invalid credential")
	}

	accessToken, refreshToken, err = GenerateJWT(dbUser)
	if err != nil {
		return nil, "", "", err
	}

	return dbUser, accessToken, refreshToken, nil
}
