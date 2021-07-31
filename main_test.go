package simplemysql

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error
	db, err = Connect("root:@tcp(127.0.0.1:3306)/demo?timeout=3s&charset=utf8&interpolateParams=true")
	if err != nil {
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	m.Run()
}
