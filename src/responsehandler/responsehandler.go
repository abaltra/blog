package responsehandler

import (
	"encoding/json"
	"net/http"
)

func EncodeJSONResponse(w http.ResponseWriter, content interface{}, status int, headers map[string]string) {
	if headers != nil {
		for header, value := range headers {
			w.Header().Set(header, value)
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(content)
}
