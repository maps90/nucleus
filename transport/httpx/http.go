package httpx

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/labstack/echo"
)

var c *Config

func init() {
	c = New()
}

// Config struct
type Config struct {
	*echo.Echo
	port string
}

// ConfFunc implements http functions
type ConfFunc func(*Config) error

// New HTTP server
func New() *Config {
	if c != nil {
		return c
	}

	return &Config{
		Echo: echo.New(),
	}
}

// Setup gets the global httpx instance.
func Setup() *Config {
	return c
}

// SetPort config
func (h *Config) SetPort(port string) {
	h.port = port
}

// Serve the HTTP server
func (h *Config) Serve() error {
	// starting HTTP service
	go func() {
		if err := h.Start(h.port); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// shutdown process capture
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	log.Println("Shutting down server... Reason:", sig)
	if err := h.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

// Stop the HTTP server
func (h *Config) Stop() error {
	return h.Shutdown(context.Background())
}

// Set HTTP implementation
func (h *Config) Set(conf ...ConfFunc) error {
	for _, v := range conf {
		if err := v(h); err != nil {
			return err
		}
	}
	return nil
}

// HealthCheck will enable the health_check URL
func (h *Config) HealthCheck(enabled bool) {
	if enabled {
		h.Echo.GET("/health_check", func(ctx echo.Context) error {
			return ctx.JSON(http.StatusOK, map[string]string{
				"request_id": ctx.Response().Header().Get(echo.HeaderXRequestID),
				"status":     "HTTP Status OK!",
			})
		})
	}
}
