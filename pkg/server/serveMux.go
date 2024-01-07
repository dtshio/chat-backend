package server

import (
	"net/http"
)

func ServeMux(echo *EchoHandler) *http.ServeMux {
  mux := http.NewServeMux()
  mux.Handle("/echo", echo)
  return mux
}
