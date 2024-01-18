package core

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type ControllerMethod http.HandlerFunc

type Controller struct {
	methods []ControllerMethod
	db *gorm.DB
}

func (c *Controller) SetDB(db *gorm.DB) {
	c.db = db
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

func generateErrorMessages(statusCode int) string {
	switch statusCode {
	case http.StatusMethodNotAllowed:
		return "Method not allowed"
	case http.StatusBadRequest:
		return "Invalid data format"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusNotFound:
		return "Not found"
	case http.StatusInternalServerError:
		return "Internal server error"
	default:
		return "Unknown error"
	}
}


func (c *Controller) Response(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		http.Error(w, generateErrorMessages(statusCode), statusCode)
		return
	}

	switch v := data.(type) {
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
