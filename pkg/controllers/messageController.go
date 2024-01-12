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

type MessageController struct {
	core.Controller
	db *gorm.DB
}

func (mc *MessageController) HandleNewMessage(w http.ResponseWriter, r *http.Request) {
	message := &models.Message{}
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	var raw json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = json.Unmarshal(raw, &message)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid Message data")
		return
	}

	if auth.VerifyToken(fmt.Sprint(message.AuthorID) + "." + token) == false {
		core.Response(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	messageService := services.NewMessageService()

	newMessage, err := messageService.CreateMessage(mc.db, message)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	res, err := json.Marshal(newMessage)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
	}

	core.Response(w, http.StatusCreated, res)
}

func (mc *MessageController) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	payload := &core.GetMessagesPayload{}
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	var raw json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = json.Unmarshal(raw, &payload)
	if err != nil {
		core.Response(w, http.StatusBadRequest, "Invalid Message data")
		return
	}

	if auth.VerifyToken(fmt.Sprint(payload.ID) + "." + token) == false {
		core.Response(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	messageService := services.NewMessageService()

	messages, err := messageService.GetMessages(mc.db, payload)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	res, err := json.Marshal(messages)
	if err != nil {
		core.Response(w, http.StatusInternalServerError, "Internal server error")
	}

	core.Response(w, http.StatusOK, res)
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
