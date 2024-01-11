package core

import (
	"net/http"
)

func Response(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

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
