package main

import (
	"database/sql/driver"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/", info)
	http.HandleFunc("/post/create", post)

	port := ":8080"
	log.Printf("start work on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func info(writer http.ResponseWriter, request *http.Request) {
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
			writer.Write([]byte(err.Error()))
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

			writer.Write([]byte("your ip is: " + ip.String()))
			return
		}
	}

	writer.Write([]byte("are you connected to the network?"))
}

func post(writer http.ResponseWriter, request *http.Request) {
	host := "localhost"
	port := "5432"
	user := "postgres"
	dbname := "postgres"
	password := "yourPassword"

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := pq.Open(dns)
	if err != nil {
		response(writer, err.Error())
	}

	st, er := db.Prepare("INSERT INTO person (id, name) VALUES (1, md5(random()::text))")
	if er != nil {
		response(writer, er.Error())
	}

	st.Exec([]driver.Value{})
}

func response(writer http.ResponseWriter, resp string) {
	writer.Write([]byte(resp))
}
