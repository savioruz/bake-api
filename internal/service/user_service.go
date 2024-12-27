package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/savioruz/bake/internal/domain/entity"
	"github.com/savioruz/bake/internal/domain/model"
	"github.com/savioruz/bake/internal/repository"
	e "github.com/savioruz/bake/pkg/error"
	"github.com/savioruz/bake/pkg/helper"
	"github.com/savioruz/bake/pkg/jwt"
	"github.com/savioruz/bake/pkg/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository    *repository.UserRepository
	AddressRepository *repository.AddressRepository
	DB                *sqlx.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	JWTService        jwt.JWTService
}

func NewUserService(
	userRepo *repository.UserRepository,
	addressRepo *repository.AddressRepository,
	db *sqlx.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	jwtService jwt.JWTService,
) *UserService {
	return &UserService{
		UserRepository:    userRepo,
		AddressRepository: addressRepo,
		DB:                db,
		Log:               log,
		Validate:          validate,
		JWTService:        jwtService,
	}
}

func (s *UserService) CreateUser(ctx context.Context, request *model.UserRegisterRequest) (*model.UserResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, e.ErrValidation
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	existingUser, err := s.UserRepository.GetByEmail(tx, request.Email)
	if err == nil {
		s.Log.Errorf("user already exists: %v", existingUser)
		return nil, e.ErrUserExists
	} else if err != sql.ErrNoRows {
		s.Log.Errorf("error getting user by email: %v", err)
		return nil, err
	}

	first := &entity.User{}
	err = s.UserRepository.GetFirst(tx, first)

	// Role
	var role string
	switch {
	case errors.Is(err, sql.ErrNoRows):
		role = "admin"
	case err != nil:
		s.Log.Errorf("error getting first user: %v", err)
		return nil, err
	default:
		role = "user"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Errorf("error hashing password: %v", err)
		return nil, err
	}

	userID := uuid.NewString()
	data := &entity.User{
		ID:        userID,
		Email:     request.Email,
		Password:  string(hashedPassword),
		Name:      request.Name,
		Phone:     request.Phone,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.UserRepository.Create(tx, data); err != nil {
		s.Log.Errorf("error creating user: %v", err)
		return nil, err
	}

	var addressResp *model.AddressResponse
	if request.Address != nil {
		address := &entity.Address{
			ID:          uuid.NewString(),
			UserID:      userID,
			AddressLine: request.Address.AddressLine,
			City:        request.Address.City,
			State:       request.Address.State,
			PostalCode:  request.Address.PostalCode,
			Country:     request.Address.Country,
		}

		if err := s.AddressRepository.Create(tx, address); err != nil {
			s.Log.Errorf("error creating address: %v", err)
			return nil, err
		}

		addressResp = &model.AddressResponse{
			ID:          address.ID,
			UserID:      address.UserID,
			AddressLine: address.AddressLine,
			City:        address.City,
			State:       address.State,
			PostalCode:  address.PostalCode,
			Country:     address.Country,
			CreatedAt:   helper.FormatTime(address.CreatedAt),
			UpdatedAt:   helper.FormatTime(address.UpdatedAt),
		}
	}

	if err := tx.Commit(); err != nil {
		s.Log.Errorf("error committing transaction: %v", err)
		return nil, err
	}

	return &model.UserResponse{
		ID:        data.ID,
		Email:     data.Email,
		Name:      data.Name,
		Phone:     data.Phone,
		Role:      data.Role,
		CreatedAt: helper.FormatTime(data.CreatedAt),
		UpdatedAt: helper.FormatTime(data.UpdatedAt),
		Address:   addressResp,
	}, nil
}

func (s *UserService) Login(ctx context.Context, request *model.UserLoginRequest) (*model.TokenResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, e.ErrValidation
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	data, err := s.UserRepository.GetByEmail(tx, request.Email)
	if err != nil {
		s.Log.Errorf("error getting user by email: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserNotFound
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(request.Password)); err != nil {
		s.Log.Errorf("error comparing password: %v", err)
		return nil, e.ErrCredential
	}

	accessToken, err := s.JWTService.GenerateAccessToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("error generating access token: %v", err)
		return nil, err
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(data.ID, data.Email, data.Role)
	if err != nil {
		s.Log.Errorf("error generating refresh token: %v", err)
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, request *model.RefreshTokenRequest) (*model.TokenResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		return nil, e.ErrValidation
	}

	claims, err := s.JWTService.ValidateToken(request.RefreshToken)
	if err != nil {
		s.Log.Errorf("failed to validate token: %v", err)
		return nil, e.ErrCredential
	}

	accessToken, err := s.JWTService.GenerateAccessToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate access token: %v", err)
		return nil, err
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		s.Log.Errorf("failed to generate refresh token: %v", err)
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Me(ctx context.Context) (*model.UserResponse, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == "" {
		return nil, e.ErrUnauthorized
	}

	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	data, err := s.UserRepository.GetByID(tx, userID)
	if err != nil {
		return nil, err
	}

	address, err := s.AddressRepository.GetByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	addressResp := &model.AddressResponse{
		ID:          address.ID,
		UserID:      address.UserID,
		AddressLine: address.AddressLine,
		City:        address.City,
		State:       address.State,
		PostalCode:  address.PostalCode,
		Country:     address.Country,
		CreatedAt:   helper.FormatTime(address.CreatedAt),
		UpdatedAt:   helper.FormatTime(address.UpdatedAt),
	}

	return &model.UserResponse{
		ID:        data.ID,
		Email:     data.Email,
		Name:      data.Name,
		Phone:     data.Phone,
		Role:      data.Role,
		CreatedAt: helper.FormatTime(data.CreatedAt),
		UpdatedAt: helper.FormatTime(data.UpdatedAt),
		Address:   addressResp,
	}, nil
}
