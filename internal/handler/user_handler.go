package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/savioruz/bake/internal/domain/model"
	"github.com/savioruz/bake/internal/service"
	e "github.com/savioruz/bake/pkg/error"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserService *service.UserService
	Log         *logrus.Logger
}

func NewUserHandler(userService *service.UserService, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Log:         log,
	}
}

// @Summary Register a new user
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserRegisterRequest true "User"
// @Success 201 {object} model.SuccessResponse[model.UserResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	var request model.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Log.Errorf("failed to decode request body: %v", err)
		e.ErrorHandler(w, r, http.StatusBadRequest, e.ErrValidation)
		return
	}

	response, err := h.UserService.CreateUser(r.Context(), &request)
	if err != nil {
		h.Log.Errorf("failed to create user: %v", err)
		switch {
		case errors.Is(err, e.ErrUserExists):
			e.ErrorHandler(w, r, http.StatusConflict, e.ErrUserExists)
		case errors.Is(err, e.ErrValidation):
			e.ErrorHandler(w, r, http.StatusBadRequest, e.ErrValidation)
		default:
			e.ErrorHandler(w, r, http.StatusInternalServerError, e.ErrInternalServer)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.NewSuccessResponse(response, nil))
}

// @Summary Login a user
// @Description Login a user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserLoginRequest true "User"
// @Success 200 {object} model.SuccessResponse[model.TokenResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	var request model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Log.Errorf("failed to decode request body: %v", err)
		e.ErrorHandler(w, r, http.StatusBadRequest, e.ErrValidation)
		return
	}

	response, err := h.UserService.Login(r.Context(), &request)
	if err != nil {
		h.Log.Errorf("failed to login: %v", err)
		switch {
		case errors.Is(err, e.ErrValidation):
			e.ErrorHandler(w, r, http.StatusBadRequest, e.ErrValidation)
		case errors.Is(err, e.ErrUserNotFound):
			e.ErrorHandler(w, r, http.StatusNotFound, e.ErrUserNotFound)
		case errors.Is(err, e.ErrCredential):
			e.ErrorHandler(w, r, http.StatusUnauthorized, e.ErrCredential)
		default:
			e.ErrorHandler(w, r, http.StatusInternalServerError, e.ErrInternalServer)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.NewSuccessResponse(response, nil))
}

// @Summary Refresh a user's token
// @Description Refresh a user's token
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.RefreshTokenRequest true "User"
// @Success 200 {object} model.SuccessResponse[model.TokenResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/refresh [post]
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	var request model.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.Log.Errorf("failed to decode request body: %v", err)
		e.ErrorHandler(w, r, http.StatusBadRequest, e.ErrValidation)
		return
	}

	response, err := h.UserService.RefreshToken(r.Context(), &request)
	if err != nil {
		h.Log.Errorf("failed to refresh token: %v", err)
		e.ErrorHandler(w, r, http.StatusInternalServerError, e.ErrInternalServer)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.NewSuccessResponse(response, nil))
}

// @Summary Get current user
// @Description Get current user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResponse[model.UserResponse]
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security ApiKeyAuth
// @Router /users/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	response, err := h.UserService.Me(r.Context())
	if err != nil {
		h.Log.Errorf("failed to get me: %v", err)
		switch {
		case errors.Is(err, e.ErrUnauthorized):
			e.ErrorHandler(w, r, http.StatusUnauthorized, e.ErrUnauthorized)
		default:
			e.ErrorHandler(w, r, http.StatusInternalServerError, e.ErrInternalServer)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.NewSuccessResponse(response, nil))
}
