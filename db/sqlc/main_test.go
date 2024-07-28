package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("connot connect to db, error: ", err)
	}

	testQueries = New(testDB)
	fmt.Println("TestMain Start ... ")
	os.Exit(m.Run()) //m.Run() just start a test runner, start running the unit test and will return exit code to pass to os.Exit
	//os.Exit will
}
