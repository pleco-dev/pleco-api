package services

import (
	"errors"
	"fmt"
	"go-auth-app/config"
	"go-auth-app/models"
	"go-auth-app/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repositories.UserRepository
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Register(user *models.User, password string) error {
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	user.Role = "user"
	return s.UserRepo.Create(user)
}

func (s *AuthService) Login(email, password string) (*AuthTokens, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	fmt.Println("LOGIN SECRET:", string(config.JWTSecret))

	// ✅ check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// ✅ generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	// ✅ generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	var jwtKey = config.JWTSecret

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	// ✅ simpan refresh token
	user.RefreshToken = refreshTokenString
	err = s.UserRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return &AuthTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *AuthService) Logout(userID uint, refreshToken string) error {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// optional: validasi refresh token
	if user.RefreshToken != refreshToken {
		return errors.New("invalid refresh token")
	}

	// hapus refresh token (invalidate session)
	user.RefreshToken = ""

	return s.UserRepo.Update(user)
}
