package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

type DataStruct struct {
	Url string
	Id  string
}

func main() {

	// Flags
	var addr = flag.String("addr", ":8008", "Port for application, default: :8008")
	flag.Parse()

	// Http server
	server := &http.Server{
		Addr:           *addr,
		Handler:        nil,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Routing
	http.HandleFunc("/", rootHandler)

	// Start HTTP server
	log.Println("Server listen on " + *addr)
	panic(server.ListenAndServe())

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON
	var data DataStruct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	log.Println(r.Method, data)

}
