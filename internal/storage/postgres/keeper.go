package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/htrandev/gophkeeper/internal/domain"
)

// Add добавляет данные пользователю.
func (p *Postgres) Add(ctx context.Context, dto domain.AddRequest) (string, error) {
	var id uuid.UUID

	query := buildAddQuery(dto.Kind)
	args := buildAddArgs(dto)

	err := p.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("repository/add: scan returned id: %w", err)
	}

	return id.String(), nil
}

// Delete удаляет данные пользователя по идентификатору.
func (p *Postgres) Delete(ctx context.Context, dto domain.DeleteRequest) error {
	// запрос на получение идентификатора пользователя, которому принадлежит информация
	ownerQuery := buildOwnerQuery(dto.Kind)

	var id uuid.UUID
	err := p.db.QueryRowContext(ctx, ownerQuery, dto.DataID).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository/delete: scan data owner id: %w", err)
	}

	if dto.UserID != id {
		return domain.ErrPermissionDenied
	}

	query := buildDeleteQuery(dto.Kind)

	if _, err := p.db.ExecContext(ctx, query, dto.DataID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("repository/delete: exec delete query: %w", err)
	}
	return nil
}

// Get получает данные пользователя по идентификатору.
func (p *Postgres) Get(ctx context.Context, dto domain.GetRequest) (domain.Data, error) {
	// запрос на получение идентификатора пользователя, которому принадлежит информация
	ownerQuery := buildOwnerQuery(dto.Kind)

	var id uuid.UUID
	err := p.db.QueryRowContext(ctx, ownerQuery, dto.DataID).Scan(&id)
	if err != nil {
		return domain.Data{}, fmt.Errorf("repository/get: scan data owner id: %w", err)
	}

	if dto.UserID != id {
		return domain.Data{}, domain.ErrPermissionDenied
	}

	return p.queryRow(ctx, dto)
}

func (p *Postgres) GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Info, error) {
	var info []domain.Info

	query := `SELECT id, 'Text' as type, description
		FROM text_data td 
			WHERE td.user_id = $1
		UNION ALL
		SELECT id, 'LogPass' as type, description
		FROM log_pass_data lpd 
			WHERE lpd.user_id = $1
		UNION ALL
		SELECT id, 'File' as type, description
		FROM file_data fd 
			WHERE fd.user_id = $1
		UNION ALL
		SELECT id, 'BankCard' as type, description
		FROM bank_card_data bcd  
			WHERE bcd.user_id = $1;`

	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("repository/getAll: query rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i domain.Info
		var kind string
		if err := rows.Scan(&i.ID, &kind, &i.Description); err != nil {
			return nil, fmt.Errorf("repository/getAll: scan row: %w", err)
		}
		i.Kind = domain.Parse(kind)
		info = append(info, i)
	}

	if rows.Err() != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("repository/getAll: rows err: %w", err)
	}

	return info, nil
}

func (p *Postgres) queryRow(ctx context.Context, dto domain.GetRequest) (domain.Data, error) {
	query := buildGetQuery(dto.Kind)
	row := p.db.QueryRowContext(ctx, query, dto.DataID)
	data := domain.Data{}

	switch dto.Kind {
	case domain.PayloadLogPass:
		logPass := domain.LogPass{}
		if err := row.Scan(&data.ID, &data.Descriptoin, &logPass.Login, &logPass.Password); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Data{}, domain.ErrNotFound
			}
			return domain.Data{}, fmt.Errorf("scan log pass: %w", err)
		}
		data.Kind = domain.PayloadLogPass
		data.LogPass = &logPass
	case domain.PayloadText:
		text := domain.Text{}
		if err := row.Scan(&data.ID, &data.Descriptoin, &text.Text); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Data{}, domain.ErrNotFound
			}
			return domain.Data{}, fmt.Errorf("scan text: %w", err)
		}
		data.Kind = domain.PayloadText
		data.Text = &text
	case domain.PayloadFile:
		file := domain.File{}
		if err := row.Scan(&data.ID, &data.Descriptoin, &file.Name, &file.Content); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Data{}, domain.ErrNotFound
			}
			return domain.Data{}, fmt.Errorf("scan file: %w", err)
		}
		data.Kind = domain.PayloadFile
		data.File = &file
	case domain.PayloadBankCard:
		bankCard := domain.BankCard{}
		if err := row.Scan(&data.ID, &data.Descriptoin, &bankCard.Holder, &bankCard.Number); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.Data{}, domain.ErrNotFound
			}
			return domain.Data{}, fmt.Errorf("scan bank card: %w", err)
		}
		data.Kind = domain.PayloadBankCard
		data.BankCard = &bankCard
	}

	if row.Err() != nil {
		return domain.Data{}, fmt.Errorf("repository/get: queryRow: row err: %w", row.Err())
	}
	return data, nil
}
