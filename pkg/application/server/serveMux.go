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

	mux.Handle("/signup", CorsMiddleware(http.HandlerFunc(userController.HandleSignUp)))
	mux.Handle("/signin", CorsMiddleware(http.HandlerFunc(userController.HandleSignIn)))
	mux.Handle("/message", CorsMiddleware(http.HandlerFunc(messageController.HandleNewMessage)))
	mux.Handle("/message/history", CorsMiddleware(http.HandlerFunc(messageController.HandleGetMessages)))
	mux.Handle("/channel", CorsMiddleware(http.HandlerFunc(channelController.HandleNewChannel)))
	mux.Handle("/channel/list", CorsMiddleware(http.HandlerFunc(channelController.HandleGetChannels)))

	return mux
}
