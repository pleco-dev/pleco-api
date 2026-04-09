package user

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func ToUserResponseList(users []User) []UserResponse {
	result := make([]UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, ToUserResponse(user))
	}

	return result
}
