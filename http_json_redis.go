package main

import (
	"encoding/json"
	"flag"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"net/http"
	"strings"
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

	redisPool := redis.Pool{
		MaxIdle:   2,
		MaxActive: 100, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	defer redisPool.Close()

	// Http server
	server := &http.Server{
		Addr:           *addr,
		Handler:        nil,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Routing
	http.Handle("/", rootHandler(redisPool))

	// Start HTTP server
	log.Println("Server listen on " + *addr)
	panic(server.ListenAndServe())

}

func rootHandler(redisPool redis.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log request
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// Set response Header
		w.Header().Set("Content-Type", "application/json")

                conn := redisPool.Get()
                defer conn.Close()


		switch r.Method {
		case "GET":
			{
				data, _ := redis.Strings(conn.Do("MGET", "testKey"))
				io.WriteString(w, data[0])
			}
		case "POST":
			{
				// Check Content-Type
				if strings.ToLower(r.Header.Get("Content-Type")) != "application/json" {
					w.WriteHeader(http.StatusUnsupportedMediaType)

					io.WriteString(w, `{"status": "415", "errors": ["request": "Unsupported Media Type"]}`)
					return
				}

				// Decode JSON
				var data DataStruct
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&data)
				if err != nil {
					log.Println(err)
					http.Error(w, http.StatusText(500), 500)
					return
				}

				if _, err = conn.Do("SET", "testKey", data); err != nil {
					log.Println(err)
					http.Error(w, http.StatusText(500), 500)
                                        return
				}

				// Pack data to JSON
				log.Println("Data:", data)
				json.NewEncoder(w).Encode(data)
			}
		default:
			{
				w.WriteHeader(http.StatusMethodNotAllowed)
				io.WriteString(w, `{"status": 405, "message": "Method Not Allowed"}`)
				return
			}
		}

	})
}
