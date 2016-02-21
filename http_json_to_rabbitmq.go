package main

import (
	"encoding/json"
	"flag"
	"github.com/streadway/amqp"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	RabbitMQUrl  = "amqp://test:test@rabbit:5672//test"
	RabbitMQueue = "task_queue"
)

type DataStruct struct {
	Url string
	Id  string
}

func main() {

	// Flags
	var addr = flag.String("addr", ":8008", "Port for application, default: :8008")
	flag.Parse()

	// RabbitMQ
	conn, err := amqp.Dial(RabbitMQUrl)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		RabbitMQueue, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		panic(err)
	}

	message := make(chan DataStruct)

	go func() {
		for {
			body, err := json.Marshal(<-message)

			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         []byte(body),
				})

			if err != nil {
				panic(err)
			}

			log.Printf(" [x] Sent %s", body)

		}
	}()

	// Http server
	server := &http.Server{
		Addr:           *addr,
		Handler:        nil,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Routing
	http.Handle("/", rootHandler(message))

	// Start HTTP server
	log.Println("Server listen on " + *addr)
	panic(server.ListenAndServe())

}

func rootHandler(mq chan DataStruct) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Log request
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// Set response Header
		w.Header().Set("Content-Type", "application/json")

		// Check Method
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, `{"status": 405, "message": "Method Not Allowed"}`)
			return
		}

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

		// Send message in channel
		mq <- data

		// Pack data to JSON
		log.Println("Data:", data)
		json.NewEncoder(w).Encode(data)

	})
}
