package auth

import userModule "go-auth-app/modules/user"

type Module struct {
	Repository AuthRepository
	Service    AuthService
	Handler    *AuthHandler
}

func BuildModule(userService *userModule.Service) *Module {
	repository := NewRepository(nil)
	service := NewService(repository, userService)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}
