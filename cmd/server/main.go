package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/CashInvoice-Golang-Assignment/internal/config"
	"github.com/CashInvoice-Golang-Assignment/internal/handler"
	"github.com/CashInvoice-Golang-Assignment/internal/middleware"
	"github.com/CashInvoice-Golang-Assignment/internal/repository"
	"github.com/CashInvoice-Golang-Assignment/internal/service"
	"github.com/CashInvoice-Golang-Assignment/internal/worker"
	"github.com/CashInvoice-Golang-Assignment/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.Load()
	db := database.Connect(cfg)
	database.RunMigrations(db)

	taskQueue := make(chan string, 100)
	wg := &sync.WaitGroup{}

	taskRepo := repository.NewMySQLTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, taskQueue)
	taskHandler := handler.NewTaskHandler(taskService)
	userRepo := repository.NewMySQLUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(
		authService,
		os.Getenv("JWT_SECRET"),
		10,
	)
	// Public
	delay := time.Duration(cfg.AutoCompleteMinutes) * time.Minute
	worker := worker.NewAutoCompleteWorker(taskRepo, taskQueue, delay, wg)
	worker.Start(ctx, 4)

	r := gin.Default()

	// Public
	auth := r.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)

	// Protected
	tasks := r.Group("/tasks")
	tasks.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	tasks.POST("", taskHandler.Create)
	tasks.GET("", taskHandler.GetAllTask)
	tasks.GET("/:id", taskHandler.GetByID)
	tasks.DELETE("/:id", taskHandler.Delete)

	// Admin-only group
	admin := auth.Group("/admin")
	admin.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	admin.Use(middleware.AdminOnly()) // we’ll create this
	admin.POST("/register", authHandler.RegisterAdmin)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received")

	// Stop accepting new requests
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Server shutdown error:", err)
	}

	cancel()
	close(taskQueue)

	wg.Wait()

	// Close DB
	err := db.Close()
	if err != nil {
		log.Println("ERROR CLOSING DB ", err)
	}

	log.Println("Graceful shutdown complete")

}
