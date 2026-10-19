package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arlengur/go_assessment/internal/handler"
	"github.com/arlengur/go_assessment/internal/migrator"
	"github.com/arlengur/go_assessment/internal/pool"
	"github.com/arlengur/go_assessment/internal/repository"
	"github.com/arlengur/go_assessment/internal/service"
)

func main() {
	migratorRunner := migrator.NewMigrator()
	err := migratorRunner.Up()
	if err != nil {
		log.Printf("Database migration error: %v\n", err)
		return
	}

	repo := repository.NewRepository()
	pool := pool.NewWorkerPool(repo, 5, 100)
	pool.Start()
	defer pool.Stop()

	service := service.NewTaskService(pool, repo)
	taskHandler := handler.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, taskHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Starting API server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server dynamic crash: %v", err)
		}
	}()

	// Listening loop for OS graceful interrupts
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Initiating graceful teardown components...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
	log.Println("Server processing stopped. Exiting.")
}
