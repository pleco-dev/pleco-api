package user

type Module struct {
	Repository Repository
	Service    *Service
	Handler    *Handler
}

func BuildModule() *Module {
	repository := NewRepository()
	service := NewService(repository)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}
