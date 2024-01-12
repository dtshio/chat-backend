package server

import (
	"net/http"

	"github.com/datsfilipe/pkg/controllers"
)

func ServeMux(
	userController *controllers.UserController,
	messageController *controllers.MessageController,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/signup", http.HandlerFunc(userController.HandleSignUp))
	mux.Handle("/signin", http.HandlerFunc(userController.HandleSignIn))
	mux.Handle("/channel", http.HandlerFunc(channelController.HandleNewChannel))

	return mux
}
