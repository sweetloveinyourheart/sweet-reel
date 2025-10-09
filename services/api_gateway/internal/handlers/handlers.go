package handlers

// Handlers holds all the request handlers
type Handlers struct {
	VideoManagement IVideoManagementHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		VideoManagement: NewVideoManagementHandler(),
	}
}
