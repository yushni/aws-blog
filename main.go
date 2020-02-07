package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net"
)

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}

	lambda.Start(do(db))
}

func do(db *DB) func() (string, error) {
	var err error

	if err = db.Create(); err != nil {
		return func() (s string, err error) {
			return "", err
		}
	}

	res, err := db.All()
	if err != nil {
		return func() (s string, err error) {
			return "", err
		}
	}

	ip, err := MyIP()
	if err != nil {
		return func() (s string, err error) {
			return "", err
		}
	}

	return func() (s string, err error) {
		return "IP: " + ip + " \n" + res, nil
	}
}

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

func MyIP() (ip string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ip, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return ip, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ip.String(), nil
		}
	}

	return ip, errors.New("you are not connected to the network")
}
