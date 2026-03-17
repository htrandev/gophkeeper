package contracts

import (
	"context"

	"github.com/google/uuid"

	"github.com/htrandev/gophkeeper/internal/domain"
)

// UserStorage предоставляет интерфейс взаимодействия с хранилищем пользователей.
type UserStorage interface {
	Register(ctx context.Context, user domain.User) (string, error)
	Login(ctx context.Context, login string) (domain.User, error)
}

// DataStorage предоставляет интерфейс взаимодействия с хранилищем секретной информации.
type DataStorage interface {
	Add(ctx context.Context, dto domain.AddRequest) (string, error)
	Delete(ctx context.Context, dto domain.DeleteRequest) error
	Get(ctx context.Context, dto domain.GetRequest) (domain.Data, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Info, error)
}
