package server

import (
	"net/http"

	"github.com/datsfilipe/pkg/controllers"
)

func ServeMux(
	userController *controllers.UserController,
	messageController *controllers.MessageController,
	channelController *controllers.ChannelController,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/signup", Middleware(http.HandlerFunc(userController.HandleSignUp)))
	mux.Handle("/signin", Middleware(http.HandlerFunc(userController.HandleSignIn)))
	mux.Handle("/message", Middleware(http.HandlerFunc(messageController.HandleNewMessage)))
	mux.Handle("/channel", Middleware(http.HandlerFunc(channelController.HandleNewChannel)))
	mux.Handle("/channel/list", Middleware(http.HandlerFunc(channelController.HandleGetChannels)))

	return mux
}
