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

type ProductHandler struct {
	ProductService *service.ProductService
	Log            *logrus.Logger
}

func NewProductHandler(productService *service.ProductService, log *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
		Log:            log,
	}
}

// @Summary Get all products
// @Description Get all products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort" Enums(id, name, description, price, stock, image, created_at, updated_at)
// @Param order query string false "Order" Enums(ASC, DESC)
// @Success 200 {object} model.SuccessResponse[[]model.ProductResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	request := h.parsePagination(r)
	response, err := h.ProductService.GetAll(r.Context(), request)
	if err != nil {
		h.Log.Errorf("failed to get all products: %v", err)
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

// @Summary Search products
// @Description Search products
// @Tags products
// @Accept json
// @Produce json
// @Param id query string false "ID"
// @Param name query string false "Name"
// @Param description query string false "Description"
// @Param price query string false "Price"
// @Param stock query string false "Stock"
// @Param image query string false "Image"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort" Enums(id, name, description, price, stock, image, created_at, updated_at)
// @Param order query string false "Order" Enums(ASC, DESC)
// @Success 200 {object} model.SuccessResponse[[]model.ProductResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /products/search [get]
func (h *ProductHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	pagination := h.parsePagination(r)
	query := &model.ProductQuery{}

	if id := r.URL.Query().Get("id"); id != "" {
		query.ID = helper.StrToPtr(id)
	}

	if name := r.URL.Query().Get("name"); name != "" {
		query.Name = helper.StrToPtr(name)
	}

	if description := r.URL.Query().Get("description"); description != "" {
		query.Description = helper.StrToPtr(description)
	}

	if price := r.URL.Query().Get("price"); price != "" {
		if priceVal, err := strconv.ParseFloat(price, 64); err == nil {
			query.Price = helper.FloatToPtr(priceVal)
		} else {
			e.ErrorHandler(w, r, http.StatusBadRequest, errors.New("invalid price format"))
			return
		}
	}

	if stock := r.URL.Query().Get("stock"); stock != "" {
		if stockVal, err := strconv.Atoi(stock); err == nil {
			query.Stock = helper.IntToPtr(stockVal)
		} else {
			e.ErrorHandler(w, r, http.StatusBadRequest, errors.New("invalid stock format"))
			return
		}
	}

	if image := r.URL.Query().Get("image"); image != "" {
		query.Image = helper.StrToPtr(image)
	}

	response, err := h.ProductService.Search(r.Context(), query, pagination)
	if err != nil {
		h.Log.Errorf("failed to search products: %v", err)
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

// @Summary Get product by ID
// @Description Get product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} model.SuccessResponse[model.ProductResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	request := helper.ParseParam(r)
	h.Log.Info("Parsed request parameter: ", request)

	response, err := h.ProductService.GetById(r.Context(), &model.GetProductRequest{ID: request})
	if err != nil {
		h.Log.Errorf("failed to get product by id: %v", err)
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

// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body model.CreateProductRequest true "Product"
// @Success 201 {object} model.SuccessResponse[model.ProductResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security ApiKeyAuth
// @Router /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	request := &model.CreateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		e.ErrorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	response, err := h.ProductService.Create(r.Context(), request)
	if err != nil {
		h.Log.Errorf("failed to create product: %v", err)
		switch {
		case errors.Is(err, e.ErrValidation):
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

// @Summary Update a product
// @Description Update a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body model.UpdateProductRequest true "Product"
// @Success 200 {object} model.SuccessResponse[model.ProductResponse]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security ApiKeyAuth
// @Router /products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	id := &model.DeleteProductRequest{
		ID: helper.ParseParam(r),
	}
	request := &model.UpdateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		e.ErrorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	response, err := h.ProductService.Update(r.Context(), id, request)
	if err != nil {
		h.Log.Errorf("failed to update product: %v", err)
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

// @Summary Delete a product
// @Description Delete a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} model.SuccessResponse[model.DeleteProductRequest]
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security ApiKeyAuth
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		e.ErrorHandler(w, r, http.StatusMethodNotAllowed, e.ErrMethodNotAllowed)
		return
	}

	id := &model.DeleteProductRequest{
		ID: helper.ParseParam(r),
	}

	response, err := h.ProductService.Delete(r.Context(), id)
	if err != nil {
		h.Log.Errorf("failed to delete product: %v", err)
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

// parsePagination is a private helper function to parse pagination parameters
func (h *ProductHandler) parsePagination(r *http.Request) *model.ProductPagination {
	pagination := &model.ProductPagination{
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
