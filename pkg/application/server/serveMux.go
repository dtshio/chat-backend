package server

import (
	"net/http"

	"github.com/datsfilipe/pkg/controllers"
)

func ServeMux(
	userController *controllers.UserController,
	messageController *controllers.MessageController,
	friendshipController *controllers.FriendshipController,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/signup", CorsMiddleware(http.HandlerFunc(userController.HandleSignUp)))
	mux.Handle("/signin", CorsMiddleware(http.HandlerFunc(userController.HandleSignIn)))

	mux.Handle("/message", CorsMiddleware(http.HandlerFunc(messageController.HandleNewMessage)))
	mux.Handle("/message/history", CorsMiddleware(http.HandlerFunc(messageController.HandleGetMessages)))

	mux.Handle("/friendship/list", CorsMiddleware(http.HandlerFunc(friendshipController.HandleGetFriendships)))
	mux.Handle("/friendship/delete", CorsMiddleware(http.HandlerFunc(friendshipController.HandleDeleteFriendship)))

	mux.Handle("/friendship-requests", CorsMiddleware(http.HandlerFunc(friendshipController.HandleNewFriendshipRequest)))
	mux.Handle("/friendship-requests/accept", CorsMiddleware(http.HandlerFunc(friendshipController.HandleNewFriendship)))
	mux.Handle("/friendship-requests/list", CorsMiddleware(http.HandlerFunc(friendshipController.HandleGetFriendshipRequests)))
	mux.Handle("/friendship-requests/delete", CorsMiddleware(http.HandlerFunc(friendshipController.HandleDeleteFriendshipRequest)))

	return mux
}
