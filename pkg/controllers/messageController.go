package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
)

type MessageController struct {
	core.Controller
	service *services.MessageService
}

func (mc *MessageController) HandleNewMessage(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	message := &models.Message{}
	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	authorID, _ := strconv.ParseUint(payload["author_id"].(string), 10, 64)
	channelID, _ := strconv.ParseUint(payload["channel_id"].(string), 10, 64)

	message.Content = payload["content"].(string)
	message.AuthorID = models.BigInt(authorID)
	message.ChannelID = models.BigInt(channelID)

	if message.Content == "" || message.AuthorID == 0 || message.ChannelID == 0 {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	messageObject, err := mc.service.CreateMessage(message)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(messageObject)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusCreated, res)
}

func (mc *MessageController) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"GET"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	if payload["channel_id"] == "" || payload["page"] == "" {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	messages, err := mc.service.GetMessages(payload)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, err)
		return
	}

	res, err := json.Marshal(messages)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	mc.Response(w, http.StatusOK, res)
}

func NewMessageController(service *services.MessageService) *MessageController {
	return &MessageController{
		service: service,
	}
}
