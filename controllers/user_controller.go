package controllers

import (
	"fmt"
	"go-auth-app/dto"
	"go-auth-app/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func (u *UserController) GetAllUsers(c *gin.Context) {
	// default values
	page := 1
	limit := 10

	// query params
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	search := c.Query("search")
	role := c.Query("role")

	users, total, err := u.UserService.GetAllUsers(page, limit, search, role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}

	var result []dto.UserResponse

	for _, user := range users {
		result = append(result, dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	response := dto.ToUserResponseList(users)

	c.JSON(200, gin.H{
		"data": response,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (u *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")

	var id uint
	fmt.Sscanf(idParam, "%d", &id)

	err := u.UserService.DeleteUser(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})
}

func (a *UserController) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// ambil user dari repository
	user, err := a.UserService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// response (hindari kirim password!)
	c.JSON(200, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}
