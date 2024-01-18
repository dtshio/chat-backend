package server

import (
	"net/http"

	"github.com/datsfilipe/pkg/controllers"
)

func Route(
	mux *http.ServeMux,
	path string,
	handler http.HandlerFunc,
	middlewares ...func(next http.Handler) http.Handler,
) {
	var handlerWithMiddlewares http.Handler = handler

	for _, middleware := range middlewares {
		handlerWithMiddlewares = middleware(handlerWithMiddlewares)
	}

	mux.Handle(path, handlerWithMiddlewares)
}

func ServeMux(
	userController *controllers.UserController,
	messageController *controllers.MessageController,
	friendshipController *controllers.FriendshipController,
) *http.ServeMux {
	mux := http.NewServeMux()

	Route(mux, "/signup", userController.HandleSignUp, CorsMiddleware)
	Route(mux, "/signin", userController.HandleSignIn, CorsMiddleware)

	Route(mux, "/message", messageController.HandleNewMessage, CorsMiddleware, AuthMiddleware)
	Route(mux, "/message/history", messageController.HandleGetMessages, CorsMiddleware, AuthMiddleware)

	Route(mux, "/friendship/list", friendshipController.HandleGetFriendships, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship/delete", friendshipController.HandleDeleteFriendship, CorsMiddleware, AuthMiddleware)

	Route(mux, "/friendship-requests", friendshipController.HandleNewFriendshipRequest, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship-requests/accept", friendshipController.HandleNewFriendship, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship-requests/list", friendshipController.HandleGetFriendshipRequests, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship-requests/delete", friendshipController.HandleDeleteFriendshipRequest, CorsMiddleware, AuthMiddleware)

	return mux
}
