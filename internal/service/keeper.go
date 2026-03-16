package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/htrandev/gophkeeper/internal/contracts"
	"github.com/htrandev/gophkeeper/internal/domain"
)

var _ contracts.KeeperService = (*KeeperService)(nil)

// GophkeeperServer определяет сервис для работы с секретными данными.
type KeeperService struct {
	authorizer contracts.Authorizer
	storage    contracts.DataStorage
}

// NewKeeper возвращает новый экземпляр KeeperService.
func NewKeeper(a contracts.Authorizer, s contracts.DataStorage) *KeeperService {
	return &KeeperService{authorizer: a, storage: s}
}

// Add добавляет данные к пользователю.
func (s *KeeperService) Add(ctx context.Context, dto domain.AddRequest) (string, error) {
	// если данные содержат секретную информацию, то хэшируем ее
	if dto.Kind == domain.PayloadLogPass {
		hash, err := s.authorizer.HashPassword(dto.LogPass.Password)
		if err != nil {
			return "", fmt.Errorf("service/add: hash reseived password: %w", err)
		}
		dto.LogPass.Password = hash
	}

	// добавляем данные в базу
	id, err := s.storage.Add(ctx, dto)
	if err != nil {
		return "", fmt.Errorf("service/add: add to storage: %w", err)
	}
	return id, nil
}

// Delete удаляет данные у пользователя по идентификатору.
func (s *KeeperService) Delete(ctx context.Context, dto domain.DeleteRequest) error {
	if err := s.storage.Delete(ctx, dto); err != nil {
		return fmt.Errorf("service/delete: %w", err)
	}
	return nil
}

// Get получает данные пользователя по идентификатору.
func (s *KeeperService) Get(ctx context.Context, dto domain.GetRequest) (domain.Data, error) {
	data, err := s.storage.Get(ctx, dto)
	if err != nil {
		return domain.Data{}, fmt.Errorf("service/get: %w", err)
	}
	return data, nil
}

// GetAll получает все данные пользователя.
func (s *KeeperService) GetAll(ctx context.Context, userId uuid.UUID) ([]domain.Info, error) {
	info, err := s.storage.GetAll(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("service/getAll: %w", err)
	}
	return info, nil
}
