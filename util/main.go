package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var (
	htmlOutput = flag.Bool("html", false, "Output results in HTML format")
)

const (
	supportedExtensions = ".zip,.js,.jsx,.ts,.tsx,.py,.java,.cpp,.c,.cs,.go,.rb,.php"
	apiKey             = "API_KEY_PLACEHOLDER" // This will be replaced during build
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: scanner [options] <file>")
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Printf("\nSupported file types: %s\n", supportedExtensions)
		os.Exit(1)
	}

	filePath := flag.Args()[0]
	if err := validateFile(filePath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	result, err := analyzeFile(filePath)
	if err != nil {
		fmt.Printf("Analysis failed: %v\n", err)
		os.Exit(1)
	}

	if *htmlOutput {
		fmt.Println(result)
	} else {
		// Convert HTML to plain text table format
		plainText := convertHTMLToPlainText(result)
		fmt.Println(plainText)
	}
}

func validateFile(filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	supportedExts := strings.Split(supportedExtensions, ",")
	
	for _, supported := range supportedExts {
		if ext == supported {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported file type: %s. Supported types: %s", ext, supportedExtensions)
}

func analyzeFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}

	// Initialize OpenAI client
	client := openai.NewClient(apiKey)

	// Prepare the messages for the chat completion
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a security code analysis expert. Analyze the provided code and identify security vulnerabilities. Format the output as an HTML table with severity, issue, location, and description columns. Do not mention OpenAI or any AI in your response.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(`Analyze the following code for security vulnerabilities. 
Provide a detailed analysis in the following format:

<table>
  <tr>
    <th>Severity</th>
    <th>Issue</th>
    <th>Location</th>
    <th>Description</th>
  </tr>
  <!-- Add rows for each vulnerability found -->
</table>

Code to analyze:
%s`, string(content)),
		},
	}

	// Create the chat completion request
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4,
		Messages: messages,
	}

	// Send the request
	resp, err := client.CreateChatCompletion(req)
	if err != nil {
		return "", fmt.Errorf("error during analysis: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no analysis results received")
	}

	return resp.Choices[0].Message.Content, nil
}

func convertHTMLToPlainText(html string) string {
	// Basic HTML to plain text conversion
	text := html
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "<h1>", "\n# ")
	text = strings.ReplaceAll(text, "<h2>", "\n## ")
	text = strings.ReplaceAll(text, "<h3>", "\n### ")
	text = strings.ReplaceAll(text, "<table>", "\n")
	text = strings.ReplaceAll(text, "</table>", "\n")
	text = strings.ReplaceAll(text, "<tr>", "")
	text = strings.ReplaceAll(text, "</tr>", "\n")
	text = strings.ReplaceAll(text, "<td>", " | ")
	text = strings.ReplaceAll(text, "</td>", "")
	text = strings.ReplaceAll(text, "<th>", " | ")
	text = strings.ReplaceAll(text, "</th>", "")
	
	// Remove HTML tags
	text = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(text, "")
	
	// Clean up extra whitespace
	text = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)
	
	return text
} 