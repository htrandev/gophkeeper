package contracts

import (
	"context"

	"github.com/google/uuid"
	"github.com/htrandev/gophkeeper/internal/domain"
)

// AuthService предоставляет интерфейс взаимодействия с сервисом авторизации и аутентификации.
type AuthService interface {
	SignUp(ctx context.Context, r domain.AuthorizationRequest) (string, error)
	SignIn(ctx context.Context, r domain.AuthorizationRequest) (string, error)
}

// KeeperService предоставляет интерфейс взаимодействия с сервисом хранения секретной информации.
type KeeperService interface {
	Add(ctx context.Context, dto domain.AddRequest) (string, error)
	Delete(ctx context.Context, dto domain.DeleteRequest) error
	Get(ctx context.Context, dto domain.GetRequest) (domain.Data, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Info, error)
}
