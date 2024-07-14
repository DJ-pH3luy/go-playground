package services

import (
	"context"

	"github.com/dj-ph3luy/go-playground/internal/dto"
	"github.com/dj-ph3luy/go-playground/internal/entities"
	"github.com/dj-ph3luy/go-playground/internal/views"
)
	

type IUserService interface {
	Create(ctx context.Context, dto dto.CreateUser) (uint, error)
	Update(ctx context.Context, dto dto.UpdateUser) (uint, error)
	Login(ctx context.Context, name string, password string) (entities.User, error)
	GetAll(ctx context.Context) ([]entities.User, error)
	GetById(ctx context.Context, id string) (entities.User, error)
	GetByName(ctx context.Context, name string) (entities.User, error) 
	GenerateJWT(user views.User) (string, error)
}