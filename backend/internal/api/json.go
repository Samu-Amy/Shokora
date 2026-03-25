package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// - JSON Validator -
func formatValidationErrors(vErr validator.ValidationErrors) map[string]string {
	out := make(map[string]string)

	for _, fieldError := range vErr {
		field := strings.ToLower(fieldError.Field())

		// if field already in out -> skip (only one error per field)
		if _, exists := out[field]; exists {
			continue
		}

		out[field] = fieldError.Tag()
	}

	return out
}

// - JSON Encoding/Decoding -
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

// Used to return error in a {"error": error} JSON
func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

// Used to return error in a {"error": error} JSON
func writeValidatorJSONError(w http.ResponseWriter, status int, vErr validator.ValidationErrors) error {

	return writeJSON(w, status, formatValidationErrors(vErr))
}

// Used to return data in a {"data": data} JSON
func (app *App) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}
