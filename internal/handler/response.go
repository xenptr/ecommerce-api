package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type validationErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write(buf.Bytes())
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func writeValidationError(w http.ResponseWriter, errs validator.ValidationErrors) {
	fields := make(map[string]string)

	for _, err := range errs {
		fields[err.Field()] = validationMessage(err)
	}

	writeJSON(w, http.StatusBadRequest, validationErrorResponse{
		Error:  "validation failed",
		Fields: fields,
	})
}

func validationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", err.Param())
	default:
		return "is invalid"
	}
}
