package dto

import "go-auth-app/models"

func ToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func ToUserResponseList(users []models.User) []UserResponse {
	var result []UserResponse

	for _, user := range users {
		result = append(result, ToUserResponse(user))
	}

	return result
}
