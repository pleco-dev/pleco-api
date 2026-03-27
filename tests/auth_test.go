package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-auth-app/controllers"
	"go-auth-app/models"
	"go-auth-app/services"
	"go-auth-app/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupRouter(controller controllers.AuthController) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	r.POST("/logout", controller.Logout)

	return r
}

func TestRegister(t *testing.T) {
	mockRepo := &mocks.MockUserRepo{}

	authService := &services.AuthService{
		UserRepo: mockRepo,
	}

	controller := controllers.AuthController{
		AuthService: authService,
	}

	router := setupRouter(controller)

	body := map[string]string{
		"name":     "Test User",
		"email":    "register_test@mail.com",
		"password": "123456",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLogin(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), 14)

	mockRepo := &mocks.MockUserRepo{
		User: &models.User{
			Email:    "login_test@mail.com",
			Password: string(hash),
			Role:     "user",
		},
	}

	authService := &services.AuthService{
		UserRepo: mockRepo,
	}

	controller := controllers.AuthController{
		AuthService: authService,
	}

	router := setupRouter(controller)

	body := map[string]string{
		"email":    "login_test@mail.com",
		"password": "123456",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLogout(t *testing.T) {
	mockRepo := &mocks.MockUserRepo{
		User: &models.User{
			ID:           1,
			RefreshToken: "some_token",
		},
	}

	authService := &services.AuthService{
		UserRepo: mockRepo,
	}

	controller := controllers.AuthController{
		AuthService: authService,
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		controller.Logout(c)
	})

	// ✅ FIX: kirim body
	body := map[string]string{
		"refresh_token": "some_token",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/logout", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), 14)

	// Use the simple mock implementation (no testify/mock required).
	mockRepo := &mocks.MockUserRepo{
		User: &models.User{
			Email:    "test@mail.com",
			Password: string(hash), // hash for "correct_password"
			Role:     "user",
		},
	}

	authService := &services.AuthService{
		UserRepo: mockRepo,
	}

	controller := controllers.AuthController{
		AuthService: authService,
	}

	router := setupRouter(controller)

	body := map[string]string{
		"email":    "test@mail.com",
		"password": "wrong_password", // ❌ salah
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 🔥 expect gagal
	assert.Equal(t, 401, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepo{
		FindByEmailErr: errors.New("user not found"),
	}

	authService := &services.AuthService{
		UserRepo: mockRepo,
	}

	controller := controllers.AuthController{
		AuthService: authService,
	}

	router := setupRouter(controller)

	body := map[string]string{
		"email":    "notfound@mail.com",
		"password": "123456",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 🔥 expect gagal
	assert.Equal(t, 401, w.Code)
}

func TestLogin_InvalidInput(t *testing.T) {
	mockRepo := new(mocks.MockUserRepo)

	authService := &services.AuthService{
		UserRepo: mockRepo,
	}

	controller := controllers.AuthController{
		AuthService: authService,
	}

	router := setupRouter(controller)

	body := map[string]string{
		"email":    "",
		"password": "",
	}

	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Empty strings are valid JSON input, so the controller proceeds to the
	// 400 validation failed
	assert.Equal(t, 400, w.Code)
}
