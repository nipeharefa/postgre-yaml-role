package lib

import (
	"context"
	"database/sql"
)

func Exec(ctx context.Context, db *sql.DB, query string) error {

	db.ExecContext(ctx, query)
	return nil
}
