package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var pl = fmt.Println

func getTokenFromFile(key string, passwdFile string) string {
	cmd := exec.Command("grep", "-r", key, passwdFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Monitor] getting token from %s faile with error: %s", passwdFile, err)
		return ""
	}
	return strings.Split(string(output), ":")[1]
}

func main() {
	log.SetOutput(os.Stdout)
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	wsPath := os.Getenv("SOCKET_PATH")
	if wsPath == "" {
		wsPath = "/"
	}
	logFile := "/shared_mointor/ocserv.log"
	passwdFile := "/shared_mointor/socket_passwd"
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			queryParams := r.URL.Query()
			user, qToken := queryParams["user"][0], queryParams["token"][0]
			token := getTokenFromFile(user, passwdFile)
			if strings.TrimSpace(token) == strings.TrimSpace(qToken) {
				log.Printf("[Monitor] Socket connection (%s) accepted", r.Host)
				return true
			}
			log.Printf("[Monitor] Socket connection (%s) rejected!", r.Host)
			return false
		},
	}

	http.HandleFunc(wsPath, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("[Monitor] Failed to upgrade connection to WebSocket:", err)
			return
		}
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println("[Monitor] Failed to create file watcher:", err)
			conn.Close()
			return
		}
		defer watcher.Close()
		absFilePath, err := filepath.Abs(logFile)
		if err != nil {
			log.Println("[Monitor] Failed to get absolute path for file:", err)
			conn.Close()
			return
		}
		file, err := os.Open(absFilePath)
		if err != nil {
			log.Println("[Monitor] Failed to open file:", err)
			conn.Close()
			return
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			log.Println("[Monitor] Failed to get file info:", err)
			conn.Close()
			return
		}
		initialReadPos := fileInfo.Size()
		err = watcher.Add(absFilePath)
		if err != nil {
			log.Println("[Monitor] Failed to watch file:", err)
			conn.Close()
			return
		}
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					if event.Op&fsnotify.Write == fsnotify.Write {
						file.Seek(initialReadPos, 0)
						reader := bufio.NewReader(file)
						for {
							line, err := reader.ReadString('\n')
							if err != nil {
								break
							}
							line = strings.TrimSpace(line)
							if line != "" {
								err = conn.WriteMessage(websocket.TextMessage, []byte(line))
								if err != nil {
									log.Println("[Monitor] Failed to send message to client:", err)
									continue
								}
							}
						}
						fileInfo, err := file.Stat()
						if err != nil {
							log.Println("[Monitor] Failed to get file info:", err)
							continue
						}
						initialReadPos = fileInfo.Size()
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("[Monitor] Watchererror:", err)
				}
			}
		}()
		_, _, err = conn.ReadMessage()
		if err != nil {
			log.Println("[Monitor] Failed to read message from client:", err)
		}
	})
	log.Printf("[Monitor] Server is running on http://%s:%s", host, port)
	log.Printf("[Monitor] Socket passwd file(%s)", passwdFile)
	log.Printf("[Monitor] Log file(%s)", logFile)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	if err != nil {
		log.Fatal("[Monitor] Openning socket error", err)
	}

}
