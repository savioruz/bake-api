package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/savioruz/bake/internal/domain/model"
	"github.com/savioruz/bake/internal/service"
	e "github.com/savioruz/bake/pkg/error"
	"github.com/savioruz/bake/pkg/helper"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	OrderService *service.OrderService
	Log          *logrus.Logger
}

func NewOrderHandler(orderService *service.OrderService, log *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
		Log:          log,
	}
}

// @Summary Get all orders
// @Description Get all orders
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort" Enums(id, user_id, product_id, address_id, quantity, total_price, status, created_at, updated_at)
// @Param order query string false "Order" Enums(ASC, DESC)
// @Success 200 {object} model.SuccessResponse[[]model.OrderResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security ApiKeyAuth
// @Router /orders [get]
func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	pagination := h.parsePagination(r)
	response, err := h.OrderService.GetAll(r.Context(), pagination)
	if err != nil {
		h.Log.Errorf("failed to get all orders: %v", err)
		switch {
		case errors.Is(err, e.ErrValidation):
			e.ErrorHandler(w, r, http.StatusBadRequest, err)
		default:
			e.ErrorHandler(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get order by ID
// @Description Get order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} model.SuccessResponse[model.OrderResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security ApiKeyAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	request := &model.GetOrderRequest{
		ID: helper.ParseParam(r),
	}

	response, err := h.OrderService.GetById(r.Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get order by id: %v", err)
		switch {
		case errors.Is(err, e.ErrValidation):
			e.ErrorHandler(w, r, http.StatusBadRequest, err)
		case errors.Is(err, e.ErrNotFound):
			e.ErrorHandler(w, r, http.StatusNotFound, err)
		default:
			e.ErrorHandler(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Create a new order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body model.CreateOrderRequest true "Order"
// @Success 201 {object} model.SuccessResponse[model.OrderResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security ApiKeyAuth
// @Router /orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	request := &model.CreateOrderRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		e.ErrorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	response, err := h.OrderService.Create(r.Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create order: %v", err)
		switch {
		case errors.Is(err, e.ErrValidation):
			e.ErrorHandler(w, r, http.StatusBadRequest, err)
		case errors.Is(err, e.ErrInsufficientStock):
			e.ErrorHandler(w, r, http.StatusBadRequest, err)
		default:
			e.ErrorHandler(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// parsePagination is a private helper function to parse pagination parameters
func (h *OrderHandler) parsePagination(r *http.Request) *model.OrderPagination {
	pagination := &model.OrderPagination{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
		Order: "desc",
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			pagination.Page = pageNum
		}
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if limitNum, err := strconv.Atoi(limit); err == nil && limitNum > 0 && limitNum <= 100 {
			pagination.Limit = limitNum
		}
	}
	if sort := r.URL.Query().Get("sort"); sort != "" {
		pagination.Sort = sort
	}
	if order := r.URL.Query().Get("order"); order != "" {
		pagination.Order = strings.ToLower(order)
	}

	return pagination
}
