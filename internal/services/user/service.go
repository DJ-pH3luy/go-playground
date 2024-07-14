package user

import (
	"context"
	"fmt"
	"time"

	"github.com/dj-ph3luy/go-playground/internal/db"
	"github.com/dj-ph3luy/go-playground/internal/dto"
	"github.com/dj-ph3luy/go-playground/internal/entities"
	"github.com/dj-ph3luy/go-playground/internal/views"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo db.IRepository[entities.User]
	secret []byte
}

func New(repo db.IRepository[entities.User]) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, dto dto.CreateUser) (uint, error) {
	entity := entities.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
	err := s.repo.Create(ctx, &entity)
	if err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (s *Service) Update(ctx context.Context, dto dto.UpdateUser) (uint, error) {
	entity, err := s.repo.GetById(ctx, dto.Id)
	if err != nil {
		return 0, err
	}
	entity.Password = dto.Password
	err = s.repo.Update(ctx, &entity)
	if err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (s *Service) Login(ctx context.Context, name string, password string) (entities.User, error) {
	user, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return user, fmt.Errorf("wrong username or password")
	}

	err = verifyPassword(password, user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *Service) GetById(ctx context.Context, id string) (entities.User, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service) GetByName(ctx context.Context, name string) (entities.User, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *Service) GetAll(ctx context.Context) ([]entities.User, error) {
	return s.repo.GetAll(ctx)
}

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) GenerateJWT(user views.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := views.UserClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
