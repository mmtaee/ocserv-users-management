package web

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	clients map[chan string]string // IP
	mu      sync.Mutex
}

func NewSSEServer() *Server {
	return &Server{
		clients: make(map[chan string]string),
	}
}

func (s *Server) AddClient(client chan string, ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client] = ip
	log.Printf("Client %v (%s) connected", client, ip)
}

// RemoveClient removes a client connection
func (s *Server) RemoveClient(client chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ip, ok := s.clients[client]; ok {
		log.Printf("Client %v (%s) disconnected", client, ip)
		delete(s.clients, client)
		close(client)
	}
}

func (s *Server) StartBroadcast(broadcaster <-chan string) {
	go func() {
		for msg := range broadcaster {
			s.mu.Lock()
			for ch := range s.clients {
				select {
				case ch <- msg:
				default:
					continue
				}
			}
			s.mu.Unlock()
		}
	}()
}

func (s *Server) SSEHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		if len(s.clients) > 3 {
			s.mu.Unlock()
			http.Error(w, "Too many clients connected", http.StatusTooManyRequests)
			return
		}
		s.mu.Unlock()

		// Setup SSE headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		clientChan := make(chan string, 10)
		ip := r.RemoteAddr
		s.AddClient(clientChan, ip)
		defer s.RemoveClient(clientChan)

		ctx := r.Context()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-clientChan:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				flusher.Flush()
			}
		}
	}
}
