package controllers

import (
	"fmt"
	"net/http"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"gorm.io/gorm"
)

type ChannelController struct {
	core.Controller
	db *gorm.DB
}

func (cc *ChannelController) HandleNewChannel(w http.ResponseWriter, r *http.Request) {
	if cc.IsAllowedMethod(r, []string{"PUT"}) == false {
		cc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if cc.IsAuthorized(r) == false {
		cc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	channelService := services.NewChannelService()
	newChannel, err := channelService.CreateChannel(cc.db, &models.Channel{})

	if err != nil {
		cc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	cc.Response(w, http.StatusCreated, fmt.Sprint("{\"id\": \"", newChannel.(models.Channel).ID, "\"}"))
}

func (cc *ChannelController) HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	if cc.IsAllowedMethod(r, []string{"GET"}) == false {
		cc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if cc.IsAuthorized(r) == false {
		cc.Response(w, http.StatusUnauthorized, nil)
		return
	}

	channelService := services.NewChannelService()
	channelRecords, err := channelService.GetChannels(cc.db, &[]models.Channel{})

	if err != nil {
		cc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	var res string
	for _, channel := range channelRecords.([]models.Channel) {
		res += fmt.Sprint("{\"id\": \"", channel.ID, "\", \"type\": \"", channel.Type, "\"},")
	}

	cc.Response(w, http.StatusCreated, fmt.Sprint("[", res[:len(res) - 1], "]"))
}

func NewChannelController(db *gorm.DB) *ChannelController {
	channelController := &ChannelController{
		db: db,
	}

	return &ChannelController{
		Controller: *core.NewController([]core.ControllerMethod{
			channelController.HandleNewChannel,
			channelController.HandleGetChannels,
		}),
		db: db,
	}
}
