package contracts

import "database/sql"

type PgDriver interface {
	Query(qString string, fields ...interface{}) (*sql.Rows, error)
}
