package main

import (
        "context"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "test-task-lo/internal/handlers"
        "test-task-lo/internal/logger"
        "test-task-lo/internal/services"
        "test-task-lo/internal/storage"
        "time"
)

func main() {
        log := logger.NewLogger(256)
        defer log.Shutdown()

        log.Info("Starting application...", nil)

        repo := storage.NewTaskStorage(log)
        svc := services.NewTaskService(repo, log)
        taskHandlers := handlers.NewTaskHandlers(svc, log)

        mux := http.NewServeMux()
        mux.HandleFunc("GET /tasks", taskHandlers.GetAllTasks)
        mux.HandleFunc("GET /tasks/{id}", taskHandlers.GetTaskByID)
        mux.HandleFunc("POST /tasks", taskHandlers.CreateTask)

        server := &http.Server{
                Addr:    ":8080",
                Handler: mux,
        }

        serverErr := make(chan error, 1)

        go func() {
                log.Info("Server is running on http://localhost:8080", nil)
                if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                        serverErr <- err
                }
        }()

        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

        select {
        case err := <-serverErr:
                log.Error("Server error", map[string]interface{}{"error": err.Error()})
                os.Exit(1)
        case sig := <-sigint:
                log.Info("Received signal, shutting down gracefully...",
                        map[string]interface{}{"signal": sig.String()})

                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()

                if err := server.Shutdown(ctx); err != nil {
                        log.Error("Server shutdown error",
                                map[string]interface{}{"error": err.Error()})
                        server.Close()
                }
        }

        log.Info("Server stopped", nil)
}
