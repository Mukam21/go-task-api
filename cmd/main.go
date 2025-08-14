package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mukam21/go-task-api/handler"
	"github.com/Mukam21/go-task-api/logger"
	"github.com/Mukam21/go-task-api/repository"
	"github.com/Mukam21/go-task-api/service"
)

func main() {
	// Инициализация логгера
	asyncLogger := logger.NewLogger(200)
	asyncLogger.Start()

	// Инициализация зависимостей
	repo := repository.NewTaskRepository()
	svc := service.NewTaskService(repo, asyncLogger)
	h := handler.NewTaskHandler(svc)

	// Роуты
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Запуск сервера
	go func() {
		log.Println("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	asyncLogger.Close()
	log.Println("Server stopped gracefully")
}
