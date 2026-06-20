package core_pgx_pool

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/roman-styazhkin/golang-todoapp/internal/core/repository/postgres/pool"
)

type Rows struct {
	pgx.Rows
}

type Row struct {
	pgx.Row
}

func (r Row) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return mapErrors(err)
	}

	return nil
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyCode = "23503"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyCode {
			return fmt.Errorf(
				"%v, %w", err,
				core_postgres_pool.ErrViolatesForeignKey,
			)
		}
	}

	return fmt.Errorf(
		"%v, %w",
		err,
		core_postgres_pool.ErrUnknown,
	)
}

type CmdTag struct {
	pgconn.CommandTag
}
