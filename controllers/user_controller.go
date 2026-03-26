package controllers

import (
	"fmt"
	"go-auth-app/repositories"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserRepo repositories.UserRepository
}

func (u *UserController) GetAllUsers(c *gin.Context) {
	allUsers, err := u.UserRepo.FindAll()
	users := make([]interface{}, 0)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}

	for _, user := range allUsers {
		if user.Role != "admin" {
			users = append(users, user)
		}
	}

	c.JSON(200, users)
}

func (u *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")

	var id uint
	fmt.Sscanf(idParam, "%d", &id)

	err := u.UserRepo.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})
}
