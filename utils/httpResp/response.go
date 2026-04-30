// myapp/utils/httpResp/response.go
package httpResp

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON writes an HTTP response with a JSON body.
//
//	w      — the response writer provided by the handler
//	code   — HTTP status code, e.g. http.StatusOK (200), http.StatusCreated (201)
//	payload — any Go value; it will be marshalled to JSON automatically
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to marshal response")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError writes an HTTP response with an error JSON body.
//
//	e.g. {"error": "Student not found"}
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}
