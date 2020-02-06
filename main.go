package main

import (
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", info)
	http.HandleFunc("/post/create", post)
	http.HandleFunc("/post/show", show)

	port := ":8082"
	log.Printf("start work on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func info(writer http.ResponseWriter, _ *http.Request) {
	ifaces, err := net.Interfaces()
	if err != nil {
		writer.Write([]byte(err.Error()))
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
			response(writer, err.Error())
			return
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

			response(writer, "your ip is: "+ip.String())
			return
		}
	}

	response(writer, "are you connected to the network?")
}

func post(writer http.ResponseWriter, _ *http.Request) {
	db, err := connect()
	if err != nil {
		response(writer, err.Error())
		return
	}

	st, err := db.Prepare("INSERT INTO person (name) VALUES (md5(random()::text))")
	if err != nil {
		response(writer, err.Error())
		return
	}

	st.Exec([]driver.Value{})
}

func show(writer http.ResponseWriter, _ *http.Request) {
	db, err := connect()
	if err != nil {
		response(writer, err.Error())
		return
	}

	st, err := db.Prepare("Select * from person")
	if err != nil {
		response(writer, err.Error())
		return
	}

	res, err := st.Query()
	if err != nil {
		response(writer, err.Error())
		return
	}

	var data string
	for res.Next() {
		var Id int
		var Name string

		if err = res.Scan(&Id, &Name); err != nil {
			response(writer, err.Error())
			return
		}

		data += fmt.Sprintf("ID: %d, Name: %s \n", Id, Name)
	}

	response(writer, data)
}

func connect() (*sqlx.DB, error) {
	host := "host"
	port := "5432"
	user := "postgres"
	dbname := "postgres"
	password := "1234567890"

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return sqlx.Open("postgres", dns)
}

func response(writer http.ResponseWriter, resp string) {
	writer.Write([]byte(resp))
}
