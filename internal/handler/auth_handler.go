
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"auth-api/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepository *user.Repository
}

func NewAuthHandler(repository *user.Repository) *AuthHandler {
	return &AuthHandler{userRepository: repository}
}

func (handler *AuthHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /auth/register", handler.Register)
}

func (handler *AuthHandler) Register(
	response http.ResponseWriter,
	request *http.Request,
) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	request.Body = http.MaxBytesReader(response, request.Body, 1<<20)

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		http.Error(response, `{"error":"Invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if input.Name == "" || len(input.Name) > 100 {
		http.Error(response, `{"error":"Name is required and must be at most 100 characters"}`, http.StatusBadRequest)
		return
	}

	if len(input.Email) > 255 {
		http.Error(response, `{"error":"Invalid email address"}`, http.StatusBadRequest)
		return
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		http.Error(response, `{"error":"Invalid email address"}`, http.StatusBadRequest)
		return
	}

	if len(input.Password) < 8 || len(input.Password) > 72 {
		http.Error(response, `{"error":"Password must be between 8 and 72 bytes"}`, http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(response, `{"error":"Could not secure password"}`, http.StatusInternalServerError)
		return
	}

	contextWithTimeout, cancel := context.WithTimeout(request.Context(), 5*time.Second)
	defer cancel()

	createdUser, err := handler.userRepository.Create(
		contextWithTimeout,
		input.Name,
		input.Email,
		string(passwordHash),
	)
	if errors.Is(err, user.ErrEmailAlreadyExists) {
		http.Error(response, `{"error":"Email is already registered"}`, http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(response, `{"error":"Could not create account"}`, http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(response).Encode(map[string]any{
		"message": "Account created successfully",
		"user":    createdUser,
	})
}
