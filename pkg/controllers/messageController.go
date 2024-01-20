package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MessageController struct {
	core.Controller
	db *gorm.DB
	log *zap.Logger
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

	service := services.NewMessageService()

	messageObject, err := service.CreateMessage(mc.db, mc.log, message)
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

	service := services.NewMessageService()

	messages, err := service.GetMessages(mc.db, mc.log, payload)
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

func NewMessageController(db *gorm.DB, log *zap.Logger) *MessageController {
	controller := &MessageController{
		db: db,
		log: log,
	}

	return &MessageController{
		Controller: *core.NewController([]core.ControllerMethod{
			controller.HandleNewMessage,
			controller.HandleGetMessages,
		}),
		db: db,
		log: log,
	}
}
