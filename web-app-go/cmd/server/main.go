package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "project/docs"

	"project/internal/handlers"
)

// @title Text Editor API
// @version 1.0
// @description A simple file-based text editor API in Go with Swagger
// @host localhost:5007
// @BasePath /api/files
func main() {
	r := gin.Default()

	// Serve static HTML
	r.Static("/public", "./public")
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// Initialize files directory
	filesDir := filepath.Join(".", "files")
	if err := os.MkdirAll(filesDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create files dir: %v", err)
	}
	handlers.SetFilesDir(filesDir)

	// API routes
	api := r.Group("/api/files")
	{
		api.GET("", handlers.ListFiles)
		api.GET("/:filename", handlers.GetFile)
		api.POST("/:filename", handlers.SaveFile)
		api.DELETE("/:filename", handlers.DeleteFile)
		api.POST("", handlers.CreateFile)
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("Server running at http://localhost:5007")
	if err := r.Run(":5007"); err != nil {
		log.Fatal(err)
	}
}
