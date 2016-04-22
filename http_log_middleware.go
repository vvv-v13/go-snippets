package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		lrw := NewLoggingResponseWriter(w)
		handler.ServeHTTP(lrw, r)
		elapsed := float64(time.Now().Sub(startTime).Nanoseconds()) / 1e6
		statusCode := lrw.statusCode
		log.Printf(`[%s] [%.3fms] %s %s %s %d`, r.RemoteAddr, elapsed, r.Method, r.URL.Path, r.Proto, statusCode)
	})

}

func main() {

	// Flags
	var addr = flag.String("addr", ":8008", "Port for application, default: :8008")
	flag.Parse()

	// Http server
	server := &http.Server{
		Addr:           *addr,
		Handler:        Logger(http.DefaultServeMux),
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
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"status": 200}`)

}
