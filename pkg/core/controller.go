package core

import (
	"encoding/json"
	"net/http"
)

type ControllerMethod http.HandlerFunc

type Controller struct {
	methods []ControllerMethod
}

func (c *Controller) GetPayload(r *http.Request) Map {
    payload := make(Map)

    switch r.Method {
    case "POST", "PUT":
        decoder := json.NewDecoder(r.Body)
        err := decoder.Decode(&payload)

        if err != nil {
            return nil
        }
    case "GET", "DELETE":
        r.ParseForm()
        for key, values := range r.Form {
            if len(values) > 0 {
                payload[key] = values[0]
            }
        }
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
