package broker

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Broker handles connected clients and sends updates
// to each of them when receives a new connection or disconnection
type Broker struct {
	clients       map[chan string]bool
	newClients    chan chan string
	closedClients chan chan string
	updates       chan string
}

func NewBroker() *Broker {
	return &Broker{
		clients:       make(map[chan string]bool),
		newClients:    make(chan (chan string)),
		closedClients: make(chan (chan string)),
		updates:       make(chan string),
	}
}

func (b *Broker) Start() {
	go func() {
		for {
			select {
			case s := <-b.newClients:
				b.clients[s] = true
				go b.sendClientsCount()
				log.Println("New client connected")

			case s := <-b.closedClients:
				delete(b.clients, s)
				close(s)
				go b.sendClientsCount()
				log.Println("Client disconnected")

			case upd := <-b.updates:
				for s := range b.clients {
					s <- upd
				}
				log.Printf("Send update to %d clients", len(b.clients))
			}
		}
	}()
}

// ServeHTTP Broker's method handles and HTTP request at the "/events/" URL.
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	updateChan := make(chan string)

	b.newClients <- updateChan

	ctx := r.Context()
	go func() {
		<-ctx.Done()
		b.closedClients <- updateChan
	}()

	// Set the headers related to event streaming (SSE).
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {
		upd, open := <-updateChan
		if !open {
			break
		}

		_, err := fmt.Fprintf(w, "data: %s\n\n", upd)
		if err != nil {
			log.Printf("Error occurred during writing to ResponseWriter, %v\n", err)
		}

		f.Flush()
	}
}

func (b *Broker) sendClientsCount() {
	b.updates <- strconv.Itoa(len(b.clients))
}