package simplemysql

import (
	"os"
	"testing"
)

var db *DB

func TestMain(m *testing.M) {
	var err error

	db, err = ConnectMySQL("root:@tcp(127.0.0.1:3306)/demo?charset=utf8&interpolateParams=true")
	if err != nil {
		os.Exit(1)
	}

	m.Run()
}
func TestQueryBuild(t *testing.T) {
	var result []*Identity
	err := db.NewSession().Select("id").Table("my_table").Where("id", 1).Order("id desc").Limit(1).Offset(0).Query(&result)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatal("null result")
	}

	t.Log(result[0].ID)
}

func TestQueryWhere(t *testing.T) {
	var result []*Identity
	err := db.NewSession().Select("id").Table("my_table").Where("id > 0", nil).Order("id desc").Limit(1).Offset(1).Query(&result)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatal("null result")
	}

	t.Log(result[0].ID)
}

func TestRawQuery(t *testing.T) {
	var result []*Identity
	err := db.NewSession().RawQuery(&result, "select * from my_table where id > ? limit ?;", 0, 5)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) == 0 {
		t.Fatal("null result")
	}

	t.Log(result[0].ID)
}

type Identity struct {
	ID         int64  `json:"id"`
	IDCardHash string `json:"idcard_hash"`
	CreateTime string `json:"create_time"`
}
