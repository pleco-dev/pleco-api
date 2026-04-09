package auth

import (
	"go-auth-app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService AuthService
}

func NewHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationError(c, utils.FormatValidationError(err))
		return
	}

	user := utils.DtoToUser(input.Name, input.Email)
	err := h.AuthService.Register(&user, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationError(c, utils.FormatValidationError(err))
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	tokens, err := h.AuthService.Login(input.Email, input.Password, deviceID, userAgent, ipAddress)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deviceID := c.GetHeader("X-Device-ID")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device id required"})
		return
	}

	err := h.AuthService.Logout(userID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout success"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, utils.FormatValidationError(err))
		return
	}

	tokens, err := h.AuthService.RefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, ok := utils.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.AuthService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.JSON(400, gin.H{"error": "token required"})
		return
	}

	err := h.AuthService.VerifyEmail(token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "email verified"})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	type ResendInput struct {
		Email string `json:"email" binding:"required,email"`
	}

	var input ResendInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationError(c, utils.FormatValidationError(err))
		return
	}

	err := h.AuthService.ResendVerification(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email resent"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, utils.FormatValidationError(err))
		return
	}

	err := h.AuthService.ForgotPassword(body.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "reset link sent"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, utils.FormatValidationError(err))
		return
	}

	err := h.AuthService.ResetPassword(body.Token, body.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "password updated"})
}

func (h *AuthHandler) SocialLogin(c *gin.Context) {
	var body struct {
		Provider string `json:"provider"`
		Token    string `json:"id_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ValidationError(c, utils.FormatValidationError(err))
		return
	}

	deviceID := "web"
	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	tokens, err := h.AuthService.SocialLogin(body.Provider, body.Token, deviceID, userAgent, ip)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tokens)
}
