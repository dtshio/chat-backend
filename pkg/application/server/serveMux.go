package server

import (
	"net/http"

	"github.com/datsfilipe/pkg/controllers"
)

func ServeMux(userController *controllers.UserController) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/signup", http.HandlerFunc(userController.HandleSignUp))
	mux.Handle("/signin", http.HandlerFunc(userController.HandleSignIn))

	return mux
}
