package simplemysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rnben/simplemysql/utils"
)

type DB struct {
	*sql.DB
}

type session struct {
	db      *DB
	table   string
	columns string
	where   []string
	order   string
	limit   int
	offset  int
	sql     string
	fields  map[string]string
}

func ConnectMySQL(dsn string) (*DB, error) {
	db, err := sql.Open("mysql", dsn)
	return &DB{db}, err
}

func (db *DB) NewSession() *session {
	return &session{
		db:     db,
		fields: make(map[string]string),
	}
}

func (sess *session) Select(columns ...string) *session {
	if len(columns) == 0 {
		sess.columns = "*"
		return sess
	}

	cc := make([]string, 0, len(columns))
	for _, c := range columns {
		if c == "*" {
			cc = append(cc, c)
			break
		}

		cc = append(cc, utils.Quote(c))
	}

	sess.columns = strings.Join(cc, ",")

	return sess
}

func (sess *session) Table(table string) *session {
	sess.table = utils.Quote(table)
	return sess
}

func (sess *session) Where(k string, v interface{}) *session {
	if k == "" {
		return sess
	}

	if v == nil {
		sess.where = append(sess.where, k)
		return sess
	}

	var where string
	switch value := v.(type) {
	case int, int64:
		where = fmt.Sprintf("`%s` = %d", k, value)
	case string:
		where = fmt.Sprintf("`%s` = '%s'", k, value)
	}

	sess.where = append(sess.where, where)

	return sess
}

func (sess *session) Limit(limit int) *session {
	sess.limit = limit
	return sess
}

func (sess *session) Offset(offset int) *session {
	sess.offset = offset
	return sess
}

func (sess *session) Order(order string) *session {
	sess.order = order
	return sess
}

// Query sql injection warn!!
func (sess *session) Query(dest interface{}) error {
	err := sess.newValidate("list").validate()
	if err != nil {
		return err
	}

	sess.fields, err = getSliceFields(dest)
	if err != nil {
		return err
	}

	rows, err := sess.db.Query(sess.build())
	if err != nil {
		return err
	}
	defer rows.Close()

	err = sess.formatOutput(rows, dest)

	return err
}

// RawQuery raw sql query
func (sess *session) RawQuery(dest interface{}, sql string, args ...interface{}) error {
	var err error

	sess.fields, err = getSliceFields(dest)
	if err != nil {
		return err
	}

	rows, err := sess.db.DB.Query(sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = sess.formatOutput(rows, dest)

	return err
}

type validate interface {
	validate() error
}

func (sess *session) newValidate(method string) validate {
	return &listValidate{sess}
}

type listValidate struct {
	*session
}

func (list *listValidate) validate() error {
	if list.table == "" {
		return errors.New("invalid table")
	}

	if list.columns == "" {
		return errors.New("invalid columns")
	}

	if list.limit == 0 {
		return errors.New("invalid limit")
	}

	return nil
}

func (sess *session) build() string {
	sess.sql = fmt.Sprintf("select %s from %s ", sess.columns, sess.table)

	if len(sess.where) != 0 {
		sess.sql += " where "
		sess.sql += strings.Join(sess.where, "and")
	}

	if sess.order != "" {
		sess.sql += " order by " + sess.order + " "
	}

	if sess.limit != 0 {
		sess.sql += " limit " + strconv.Itoa(sess.limit) + " "
	}

	if sess.offset != 0 {
		sess.sql += " offset " + strconv.Itoa(sess.offset) + " "
	}

	return sess.sql
}
