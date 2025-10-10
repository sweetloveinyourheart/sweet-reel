package handlers

// Handlers holds all the request handlers
type Handlers struct {
	AuthHandler     IAuthHandler
	VideoManagement IVideoManagementHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		VideoManagement: NewVideoManagementHandler(),
		AuthHandler:     NewAuthHandler(),
	}
}
