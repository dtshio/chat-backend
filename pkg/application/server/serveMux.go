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
	mux.Handle("/signup", http.HandlerFunc(userController.HandleSignUp))
	mux.Handle("/signin", http.HandlerFunc(userController.HandleSignIn))
	mux.Handle("/message", http.HandlerFunc(messageController.HandleNewMessage))
	mux.Handle("/channel", http.HandlerFunc(channelController.HandleNewChannel))

	return mux
}
