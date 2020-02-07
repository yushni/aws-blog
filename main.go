package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	db, err := NewDB()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", do(db))

	port := ":8081"
	log.Printf("start work on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func do(db *DB) func(http.ResponseWriter, *http.Request) {
	var err error
	respErr := func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(err.Error()))
	}

	if err = db.Create(); err != nil {
		return respErr
	}

	res, err := db.All()
	if err != nil {
		return respErr
	}

	ip, err := MyIP()
	if err != nil {
		return respErr
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("IP: " + ip + " \n" + res))
	}
}
