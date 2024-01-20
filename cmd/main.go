package main

import (
	"net/http"

	"github.com/datsfilipe/pkg/application/database"
	"github.com/datsfilipe/pkg/application/server"
	"github.com/datsfilipe/pkg/controllers"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		  return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			server.Init,
			server.ServeMux,
			controllers.NewUserController,
			controllers.NewFriendshipController,
			controllers.NewMessageController,
			controllers.NewGroupController,
			database.Open,
			zap.NewDevelopment,
		),
		fx.Invoke(func(
			_ *http.Server,
			uc *controllers.UserController,
			fc *controllers.FriendshipController,
			mc *controllers.MessageController,
			gc *controllers.GroupController,
			db *gorm.DB,
			log *zap.Logger,
		) {
			uc.SetDB(db)
			uc.SetLogger(log)

			fc.SetDB(db)
			fc.SetLogger(log)

			mc.SetDB(db)
			mc.SetLogger(log)

			gc.SetDB(db)
			gc.SetLogger(log)
		}),
	).Run()
}
