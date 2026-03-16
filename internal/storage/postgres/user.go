package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/htrandev/gophkeeper/internal/domain"
)

// Register добавляет нового пользователя в базу.
func (p *Postgres) Register(ctx context.Context, u domain.User) (string, error) {
	var id uuid.UUID

	query := `INSERT INTO users (login, password)
		VALUES($1, $2)
	RETURNING id;`

	err := p.db.QueryRowContext(ctx, query,
		u.Login,
		u.HashPassword,
	).Scan(&id)
	if err != nil {
		if isUniqueErr(err) {
			return "", domain.ErrNotUniqueLogin
		}
		return "", fmt.Errorf("repository/register: scan returned id: %w", err)
	}

	return id.String(), nil
}

// Login получает информацию о пользователе из базы.
func (r *Postgres) Login(ctx context.Context, login string) (domain.User, error) {
	var u domain.User

	query := `SELECT 
		id, login, password
	FROM users
	WHERE login = $1;`

	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&u.ID,
		&u.Login,
		&u.HashPassword,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("repository/login: scan: %w", err)
	}

	return u, nil
}
