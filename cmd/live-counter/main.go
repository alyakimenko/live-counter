package main

import (
	"log"
	"net/http"

	"github.com/alyakimenko/live-counter/internal/broker"
)

func main() {
	b := broker.NewBroker()
	b.Start()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.Handle("/events/", b)

	log.Println("Server is running on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}