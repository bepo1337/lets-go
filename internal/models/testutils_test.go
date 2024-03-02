package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDb(t *testing.T) *sql.DB {
	//Create new SQL Connection (*sql.DB)
	dataSourceString := "root:admin@/test_snippetbox?parseTime=true&multiStatements=true"
	db, err := sql.Open("mysql", dataSourceString)
	if err != nil {
		t.Fatal(err)
	}
	// Execute setup.sql
	script := readSqlFromFile(t, "./testdata/setup.sql")

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
	// Register cleanup function (teardown.sql execution) and close connection pool
	t.Cleanup(func() {
		script = readSqlFromFile(t, "./testdata/teardown.sql")
		_, err = db.Exec(script)
		if err != nil {
			t.Fatal(err)
		}
		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}

func readSqlFromFile(t *testing.T, fileName string) string {
	script, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return string(script)
}
