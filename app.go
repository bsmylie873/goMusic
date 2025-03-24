package main

import (
	"database/sql"
	"goMusic/db"
	"log"
	"time"
)

func Setup(dbFile string) {
	d, err := sql.Open("sqlite3", "./"+dbFile+"?_journal=WAL&_timeout=5000")
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(1) // SQLite only supports one writer
	d.SetMaxIdleConns(1)
	d.SetConnMaxLifetime(time.Hour)

	if err = d.Ping(); err != nil {
		panic(err)
	}

	db.DB = d

	err = db.InitDB(dbFile)

	err = db.SeedDB()
	if err != nil {
		log.Fatal("Database seeding failed: ", err)
	}
}

func GetDB() *sql.DB {
	return db.DB
}

func CloseDB() error {
	return db.DB.Close()
}
