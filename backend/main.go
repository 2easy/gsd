package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Embed frontend files
//
//go:embed dist/*
var embeddedFiles embed.FS

func main() {
	dbPath := flag.String("db", "./gsd.db", "path to the SQLite database file")
	port := flag.String("port", "8080", "port to run the server on")
	flag.Parse()

	InitDB(*dbPath)    // Initialize SQLite database
	r := gin.Default() // Includes Logger and Recovery middleware

	// API routes
	api := r.Group("/api")
	{
		api.GET("/projects", GetProjects)
		api.POST("/projects", CreateProject)
		api.PATCH("/projects/:id", UpdateProject)
		api.DELETE("/projects/:id", DeleteProject)
		// Next Actions
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
			// Try index.html for SPA routing
			if content, err = embeddedFiles.ReadFile("dist/index.html"); err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			// For SPA routing, always serve as html
			c.Header("Content-Type", "text/html")
			c.Data(http.StatusOK, "text/html", content)
			return
		}

		// Set the correct content type based on file extension
		ext := filepath.Ext(path)
		var contentType string
		switch ext {
		case ".js":
			contentType = "text/javascript"
		case ".css":
			contentType = "text/css"
		case ".html":
			contentType = "text/html"
		case ".svg":
			contentType = "image/svg+xml"
		case ".json":
			contentType = "application/json"
		default:
			contentType = mime.TypeByExtension(ext)
			if contentType == "" {
				contentType = http.DetectContentType(content)
			}
		}

		c.Header("Content-Type", contentType)
		c.Data(http.StatusOK, contentType, content)
	})

	fmt.Printf("Server running on http://localhost:%s\n", *port)
	log.Fatal(r.Run(":" + *port))
}
