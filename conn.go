package simplemysql

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/rnben/log/v2"
)

func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Errorf("connect DB failed, err:%v\n", err)
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	return db, nil
}
