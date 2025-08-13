package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	LabstackLog "github.com/labstack/gommon/log"
	"log"
	"net/http"
	"ocserv-service/group"
	"ocserv-service/occtl"
	"ocserv-service/user"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	ocpasswdPath = "/etc/ocserv/ocpasswd"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()

	server := "0.0.0.0:8081"

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	registerRoutes(e)

	verboseLog(e, server)

	go func() {
		e.Logger.SetLevel(LabstackLog.WARN)
		e.HideBanner = true

		if err := e.Start(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server due to error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("server forced to shutdown:", err)
	}

	log.Println("server shutdown complete")
}

func registerRoutes(e *echo.Echo) {
	g := e.Group("/api")
	occtl.Routes(g)
	user.Routes(g)
	group.Routes(g)
}

func verboseLog(e *echo.Echo, service string) {
	paths := e.Routes()
	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i].Path < paths[j].Path
	})
	maxNameLen := len("ROUTE NAME")
	maxPathLen := len("PATH")
	maxMethodLen := len("METHOD")
	for _, path := range paths {
		if len(path.Name) > maxNameLen {
			maxNameLen = len(path.Name)
		}
		if len(path.Path) > maxPathLen {
			maxPathLen = len(path.Path)
		}
		if len(path.Method) > maxMethodLen {
			maxMethodLen = len(path.Method)
		}
	}

	headerFormat := fmt.Sprintf("\n%%-%ds %%-%ds %%-%ds\n", maxNameLen+5, maxMethodLen, maxPathLen)
	log.Printf(headerFormat, "ROUTE NAME", "METHOD", "PATH")
	log.Println(strings.Repeat("-", maxNameLen+maxPathLen+maxMethodLen+3))

	rowFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds\n", maxNameLen+5, maxMethodLen, maxPathLen)
	for _, path := range paths {
		if !strings.HasSuffix(path.Name, ".init.func1") {
			log.Printf(
				rowFormat,
				strings.TrimSuffix(path.Name, "-fm"),
				path.Method,
				fmt.Sprintf("http://%s%s/", service, path.Path),
			)
		}
	}
}
