package core_pgx_pool

import (
	"errors"

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
		if errors.Is(err, pgx.ErrNoRows) {
			return core_postgres_pool.ErrNoRows
		}

		return err
	}

	return nil
}

type CmdTag struct {
	pgconn.CommandTag
}
