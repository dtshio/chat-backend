package services

import (
	"fmt"

	"github.com/datsfilipe/pkg/application/database"
	"github.com/datsfilipe/pkg/core"
	"github.com/datsfilipe/pkg/models"
	"github.com/datsfilipe/pkg/repositories"
)

type MessageService struct {
	core.Service
}

func (ms *MessageService) CreateMessage(data interface{}) (interface{}, error) {
	Message := data.(*models.Message)
	if Message == nil {
		return nil, fmt.Errorf("Data is not a valid Message")
	}

	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	messageRepo := repositories.NewMessageRepository()

	Message.ID = core.GenerateID()

	return messageRepo.CreateMessage(db, Message)
}

func NewMessageService() *MessageService {
	messageService := &MessageService{}

	return &MessageService{
		Service: *core.NewService(
			[]core.ServiceMethod{
				messageService.CreateMessage,
			},
		),
	}
}
