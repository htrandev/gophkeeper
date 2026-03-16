package domain

import "github.com/google/uuid"

// AuthorizationRequest определяет формат запроса для регистрации и авторизации пользователя.
type AuthorizationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// AddRequest определяет формат запроса для добавления данных пользователем.
type AddRequest struct {
	Token       string
	Description string
	UserID      uuid.UUID
	Kind        PayloadKind
	LogPass     *LogPass
	Text        *Text
	File        *File
	BankCard    *BankCard
}

// DeleteRequest определяет формат запроса для удаления данных пользователя.
type DeleteRequest struct {
	Token  string
	DataID string
	Kind   PayloadKind
	UserID uuid.UUID
}

// GetRequest определяет формат запроса для получения данных пользователя.
type GetRequest struct {
	Token  string
	DataID string
	Kind   PayloadKind
	UserID uuid.UUID
}

type GetAllRequest struct{
	Token  string
	UserID uuid.UUID
}