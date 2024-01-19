package core

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ControllerMethod http.HandlerFunc

type Controller struct {
	methods []ControllerMethod
	db *gorm.DB
	log *zap.Logger
}

func (c *Controller) SetDB(db *gorm.DB) {
	c.db = db
}

func (c *Controller) SetLogger(log *zap.Logger) {
	c.log = log
}

func (c *Controller) GetPayload(r *http.Request) Map {
	payload := *new(Map)

	var raw json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&raw)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(raw, &payload)
	if err != nil {
		return nil
	}

	return payload
}

func generateMessage(statusCode int) string {
	switch statusCode {
	case http.StatusMethodNotAllowed:
		return "Method not allowed"
	case http.StatusBadRequest:
		return "Bad request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusNotFound:
		return "Not found"
	case http.StatusInternalServerError:
		return "Internal server error"
	case http.StatusCreated:
		return "Created"
	case http.StatusOK:
		return "OK"
	default:
		return "Unknown error"
	}
}


func (c *Controller) Response(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		http.Error(w, generateMessage(statusCode), statusCode)
		return
	}

	switch v := data.(type) {
	case error:
		http.Error(w, v.Error(), statusCode)
		return
	case string:
		w.WriteHeader(statusCode)
		w.Write([]byte(v))
	case []byte:
		w.WriteHeader(statusCode)
		w.Write(v)
	default:
		http.Error(w, "Unsupported data type", http.StatusInternalServerError)
	}
}

func (c *Controller) IsAllowedMethod(r *http.Request, methods []string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}

	return false
}

func NewController(methods []ControllerMethod) *Controller {
	return &Controller{
		methods: methods,
	}
}
