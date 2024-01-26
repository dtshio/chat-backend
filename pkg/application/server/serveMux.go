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
	groupController *controllers.GroupController,
) *http.ServeMux {
	mux := http.NewServeMux()

	Route(mux, "/signup", userController.HandleSignUp, CorsMiddleware)
	Route(mux, "/signin", userController.HandleSignIn, CorsMiddleware)
	Route(mux, "/user/delete", userController.HandleDeleteUser, CorsMiddleware)

	Route(mux, "/profile", userController.HandleCreateProfile, CorsMiddleware)

	Route(mux, "/message", messageController.HandleNewMessage, CorsMiddleware, AuthMiddleware)
	Route(mux, "/message/history", messageController.HandleGetMessages, CorsMiddleware, AuthMiddleware)

	Route(mux, "/friendship/list", friendshipController.HandleGetFriendships, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship/delete", friendshipController.HandleDeleteFriendship, CorsMiddleware, AuthMiddleware)

	Route(mux, "/friendship-requests", friendshipController.HandleNewFriendshipRequest, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship-requests/accept", friendshipController.HandleNewFriendship, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship-requests/list", friendshipController.HandleGetFriendshipRequests, CorsMiddleware, AuthMiddleware)
	Route(mux, "/friendship-requests/delete", friendshipController.HandleDeleteFriendshipRequest, CorsMiddleware, AuthMiddleware)

	Route(mux, "/group", groupController.HandleNewGroup, CorsMiddleware, AuthMiddleware)
	Route(mux, "/group/list", groupController.HandleGetGroups, CorsMiddleware, AuthMiddleware)
	Route(mux, "/group/list-by-profile", groupController.HandleGetGroupsByProfile, CorsMiddleware, AuthMiddleware)
	Route(mux, "/group/delete", groupController.HandleDeleteGroup, CorsMiddleware, AuthMiddleware)

	Route(mux, "/group/member", groupController.HandleAddMember, CorsMiddleware, AuthMiddleware)
	Route(mux, "/group/member/get", groupController.HandleGetMember, CorsMiddleware, AuthMiddleware)
	Route(mux, "/group/member/list", groupController.HandleGetMembers, CorsMiddleware, AuthMiddleware)
	Route(mux, "/group/member/delete", groupController.HandleRemoveMember, CorsMiddleware, AuthMiddleware)

	return mux
}
