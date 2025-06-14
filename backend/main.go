package main

import (
	"ocserv-bakend/cmd"
	_ "ocserv-bakend/docs"
)

// @title Ocserv User management Example Api
// @version 1.0
// @description This is a sample Ocserv User management Api server.
// @BasePath /api
func main() {
	cmd.Execute()
}
