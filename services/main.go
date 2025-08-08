package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"services/internal/processor"
	"services/internal/stream"
	"services/internal/web"
	"services/pkg/database"
	"syscall"
)

func main() {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	database.Connect(false)

	service := "ocserv"

	stringStream := make(chan string, 100)
	broadcaster := make(chan string, 1)

	ctx, cancel := context.WithCancel(context.Background())

	go processor.Processor(ctx, stringStream, broadcaster)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("\nReceived signal: %s\n", sig)
		cancel()
	}()

	if os.Getenv("DOCKERIZED") == "true" {
		log.Println("Running in Docker – using Docker logs")

		go func() {
			if err := stream.DockerStreamLogs(ctx, service, stringStream); err != nil {
				log.Println(err)
			}
		}()
	} else {
		log.Println("Running on host – using systemd logs")
		if err := stream.SystemdStreamLogs(ctx, service, stringStream); err != nil {
			log.Printf("Systemd log error: %v\n", err)
		}
	}

	go func() {
		processor.CalculateUserStats(ctx, stringStream)
	}()

	sseServer := web.NewSSEServer()
	sseServer.StartBroadcast(broadcaster)

	go func() {
		http.HandleFunc("/logs", sseServer.SSEHandler())
		log.Println("Starting server on 0.0.0.0:8081")
		if err := http.ListenAndServe("0.0.0.0:8081", nil); err != nil {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()
	
	<-ctx.Done() // wait for context cancellation
	log.Println("Main exited")

}
