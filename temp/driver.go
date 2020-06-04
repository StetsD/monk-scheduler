package temp

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stetsd/monk-conf"
)

type DbDriver struct {
	dialect string
	db      *sqlx.DB
}

var dbDriver DbDriver

func NewDbDriver(conf config.Config) (*DbDriver, error) {

	if dbDriver.dialect != "" {
		return &dbDriver, nil
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Get(config.DbHost),
		conf.Get(config.DbPort),
		conf.Get(config.DbUser),
		conf.Get(config.DbPass),
		conf.Get(config.DbName),
	)

	db, err := sqlx.Connect("postgres", connString)

	if err != nil {
		return nil, err
	}

	dbDriver = DbDriver{
		dialect: "postgres",
		db:      db,
	}

	return &dbDriver, nil
}

func (dbd *DbDriver) Query(qString string, fields ...interface{}) (*sql.Rows, error) {
	var arguments []interface{} = make([]interface{}, len(fields))

	for i, field := range fields {
		arguments[i] = field
	}

	rows, err := dbd.db.Query(qString,
		arguments...,
	)

	if err != nil {
		return rows, err
	}

	return rows, nil
}
