package service

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/savioruz/bake/internal/domain/entity"
	"github.com/savioruz/bake/internal/domain/model"
	"github.com/savioruz/bake/internal/repository"
	e "github.com/savioruz/bake/pkg/error"
	"github.com/savioruz/bake/pkg/helper"
	"github.com/sirupsen/logrus"
)

type OrderService struct {
	OrderRepository   *repository.OrderRepository
	ProductRepository *repository.ProductRepository
	AddressRepository *repository.AddressRepository
	DB                *sqlx.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
	addressRepo *repository.AddressRepository,
	db *sqlx.DB,
	log *logrus.Logger,
	validate *validator.Validate,
) *OrderService {
	return &OrderService{
		OrderRepository:   orderRepo,
		ProductRepository: productRepo,
		AddressRepository: addressRepo,
		DB:                db,
		Log:               log,
		Validate:          validate,
	}
}

func (s *OrderService) Create(ctx context.Context, request *model.CreateOrderRequest) (*model.SuccessResponse[*model.OrderResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		s.Log.Errorf("validation error for request: %v", err)
		return nil, e.ErrValidation
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		s.Log.Errorf("error beginning transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			s.Log.Errorf("rolling back transaction due to error: %v", err)
			tx.Rollback()
			return
		}
	}()

	address, err := s.AddressRepository.GetByUserID(tx, request.UserID)
	if err != nil {
		s.Log.Errorf("error getting user address: %v", err)
		return nil, err
	}

	product, err := s.ProductRepository.GetByID(tx, request.ProductID)
	if err != nil {
		s.Log.Errorf("error getting product: %v", err)
		return nil, err
	}

	if product.Stock < request.Quantity {
		s.Log.Error("insufficient stock")
		return nil, e.ErrInsufficientStock
	}

	totalPrice := product.Price * float64(request.Quantity)

	now := time.Now()
	order := &entity.Order{
		ID:         uuid.NewString(),
		UserID:     request.UserID,
		ProductID:  request.ProductID,
		AddressID:  address.ID,
		Quantity:   request.Quantity,
		TotalPrice: totalPrice,
		Status:     "PENDING",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.OrderRepository.Create(tx, order); err != nil {
		s.Log.Errorf("error creating order: %v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	orderResponse := &model.OrderResponse{
		ID:         order.ID,
		ProductID:  order.ProductID,
		Quantity:   order.Quantity,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  helper.FormatTime(order.CreatedAt),
		UpdatedAt:  helper.FormatTime(order.UpdatedAt),
		Product:    *product,
		Address:    *address,
	}

	return &model.SuccessResponse[*model.OrderResponse]{
		Data: &orderResponse,
	}, nil
}

func (s *OrderService) GetAll(ctx context.Context, request *model.OrderPagination) (*model.SuccessResponse[[]*model.OrderResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		s.Log.Errorf("validation error for request: %v", err)
		return nil, e.ErrValidation
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		s.Log.Errorf("error beginning transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			s.Log.Errorf("rolling back transaction due to error: %v", err)
			tx.Rollback()
			return
		}
	}()

	orders, total, err := s.OrderRepository.GetAll(tx, request)
	if err != nil {
		s.Log.Errorf("error getting all orders: %v", err)
		return nil, err
	}

	orderResponses := make([]*model.OrderResponse, len(orders))
	for i, order := range orders {
		product, err := s.ProductRepository.GetByID(tx, order.ProductID)
		if err != nil {
			s.Log.Errorf("error getting product for order %s: %v", order.ID, err)
			return nil, err
		}

		address, err := s.AddressRepository.GetByUserID(tx, order.UserID)
		if err != nil {
			s.Log.Errorf("error getting address for order %s: %v", order.ID, err)
			return nil, err
		}

		orderResponses[i] = &model.OrderResponse{
			ID:         order.ID,
			UserID:     order.UserID,
			ProductID:  order.ProductID,
			AddressID:  order.AddressID,
			Quantity:   order.Quantity,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  helper.FormatTime(order.CreatedAt),
			UpdatedAt:  helper.FormatTime(order.UpdatedAt),
			Product:    *product,
			Address:    *address,
		}
	}

	response := model.SuccessResponse[[]*model.OrderResponse]{
		Data: &orderResponses,
		Paginate: &model.Paginate{
			Page:       request.Page,
			Limit:      request.Limit,
			TotalPages: helper.CalculateTotalPages(total, request.Limit),
			TotalItems: total,
		},
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	return &response, nil
}

func (s *OrderService) GetById(ctx context.Context, request *model.GetOrderRequest) (*model.SuccessResponse[*model.OrderResponse], error) {
	if err := s.Validate.Struct(request); err != nil {
		s.Log.Errorf("validation error for request: %v", err)
		return nil, e.ErrValidation
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		s.Log.Errorf("error beginning transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			s.Log.Errorf("rolling back transaction due to error: %v", err)
			tx.Rollback()
			return
		}
	}()

	order, err := s.OrderRepository.GetByID(tx, request.ID)
	if err != nil {
		s.Log.Errorf("error getting order by id: %v", err)
		return nil, err
	}

	product, err := s.ProductRepository.GetByID(tx, order.ProductID)
	if err != nil {
		s.Log.Errorf("error getting product by id: %v", err)
		return nil, err
	}

	address, err := s.AddressRepository.GetByUserID(tx, order.UserID)
	if err != nil {
		s.Log.Errorf("error getting address by user id: %v", err)
		return nil, err
	}

	orderResponse := &model.OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		ProductID:  order.ProductID,
		AddressID:  order.AddressID,
		Quantity:   order.Quantity,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  helper.FormatTime(order.CreatedAt),
		UpdatedAt:  helper.FormatTime(order.UpdatedAt),
		Product:    *product,
		Address:    *address,
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	return &model.SuccessResponse[*model.OrderResponse]{
		Data: &orderResponse,
	}, nil
}
