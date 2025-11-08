package handlers

// Handlers holds all the request handlers
type Handlers struct {
	AuthHandler    IAuthHandler
	ChannelHandler IChannelHandler
	Video          IVideoHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		Video:          NewVideoHandler(),
		ChannelHandler: NewChannelHandler(),
		AuthHandler:    NewAuthHandler(),
	}
}
