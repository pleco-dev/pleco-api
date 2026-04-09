package tests

import (
	"go-auth-app/middleware"
	user "go-auth-app/modules/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type stubUserRepository struct{}

func (s *stubUserRepository) Create(_ *user.User) error                { return nil }
func (s *stubUserRepository) FindByEmail(_ string) (*user.User, error) { return nil, nil }
func (s *stubUserRepository) FindByID(_ uint) (*user.User, error)      { return nil, nil }
func (s *stubUserRepository) Update(_ *user.User) error                { return nil }
func (s *stubUserRepository) FindAll() ([]user.User, error)            { return nil, nil }
func (s *stubUserRepository) FindAllWithFilter(_, _ int, _, _ string) ([]user.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserRepository) Delete(_ uint) error { return nil }

func TestAdminForbidden(t *testing.T) {
	userService := &user.Service{UserRepo: &stubUserRepository{}}
	handler := user.Handler{UserService: userService}

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})

	r.GET("/admin/users", middleware.AdminOnly(), func(c *gin.Context) {
		handler.GetAllUsers(c)
	})

	req, _ := http.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}
