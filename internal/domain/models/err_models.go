package models

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

const (
	UniqueViolationErr = pq.ErrorCode("23505")
)

func PostgresHandleError(err error) error {
	if err == nil {
		return err
	}
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr
	}
	return err
}
