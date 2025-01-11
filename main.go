package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	router := Routes() //Makes an instnace of routes

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, //Ensures we only serve connections from this origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(router)
	http.Handle("/", handler)
	log.Println("Listening on port :3000")
	http.ListenAndServe(":3000", nil)
}
