package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type AnalysisResult struct {
	ID     string      `json:"id"`
	Data   interface{} `json:"data"`
	Status string      `json:"status"` // "processing", "completed", "error"
}

var (
	results = make(map[string]AnalysisResult)
	mutex   = &sync.RWMutex{}
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal(err)
	}

	r.POST("/upload", handleUpload)
	r.GET("/results/:id", handleGetResults)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Generate unique ID for this analysis
	id := uuid.New().String()

	// Save the file
	filename := filepath.Join("uploads", id+"_"+file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Initialize result with processing status
	mutex.Lock()
	results[id] = AnalysisResult{
		ID:     id,
		Status: "processing",
	}
	mutex.Unlock()

	// Start analysis in background
	go analyzeFile(id, filename)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func handleGetResults(c *gin.Context) {
	id := c.Param("id")
	
	mutex.RLock()
	result, exists := results[id]
	mutex.RUnlock()
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Results not found"})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

func analyzeFile(id string, filename string) {
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		mutex.Lock()
		results[id] = AnalysisResult{
			ID:     id,
			Data:   map[string]string{"error": "Failed to read file"},
			Status: "error",
		}
		mutex.Unlock()
		return
	}

	// Initialize OpenAI client
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Ask for HTML formatted table
	prompt := fmt.Sprintf("Analyze this code for any security vulnerabilities. Format your response as an HTML table with proper <table>, <tr>, <th>, and <td> tags. Make sure to use <thead> and <tbody> sections. The table should have the following columns: Issue, Description, Severity (High/Medium/Low), and Recommendation. Here's the code to analyze:\n\n%s", string(content))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		log.Printf("Error calling OpenAI API: %v", err)
		mutex.Lock()
		results[id] = AnalysisResult{
			ID:     id,
			Data:   map[string]string{"error": "Failed to analyze code"},
			Status: "error",
		}
		mutex.Unlock()
		return
	}

	// Store the results
	mutex.Lock()
	results[id] = AnalysisResult{
		ID:     id,
		Data:   resp.Choices[0].Message.Content,
		Status: "completed",
	}
	mutex.Unlock()

	// Clean up the uploaded file
	os.Remove(filename)
} 