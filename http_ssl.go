package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"time"
)


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
	panic(server.ListenAndServeTLS("ssl/myssl.crt", "ssl/myssl.key"))

}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	// Log request
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

	// Set response Header
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"status": "success"}`)

}
