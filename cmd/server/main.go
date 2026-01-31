package main

import (
	"os"
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
	cfg := config.Load()
	db := database.Connect(cfg)
	database.RunMigrations(db)

	taskQueue := make(chan string, 100)

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
	worker := worker.NewAutoCompleteWorker(taskRepo, taskQueue, delay)
	worker.Start(4)

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

	r.Run(":8080")
}
