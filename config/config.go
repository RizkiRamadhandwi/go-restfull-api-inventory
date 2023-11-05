package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "inventory"
)

func ConnectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)                  // jumlah koneksi minimal yang dibuat
	db.SetMaxOpenConns(100)                 // jumlah koneksi maksimal yang dibuat
	db.SetConnMaxIdleTime(5 * time.Minute)  // berapa lama koneksi yang sudah tidak digunakan akan dihapus
	db.SetConnMaxLifetime(60 * time.Minute) // berapa lama koneksi boleh digunakan

	return db

}
