package auth

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SocialLoginRequest struct {
	Provider string `json:"provider" binding:"required,oneof=google facebook apple"`
	Token    string `json:"token"`
	IDToken  string `json:"id_token"`
}
