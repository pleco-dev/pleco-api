package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserService *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{UserService: userService}
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	search := c.Query("search")
	role := c.Query("role")

	users, total, err := h.UserService.GetAllUsers(page, limit, search, role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(200, gin.H{
		"data": ToUserResponseList(users),
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user id"})
		return
	}

	if err := h.UserService.DeleteUser(uint(id64)); err != nil {
		c.JSON(500, gin.H{"error": "Failed"})
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})
}
