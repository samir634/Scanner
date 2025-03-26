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
	systemPrompt := `You are a security code analysis expert. Your task is to analyze code for security vulnerabilities.
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
   - Always provide actionable recommendations`

	analysisPrompt := fmt.Sprintf(`Analyze the following code and present the findings in this exact HTML table format:

<table>
  <thead>
    <tr>
      <th>Severity</th>
      <th>Category</th>
      <th>Location</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <!-- For each finding, create a row with:
    - Severity: Must be exactly as follows:
      <span class="severity-high">HIGH</span>
      <span class="severity-medium">MEDIUM</span>
      <span class="severity-low">LOW</span>
    - Category: Use exact category codes (AUTH, INPUT, etc.)
    - Location: "Line X-Y in functionName()" or "Line X in fileName"
    - Description: Follow this exact format:
      Issue: [Brief title]
      Component: [Affected code element]
      Impact: [Security impact]
      Fix: [Specific solution]
    -->
  </tbody>
</table>

Code to analyze:

%s`, string(content))

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: analysisPrompt,
		},
	}

	// Create the chat completion request
	req := openai.ChatCompletionRequest{
		Model:    "gpt-4-1106-preview",
		Messages: messages,
		Temperature: 0.1, // Add low temperature for more consistent output
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

	// Handle severity levels with ANSI color codes
	text = strings.ReplaceAll(text, `<span class="severity-high">HIGH</span>`, "\033[1;31mHIGH\033[0m")     // Bold Red
	text = strings.ReplaceAll(text, `<span class="severity-medium">MEDIUM</span>`, "\033[1;33mMEDIUM\033[0m") // Bold Yellow
	text = strings.ReplaceAll(text, `<span class="severity-low">LOW</span>`, "\033[1;32mLOW\033[0m")       // Bold Green

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
	
	// Remove remaining HTML tags
	text = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(text, "")
	
	// Clean up extra whitespace
	text = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)
	
	return text
} 