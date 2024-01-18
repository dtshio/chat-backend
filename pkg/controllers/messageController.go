package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"gorm.io/gorm"
)

type MessageController struct {
	core.Controller
	db *gorm.DB
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

	messageService := services.NewMessageService()

	newMessage, err := messageService.CreateMessage(mc.db, message)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	jsonMessage, err := json.Marshal(newMessage)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
	}

	mc.Response(w, http.StatusCreated, jsonMessage)
}

func (mc *MessageController) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	if mc.IsAllowedMethod(r, []string{"POST"}) == false {
		mc.Response(w, http.StatusMethodNotAllowed, nil)
		return
	}

	payload := mc.GetPayload(r)

	if payload == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	if payload["channel_id"] == nil || payload["page"] == nil {
		mc.Response(w, http.StatusBadRequest, nil)
		return
	}

	messageService := services.NewMessageService()
	messages, err := messageService.GetMessages(mc.db, payload)

	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
		return
	}

	jsonMessages, err := json.Marshal(messages)
	if err != nil {
		mc.Response(w, http.StatusInternalServerError, nil)
	}

	mc.Response(w, http.StatusOK, jsonMessages)
}

func NewMessageController(db *gorm.DB) *MessageController {
	messageController := &MessageController{
		db: db,
	}

	return &MessageController{
		Controller: *core.NewController([]core.ControllerMethod{
			messageController.HandleNewMessage,
			messageController.HandleGetMessages,
		}),
		db: db,
	}
}
