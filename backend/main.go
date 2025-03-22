package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Embed frontend files
//
//go:embed dist/*
var embeddedFiles embed.FS

func main() {
	InitDB() // Initialize SQLite database

	r := gin.Default() // Includes Logger and Recovery middleware

	// API routes
	api := r.Group("/api")
	{
		api.GET("/projects", GetProjects)
		api.POST("/projects", CreateProject)
		api.PATCH("/projects/:id", UpdateProject)
		api.DELETE("/projects/:id", DeleteProject)
		// Next Acttions
		api.GET("/next-actions", GetNextActions)
		api.POST("/next-actions", CreateNextAction)
		api.PATCH("next-actions/:id", UpdateNextAction)
		api.DELETE("next-actions/:id", DeleteNextAction)
	}

	// Serve embedded Vue app with proper MIME types
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasSuffix(path, "/") || path == "/" {
			path = "index.html"
		}
		// Remove leading slash and join with dist
		path = strings.TrimPrefix(path, "/")
		filePath := "dist/" + path

		content, err := embeddedFiles.ReadFile(filePath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		// Let Gin detect and set the correct Content-Type
		c.Data(http.StatusOK, "", content)
	})

	fmt.Println("Server running on http://localhost:8081")
	log.Fatal(r.Run(":8081"))
}
