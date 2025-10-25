package handlers

// Handlers holds all the request handlers
type Handlers struct {
	AuthHandler     IAuthHandler
	UserHandler     IUserHandler
	VideoManagement IVideoManagementHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		VideoManagement: NewVideoManagementHandler(),
		UserHandler:     NewUserHandler(),
		AuthHandler:     NewAuthHandler(),
	}
}
