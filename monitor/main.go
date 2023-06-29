package main

import (
	"bufio"
	"log"
	"os/exec"
	"fmt"
	"os"
	"github.com/gorilla/websocket"
)
var pl = fmt.Println
func main() {
	logfile := os.Getenv("LOG_FILE")
	ws_server := os.Getenv("WS_SERVER")
	if ws_server == "" {
		log.Fatal("WS_SERVER environment variable not set")
	}
	if logfile == "" {
		log.Fatal("LOG_FILE environment variable not set")
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
	conn, _, err := websocket.DefaultDialer.Dial(ws_server, nil)
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
		err = conn.WriteMessage(websocket.TextMessage, []byte(line))
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
