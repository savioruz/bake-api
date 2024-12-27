package service

import (
	"context"
	"database/sql"
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

type ProductService struct {
	ProductRepository *repository.ProductRepository
	DB                *sqlx.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
}

func NewProductService(
	productRepo *repository.ProductRepository,
	db *sqlx.DB,
	log *logrus.Logger,
	validate *validator.Validate,
) *ProductService {
	return &ProductService{
		ProductRepository: productRepo,
		DB:                db,
		Log:               log,
		Validate:          validate,
	}
}

func (s *ProductService) GetAll(ctx context.Context, request *model.ProductPagination) (*model.SuccessResponse[[]*model.ProductResponse], error) {
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

	products, total, err := s.ProductRepository.GetAll(tx, request)
	if err != nil {
		s.Log.Errorf("error getting all products: %v", err)
		return nil, err
	}

	productResponses := make([]*model.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = &model.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
			CreatedAt:   helper.FormatTime(product.CreatedAt),
			UpdatedAt:   helper.FormatTime(product.UpdatedAt),
		}
	}

	response := model.SuccessResponse[[]*model.ProductResponse]{
		Data: &productResponses,
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

func (s *ProductService) Search(ctx context.Context, query *model.ProductQuery, pagination *model.ProductPagination) (*model.SuccessResponse[[]*model.ProductResponse], error) {
	if err := s.Validate.Struct(query); err != nil {
		s.Log.Errorf("validation error for query: %v", err)
		return nil, e.ErrValidation
	}
	if err := s.Validate.Struct(pagination); err != nil {
		s.Log.Errorf("validation error for pagination: %v", err)
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

	products, total, err := s.ProductRepository.Search(tx, query, pagination)
	if err != nil {
		s.Log.Errorf("error searching products: %v", err)
		return nil, err
	}

	productResponses := make([]*model.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = &model.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
			Image:       product.Image,
			CreatedAt:   helper.FormatTime(product.CreatedAt),
			UpdatedAt:   helper.FormatTime(product.UpdatedAt),
		}
	}

	response := model.SuccessResponse[[]*model.ProductResponse]{
		Data: &productResponses,
		Paginate: &model.Paginate{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalPages: helper.CalculateTotalPages(total, pagination.Limit),
			TotalItems: total,
		},
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	return &response, nil
}

func (s *ProductService) GetById(ctx context.Context, request *model.GetProductRequest) (*model.SuccessResponse[*model.ProductResponse], error) {
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

	var data *entity.Product
	data, err = s.ProductRepository.GetByID(tx, request.ID)
	if err != nil {
		s.Log.Errorf("error getting product by id: %v", err)
		return nil, err
	}

	productResponse := &model.ProductResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Stock:       data.Stock,
		Image:       data.Image,
		CreatedAt:   helper.FormatTime(data.CreatedAt),
		UpdatedAt:   helper.FormatTime(data.UpdatedAt),
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	return &model.SuccessResponse[*model.ProductResponse]{
		Data: &productResponse,
	}, nil
}

func (s *ProductService) Create(ctx context.Context, request *model.CreateProductRequest) (*model.SuccessResponse[*model.ProductResponse], error) {
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

	data := &entity.Product{
		ID:          uuid.NewString(),
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		Image:       request.Image,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.ProductRepository.Create(tx, data); err != nil {
		s.Log.Errorf("error creating product: %v", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	productResponse := &model.ProductResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Stock:       data.Stock,
		Image:       data.Image,
		CreatedAt:   helper.FormatTime(data.CreatedAt),
		UpdatedAt:   helper.FormatTime(data.UpdatedAt),
	}

	return &model.SuccessResponse[*model.ProductResponse]{
		Data: &productResponse,
	}, nil
}

func (s *ProductService) Update(ctx context.Context, id *model.DeleteProductRequest, request *model.UpdateProductRequest) (*model.SuccessResponse[*model.ProductResponse], error) {
	if err := s.Validate.Struct(id); err != nil {
		s.Log.Errorf("validation error for id: %v", err)
		return nil, e.ErrValidation
	}
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

	data := &entity.Product{
		ID:          id.ID,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		Image:       request.Image,
		UpdatedAt:   time.Now(),
	}

	if err := s.ProductRepository.Update(tx, data); err != nil {
		s.Log.Errorf("error updating product: %v", err)
		if err == sql.ErrNoRows {
			return nil, e.ErrNotFound
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	productResponse := &model.ProductResponse{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
		Stock:       data.Stock,
		Image:       data.Image,
		CreatedAt:   helper.FormatTime(data.CreatedAt),
		UpdatedAt:   helper.FormatTime(data.UpdatedAt),
	}

	return &model.SuccessResponse[*model.ProductResponse]{
		Data: &productResponse,
	}, nil
}

func (s *ProductService) Delete(ctx context.Context, request *model.DeleteProductRequest) (*model.SuccessResponse[*model.DeleteProductRequest], error) {
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

	if err := s.ProductRepository.Delete(tx, request.ID); err != nil {
		s.Log.Errorf("error deleting product: %v", err)
		if err == sql.ErrNoRows {
			return nil, e.ErrNotFound
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	return &model.SuccessResponse[*model.DeleteProductRequest]{
		Data: &request,
	}, nil
}
