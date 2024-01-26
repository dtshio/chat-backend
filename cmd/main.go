package main

import (
	"net/http"

	"github.com/datsfilipe/pkg/application/database"
	"github.com/datsfilipe/pkg/application/server"
	"github.com/datsfilipe/pkg/controllers"
	"github.com/datsfilipe/pkg/repositories"
	"github.com/datsfilipe/pkg/services"
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
			controllers.NewUserController,
			controllers.NewFriendshipController,
			controllers.NewMessageController,
			controllers.NewGroupController,
			services.NewUserService,
			services.NewFriendshipService,
			services.NewMessageService,
			services.NewGroupService,
			services.NewChannelService,
			repositories.NewUserRepository,
			repositories.NewFriendshipRepository,
			repositories.NewMessageRepository,
			repositories.NewGroupRepository,
			repositories.NewChannelRepository,
			database.Open,
			zap.NewDevelopment,
		),
		fx.Invoke(func(
			*http.Server,
		) {
			// Do nothing
		}),
	).Run()
}
