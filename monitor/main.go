package main

import (
	"bufio"
	"log"
	"os/exec"
	"os"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type message struct {
	WSToken string `json:"token"`
	Text    string `json:"text"`
}


func main() {
	logfile := os.Getenv("LOG_FILE")
	wsServer := os.Getenv("WS_SERVER")
	wsToken := os.Getenv("WS_TOKEN")
	if logfile == "" {
		log.Fatal("LOG_FILE environment variable not set")
	}
	if wsServer == "" {
		log.Fatal("WS_SERVER environment variable not set")
	}
	if wsToken == "" {
		log.Fatal("WS_TOKEN environment variable not set")
	}
	cmd := exec.Command("tail", "-f", "-n2",logfile)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Failed to create stdout pipe:", err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal("Failed to start command:", err)
	}
	lastLogEntry := "start script"
	conn, _, err := websocket.DefaultDialer.Dial(wsServer, nil)
	if err != nil {
		log.Fatal("Failed to connect to websocket:", err)
	}
	defer conn.Close()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if lastLogEntry == line {
			continue
		}
		lastLogEntry = line
		message := message {
			WSToken: wsToken,
			Text:  line,
		}
		jsonData, err := json.Marshal(message)
		if err != nil {
			log.Fatal("Failed to connect to websocket:", err)
		}
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Failed to send message via websocket:", err)
			break
		}
	}
	err = cmd.Wait()
	if err != nil {
		log.Println("Command execution ended with error:", err)
	}
}
