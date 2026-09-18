package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/xenptr/ecommerce-api/internal/dto"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.authService.Register(r.Context(), req)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			writeError(w, http.StatusInternalServerError, "error in validating data")
			return
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			writeValidationError(w, validationErrors)
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.authService.Login(r.Context(), req)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			writeError(w, http.StatusInternalServerError, "error in validating data")
			return
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			writeValidationError(w, validationErrors)
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"message": "user logged in successfully",
	})
}
