package server

import (
	"net/http"

	"github.com/datsfilipe/pkg/controllers"
)

func ServeMux() *http.ServeMux {
	userController := controllers.NewUserController()

	mux := http.NewServeMux()
	mux.Handle("/signup", http.HandlerFunc(userController.HandleSignUp))

	return mux
}
