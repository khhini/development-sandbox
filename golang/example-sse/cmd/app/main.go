package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Broker struct {
	clients map[chan string]bool
	lock    sync.Mutex
}

func newBroker() *Broker {
	return &Broker{clients: make(map[chan string]bool)}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msgChan := make(chan string)
	b.lock.Lock()
	b.clients[msgChan] = true
	b.lock.Unlock()

	defer func() {
		b.lock.Lock()
		delete(b.clients, msgChan)
		b.lock.Unlock()
		close(msgChan)
	}()

	for msg := range msgChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func (b *Broker) Broadcast(msg string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for client := range b.clients {
		client <- msg
	}
}

func main() {
	broker := newBroker()
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		broker.Broadcast(msg)
		fmt.Fprintf(w, "Sent: %s", msg)
	})

	http.HandleFunc("/events-page", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/events.html")
	})

	http.HandleFunc("/events", broker.ServeHTTP)
	http.Handle("/", http.FileServer(http.Dir("public")))

	log.Println("📡 Broadcasting SSE server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func tickHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			fmt.Fprintf(w, "data: %s\n\n", t.Format(time.RFC1123))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			fmt.Println("Cleint closed connection")
			return
		}
	}
}
