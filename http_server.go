package main

import (
        "flag"
        "log"
        "net/http"
        "time"
)


func main() {

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
    log.Println("bubu")
}

