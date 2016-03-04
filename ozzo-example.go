package main

import (
        "errors"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/auth"
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

	// Auth handler
        router.Post("/api/auth", func(c *routing.Context) error { return authHandler(c) })

	// serve RESTful APIs
	api := router.Group("/api")

	api.Use(
		content.TypeNegotiator(content.JSON),
		auth.Bearer(func(c *routing.Context, token string) (auth.Identity, error) { return jwtMiddleware(c, token) }),
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

func jwtMiddleware(c *routing.Context, token string) (auth.Identity, error) {
        // JWT processing must be here
	if token == "Aladdin:open sesame" {
		return auth.Identity("demo"), nil
	}
	return nil, errors.New("invalid credential")
}


func authHandler(c *routing.Context) error {
	return c.Write("password")
}

func usersGet(c *routing.Context) error {
        log.Println(c.Get(auth.User))
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
