package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bimross/skills-mcp-server/internal/app"
	"github.com/bimross/skills-mcp-server/internal/config"
)

func main() {
	cfg := config.FromEnv()
	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("build app: %v", err)
	}

	httpServer := &http.Server{
		Addr:              cfg.ListenAddress(),
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("skills-mcp-server listening on %s", cfg.ListenAddress())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
