package main

import (
	"bnr/posts"
	"bnr/server"
	"bnr/users"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type config struct {
	Port string
}

var env config

func main() {

	server := server.NewServer()
	server.AddRoutes(users.GetRoutes())
	server.AddRoutes(posts.GetRoutes())

	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"})
	header := handlers.AllowedHeaders([]string{"Content-type", "*"})
	origins := handlers.AllowedOrigins([]string{})

	env.Port = os.Getenv("PORT")
	if env.Port == "" {
		e, err := os.Open("config.json")
		if err != nil {
			fmt.Printf("no config.json file found: %s\ndefaulting to OS ENV 'PORT'\n", err.Error())
		} else {
			json.NewDecoder(e).Decode(&env)
			e.Close()
		}
	}

	if env.Port != "" {
		fmt.Println(env.Port)
		log.Fatal(http.ListenAndServe(":"+env.Port, handlers.CORS(methods, header, origins)(server.Router)))
	} else {
		panic("PORT not set")
	}

}
