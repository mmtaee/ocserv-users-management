package routing

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	LabstackLog "github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"ocserv-bakend/internal/providers/routing"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/routing/middlewares"
	"sort"
	"strings"
	"time"
)

var e *echo.Echo

func Serve(debug bool) {
	cfg := config.Get()
	server := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	e = echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.TimeoutMiddleware(10 * time.Second))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodPut,
			http.MethodHead,
			http.MethodOptions,
		},
	}))

	routing.Register(e)

	if debug {
		e.Debug = true
		e.Logger.SetLevel(LabstackLog.DEBUG)
		verboseLog(server)
	} else {
		e.Logger.SetLevel(LabstackLog.WARN)
		e.HideBanner = true
	}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "Healthy",
		})
	})

	if debug {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
	}))

	if err := e.Start(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("shutting down the server", err)
	}
}

func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("server shutdown complete")
}

func verboseLog(service string) {
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
