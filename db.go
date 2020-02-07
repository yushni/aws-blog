package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB() (*DB, error) {
	host := "rds-test.cnhtbv3maxil.eu-central-1.rds.amazonaws.com"
	port := "5432"
	user := "postgres"
	dbname := "postgres"
	password := "1234567890"

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (db *DB) Create() error {
	st, err := db.Prepare("INSERT INTO person (name) VALUES (md5(random()::text))")
	if err != nil {
		return err
	}

	if _, err = st.Exec(); err != nil {
		return err
	}

	return nil
}

func (db *DB) All() (posts string, err error) {
	st, err := db.Prepare("Select * from person")
	if err != nil {
		return posts, err
	}

	res, err := st.Query()
	if err != nil {
		return posts, err
	}

	for res.Next() {
		var Id int
		var Name string

		if err = res.Scan(&Id, &Name); err != nil {
			return posts, err
		}

		posts += fmt.Sprintf("ID: %d, Name: %s \n", Id, Name)
	}

	return posts, err
}
