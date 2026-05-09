package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dataPath string) {
	var err error
	DB, err = sql.Open("sqlite", dataPath)
	if err != nil {
		log.Fatalf("Gagal buka database: %v", err)
	}

	DB.SetMaxOpenConns(1)

	if err = DB.Ping(); err != nil {
		log.Fatalf("Gagal ping database: %v", err)
	}

	runSchema()
	log.Println("Database siap:", dataPath)
}

func runSchema() {
	_, filename, _, _ := runtime.Caller(0)
	schemaPath := filepath.Join(filepath.Dir(filename), "schema.sql")

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Gagal baca schema.sql: %v", err)
	}

	if _, err := DB.Exec(string(schema)); err != nil {
		log.Fatalf("Gagal jalankan schema: %v", err)
	}
}
