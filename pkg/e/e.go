package e

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows
var (
	ErrPassword     = fmt.Errorf("incorrect password")
	ErrExpiredToken = fmt.Errorf("token has expired")
	ErrInvalidToken = fmt.Errorf("token is invalid")
)
var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
