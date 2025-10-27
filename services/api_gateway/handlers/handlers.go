package handlers

// Handlers holds all the request handlers
type Handlers struct {
	AuthHandler     IAuthHandler
	ChannelHandler  IChannelHandler
	VideoManagement IVideoManagementHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		VideoManagement: NewVideoManagementHandler(),
		ChannelHandler:  NewChannelHandler(),
		AuthHandler:     NewAuthHandler(),
	}
}
