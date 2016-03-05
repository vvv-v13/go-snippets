package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	SecretKey = "SuperSecretKey"
	Username  = "user"
	Password  = "pass"
)

type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	server := &http.Server{
		Addr:           ":1234",
		Handler:        nil,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/api/auth", handlerAuth)
	http.HandleFunc("/api/secret", handlerSecret)
	log.Println("Server listen on 1234")
	panic(server.ListenAndServe())
}

func handlerAuth(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Println("Request:", ip, r.URL.Path)

	log.Println(r.Header.Get("Content-Type"))

	var as AuthStruct

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&as)

	if (err != nil) || (as.Username != Username) || (as.Password != Password) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"description": "Invalid credentials", "status_code": 401, "error": "Bad Request"}`)
		return
	}

	userId := uuid.NewUUID()

	// Create JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["userid"] = userId

	// Expire in 120 minutes
	token.Claims["exp"] = time.Now().Add(time.Minute * 120).Unix()
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
	}

	data := map[string]string{
		"token": tokenString,
	}

	json.NewEncoder(w).Encode(data)
        //result, err := json.Marshal(data)
        //w.Write(result)
	
}

func handlerSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	payload, err := jwt_auth(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"description": "Invalid credentials", "status_code": 401, "error": "Bad Request"}`)
		return
	}

	data := map[string]string{
		"id": payload,
	}

	//result, err := json.Marshal(data)
	//w.Write(result)
	json.NewEncoder(w).Encode(data)
}

func jwt_auth(r *http.Request) (string, error) {

	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Println("Request:", ip, r.URL.Path)

	token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) { return []byte(SecretKey), nil })

	var payload string

	if err == nil && token.Valid {
		payload = token.Claims["userid"].(string)
	}

	return payload, err

}
