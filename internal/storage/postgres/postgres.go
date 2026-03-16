package postgres

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Postgres определяет хранилище для пользователей и секретной информации.
type Postgres struct {
	maxRetry int
	db       *sql.DB
}

// New возвращает новый экземпляр Postgres.
func New(maxRetry int, db *sql.DB) *Postgres {
	return &Postgres{maxRetry: maxRetry, db: db}
}

func isUniqueErr(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		return true
	}
	return false
}

func isSerializationFailure(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.SerializationFailure
}

func isPgConnErr(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	}
	return false
}

func isRetryable(err error) bool {
	return isPgConnErr(err) || isSerializationFailure(err)
}
