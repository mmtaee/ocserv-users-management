package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"services/internal/processor"
	"services/internal/stream"
	"services/internal/web"
	"services/pkg/config"
	"services/pkg/database"
	"syscall"
)

func main() {
	service := "ocserv"
	stringStream := make(chan string, 100)
	broadcaster := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	debug := false

	if debugEnv := os.Getenv("DEBUG") == "true"; debugEnv {
		debug = true
	}

	config.Init(debug)
	database.Connect(debug)

	processor.Init()
	go processor.Processor(ctx, stringStream, broadcaster)
	go processor.CalculateUserStats(ctx, stringStream)
	go processor.UserExpiryCron(ctx)

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

	sseServer := web.NewSSEServer()
	sseServer.StartBroadcast(broadcaster)

	go func() {
		cfg := config.Get()
		server := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		http.HandleFunc("/logs", sseServer.SSEHandler())
		log.Println("Starting server on ", server)
		if err := http.ListenAndServe(server, nil); err != nil {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	<-ctx.Done() // wait for context cancellation
	log.Println("Main exited")

}
