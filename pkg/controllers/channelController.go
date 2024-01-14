package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/datsfilipe/pkg/application/auth"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"gorm.io/gorm"
)

type ChannelController struct {
	core.Controller
	db *gorm.DB
}

type HandleNewChannelPaylaod struct {
	ID string `json:"id"`
}

func (cc *ChannelController) HandleNewChannel(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	var raw json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	payload := HandleNewChannelPaylaod{}
	err = json.Unmarshal(raw, &payload)

	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid Channel data")
		return
	}

	if auth.VerifyToken(fmt.Sprint(payload.ID) + "." + token) == false {
		core.Response(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	channelService := services.NewChannelService()

	channel := &models.Channel{}
	newChannel, err := channelService.CreateChannel(cc.db, channel)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	core.Response(w, http.StatusCreated, fmt.Sprint("{\"id\": \"", newChannel.(models.Channel).ID, "\"}"))
}

func (cc *ChannelController) HandleGetChannels(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	var raw json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	payload := HandleNewChannelPaylaod{}
	err = json.Unmarshal(raw, &payload)

	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid Channel data")
		return
	}

	if auth.VerifyToken(fmt.Sprint(payload.ID) + "." + token) == false {
		core.Response(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	channelService := services.NewChannelService()

	channels := &[]models.Channel{}
	channelRecords, err := channelService.GetChannels(cc.db, channels)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var res string
	for _, channel := range channelRecords.([]models.Channel) {
		res += fmt.Sprint("{\"id\": \"", channel.ID, "\", \"type\": \"", channel.Type, "\"},")
	}

	core.Response(w, http.StatusCreated, fmt.Sprint("[", res[:len(res) - 1], "]"))
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
