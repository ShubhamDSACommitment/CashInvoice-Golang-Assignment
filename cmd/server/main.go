package main

import (
	"github.com/CashInvoice-Golang-Assignment/internal/handler"
	"github.com/CashInvoice-Golang-Assignment/internal/repository"
	"github.com/CashInvoice-Golang-Assignment/internal/service"
	"github.com/CashInvoice-Golang-Assignment/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	//
	//db := database.Connect(cfg.DBUrl)
	//
	//queue := make(chan string, 100)
	//
	//taskRepo := repository.NewPostgresTaskRepo(db)
	//taskService := service.NewTaskService(taskRepo, queue)
	//
	//worker.StartAutoCompleteWorker(
	//	taskRepo,
	//	queue,
	//	time.Duration(cfg.AutoCompleteMinutes)*time.Minute,
	//)
	db := database.Connect()
	database.RunMigrations(db)

	// ---------------- Channel (Worker Queue) ----------------
	taskQueue := make(chan string, 100)

	// ---------------- Dependency Injection ----------------
	taskRepo := repository.NewMySQLTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, taskQueue)
	taskHandler := handler.NewTaskHandler(taskService)
	r := gin.Default()

	//auth := r.Group("/auth")
	//auth.POST("/login", handler.Login)
	//
	tasks := r.Group("/tasks")
	//tasks.Use(middleware.JWTMiddleware(cfg.JWTSecret))

	tasks.POST("", taskHandler.Create)
	tasks.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	//tasks.GET(":id", taskHandler.GetByID)
	//tasks.DELETE(":id", taskHandler.Delete)

	r.Run(":8080")
}
