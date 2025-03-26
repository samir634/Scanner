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
	prompt := fmt.Sprintf(`You are a security code analysis expert. Your task is to analyze code for security vulnerabilities.
You must follow these exact rules:

1. Always use these exact severity levels with these exact criteria:
   HIGH: Immediate security threat, direct exploitation possible, sensitive data exposure
   MEDIUM: Security weakness that requires specific conditions to exploit
   LOW: Minor security concern or best practice violation

2. For each issue found, you must provide:
   - Exact line numbers where the issue occurs
   - Specific function names, variables, or code patterns involved
   - Clear steps to reproduce the vulnerability
   - Concrete examples of how it could be exploited
   - Specific, actionable fix recommendations

3. You must categorize each issue into exactly one of these categories:
   AUTH: Authentication and Authorization issues
   INPUT: Input validation and sanitization
   CRYPTO: Cryptographic issues
   EXPOSURE: Data exposure and privacy
   INJECTION: Code injection vulnerabilities
   CONFIG: Configuration and deployment issues
   DEPS: Dependency and library issues
   ACCESS: Access control problems
   LOGGING: Error handling and logging issues
   SESSION: Session management
   FILES: File operation security
   NETWORK: Network security issues
   MEMORY: Memory management
   LOGIC: Business logic flaws

4. Format Requirements:
   - Use exact column names as specified
   - Keep descriptions concise but complete
   - Include specific code references
   - Always provide actionable recommendations

Present your findings in this exact HTML table format:

<table>
  <thead>
    <tr>
      <th>Category</th>
      <th>Severity</th>
      <th>Issue</th>
      <th>Affected Components</th>
      <th>Description</th>
      <th>Recommendation</th>
    </tr>
  </thead>
  <tbody>
    <!-- For each finding, create a row with:
    - Category: Use exact category codes (AUTH, INPUT, etc.)
    - Severity: Must be exactly as follows:
      <span class="severity-high">HIGH</span>
      <span class="severity-medium">MEDIUM</span>
      <span class="severity-low">LOW</span>
    - Issue: Brief, specific title
    - Affected Components: "Line X-Y in functionName()" or specific variables/functions
    - Description: Follow this exact format:
      Impact: [Security impact]
      Vulnerability: [How it can be exploited]
      Context: [Relevant code context]
    - Recommendation: Specific, actionable steps to fix
    -->
  </tbody>
</table>

Code to analyze:

%s`, string(content))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4-1106-preview",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1, // Add low temperature for more consistent output
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