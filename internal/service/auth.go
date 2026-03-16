package service

import (
	"context"
	"fmt"

	"github.com/htrandev/gophkeeper/internal/contracts"
	"github.com/htrandev/gophkeeper/internal/domain"
)

var _ contracts.AuthService = (*AuthService)(nil)

// AuthServer определяет сервис для работы с пользователем.
type AuthService struct {
	authorizer contracts.Authorizer
	storage    contracts.UserStorage
}

// NewAuth возвращает новый экземпляр AuthService.
func NewAuth(a contracts.Authorizer, s contracts.UserStorage) *AuthService {
	return &AuthService{
		authorizer: a,
		storage:    s,
	}
}

// SignUp регистрирует нового пользователя.
func (s *AuthService) SignUp(ctx context.Context, r domain.AuthorizationRequest) (string, error) {
	// хэшируем пароль
	hash, err := s.authorizer.HashPassword(r.Password)
	if err != nil {
		return "", fmt.Errorf("service: hash password: %w", err)
	}

	// добавляем нового пользователя в базу
	id, err := s.storage.Register(ctx, domain.User{
		Login:        r.Login,
		HashPassword: hash,
	})
	if err != nil {
		return "", fmt.Errorf("service: register new user: %w", err)
	}

	// получаем токен
	token, err := s.authorizer.Token(id)
	if err != nil {
		return "", fmt.Errorf("service: create new token: %w", err)
	}
	return token, nil
}

// SignIn авторизует пользователя.
func (s *AuthService) SignIn(ctx context.Context, r domain.AuthorizationRequest) (string, error) {
	// получаем пользователя из базы
	u, err := s.storage.Login(ctx, r.Login)
	if err != nil {
		return "", fmt.Errorf("service: get user: %w", err)
	}

	// валидируем пароль
	if !s.authorizer.ValidatePassword(u.HashPassword, r.Password) {
		return "", domain.ErrIncorrectPassword
	}

	// получаем токен
	token, err := s.authorizer.Token(u.ID.String())
	if err != nil {
		return "", fmt.Errorf("service: create new token: %w", err)
	}
	return token, nil
}
