package main

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/datsfilipe/pkg/server"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		  return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			server.Init,
			server.ServeMux,
			server.CreateEchoHandler,
			zap.NewProduction,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
