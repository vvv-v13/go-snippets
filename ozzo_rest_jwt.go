package main

import (
	"errors"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	"github.com/vvv-v13/ozzo-jwt"
	"log"
	"net/http"
	"time"
)

func main() {

	jwtConfig := jwt.JWTConfig{
		Alg:    "HS256",
		Secret: "super_secret",
		Expires: time.Minute * 120,
	}

	router := routing.New()
	router.Use(
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON),
		fault.Recovery(log.Printf),
	)

	router.Post("/api/auth", func(c *routing.Context) error { return authHandler(c, jwtConfig) })

	api := router.Group("/api")
	api.Use(
		jwt.JWT(func(c *routing.Context, payload jwt.JWTPayload) (jwt.Payload, error) {
			return identity(c, payload)
		}, jwtConfig),
	)

	api.Get("/posts", func(c *routing.Context) error { return posts(c) })

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

func authHandler(c *routing.Context, jwtConfig jwt.JWTConfig) error {

	payload := make(jwt.JWTPayload)
	payload["id"] = 1
	payload["role"] = "user"

	token, err := jwt.CreateToken(jwtConfig, payload)
	if err != nil {
		return routing.NewHTTPError(http.StatusInternalServerError)
	}

	data := map[string]string{
		"token": token,
	}

	return c.Write(data)
}

func identity(c *routing.Context, payload jwt.JWTPayload) (jwt.Identity, error) {
	if id, ok := payload["id"]; ok {
		return jwt.Identity(id), nil
	}
	return nil, errors.New("invalid credential")
}

type Post struct {
        Id      int64  `json:"id"`
        Message string `json:"message"`
}

func posts(c *routing.Context) error {
	log.Println("User id:", c.Get(jwt.User))
	var posts []Post

	post := Post{
		Id:      1,
		Message: "Message",
	}
	posts = append(posts, post)
	return c.Write(posts)
}

