package main

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"log"
	"net/http"
	"time"
)

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func main() {
	router := routing.New()

	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	// serve RESTful APIs
	api := router.Group("/api")

	api.Use(
		content.TypeNegotiator(content.JSON),
	)

	api.Get("/users", func(c *routing.Context) error { return usersGet(c) })
	api.Post("/users", func(c *routing.Context) error { return usersPost(c) })
	api.Put(`/users/<id:\d+>`, func(c *routing.Context) error { return usersPut(c) })

	// Http server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Router
	http.Handle("/", router)

	// Start HTTP server
	log.Println("Server listen on 8080")
	panic(server.ListenAndServe())
}

func usersGet(c *routing.Context) error {
	var users []User

	user := User{
		Id:   123,
		Name: "User",
	}
	users = append(users, user)
	return c.Write(users)
}

func usersPost(c *routing.Context) error {
	return c.Write("create a new user")
}

func usersPut(c *routing.Context) error {
	return c.Write("update user " + c.Param("id"))
}
