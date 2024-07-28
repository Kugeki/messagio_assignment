package pgstore

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"messagio_assignment/internal/domain"
)

func ErrCreateIntoDomain(dbErr error) (domainErr error) {
	var pgErr *pgconn.PgError
	if errors.As(dbErr, &pgErr) {
		// https://www.postgresql.org/docs/16/errcodes-appendix.html
		if pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrAlreadyExists
		}
	}

	if errors.Is(dbErr, pgx.ErrNoRows) {
		return domain.ErrNotCreated
	}

	return dbErr
}

func ErrGetIntoDomain(dbErr error) (domainErr error) {
	if errors.Is(dbErr, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}

	return dbErr
}
