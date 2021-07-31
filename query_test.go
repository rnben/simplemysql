package simplemysql

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rnben/log/v2"
)

func TestSQLX(t *testing.T) {
	var pageSize = 500

	for pageIndex := 0; ; pageIndex++ {
		var identities []MyTable

		if err := db.Select(&identities,
			"select * from my_table order by id asc limit ? offset ?", pageSize, pageIndex*pageSize); err != nil {
			t.Fatal(err)
		}

		if len(identities) == 0 {
			break
		}

		t.Logf("pageIndex:%d, count:%d", pageIndex, len(identities))
	}
}

func TestRawQueryContext(t *testing.T) {
	var pageSize = 500

	var result []MyTable

	for pageIndex := 0; ; pageIndex++ {
		var tmp []MyTable
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3) // funny
		rows, err := db.QueryContext(ctx, "select id, name, create_time from my_table order by id asc limit ? offset ?", pageSize, pageIndex*pageSize)
		if err != nil {
			t.Fatal(err)
		}

		for rows.Next() {
			var tableItem MyTable

			rows.Scan(&tableItem.ID, &tableItem.Name, &tableItem.CreateTime)
			time.Sleep(time.Millisecond * 100) // rows.Err()  will receive "context deadline exceeded" error

			tmp = append(tmp, tableItem)
		}
		if err = rows.Err(); err != nil {
			t.Error(err)
		}
		if err = rows.Close(); err != nil {
			log.Error(err)
		}

		if len(tmp) == 0 {
			break
		}

		t.Log(pageIndex, len(tmp))

		result = append(result, tmp...)
	}

	t.Log("total", len(result))

}

func TestQueryRow(t *testing.T) {
	var tableItem MyTable
	err := db.Get(&tableItem, "select * from my_table where id = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tableItem)
}

func TestRowsScanInterface(t *testing.T) {
	var result [][]interface{}

	rows, err := db.Query("select id, name from my_table limit ?", 10)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	for rows.Next() {
		colVal := make([]interface{}, len(cols))
		colValPointer := make([]interface{}, len(cols))
		for i, _ := range colVal {
			colValPointer[i] = &colVal[i]
		}

		rows.Scan(colValPointer...)

		result = append(result, colVal)
	}
	if err = rows.Err(); err != nil {
		log.Error(err)
	}

	t.Log(len(result))
}

func TestInsert(t *testing.T) {
	for i := 0; i < 1<<17; i++ {
		ret, err := db.Exec("insert into my_table(name) values(?)", strconv.FormatUint(uint64(rand.Int63()), 16))
		if err != nil {
			t.Fatal(err)
		}

		t.Log(ret.LastInsertId())
	}
}

func BenchmarkUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuid.New().String()
	}
}
func BenchmarkFormatUnit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strconv.FormatUint(uint64(rand.Int63()), 16)
	}
}

type MyTable struct {
	ID         int64  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	CreateTime string `db:"create_time"`
}
