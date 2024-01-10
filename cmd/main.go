package main

import (
	"net/http"

	"github.com/datsfilipe/pkg/application/server"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		  return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			server.Init,
			server.ServeMux,
			zap.NewProduction,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
