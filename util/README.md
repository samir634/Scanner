# Code Analysis CLI Utility

This command-line utility provides advanced code analysis to scan files for security vulnerabilities.

## Features

- Advanced code analysis engine
- Supports multiple file types (ZIP, JavaScript, TypeScript, Python, Java, C++, C, C#, Go, Ruby, PHP)
- Output in either HTML or plain text table format
- Detailed security vulnerability analysis

## Building the Binary

The utility can be built into a standalone binary that can be deployed to any Linux server. The API key is embedded in the binary during the build process.

### Using Bash (Linux/Mac)

```bash
# Make the build script executable
chmod +x build.sh

# Build the binary with your API key
./build.sh your-api-key
```

### Using PowerShell (Windows)

```powershell
# Build the binary with your API key
.\build.ps1 -ApiKey your-api-key
```

After building, you'll get a `scanner` binary that can be copied to any Linux server and run directly.

## Usage

Before running the scanner, set your API key as an environment variable:

```bash
# Windows
set OPENAI_API_KEY=your-api-key

# Linux/Mac
export OPENAI_API_KEY=your-api-key
```

Basic usage:
```bash
./scanner <file>
```

Options:
- `-html`: Output results in HTML format (default: plain text table)

Examples:
```bash
# Scan a JavaScript file with HTML output
./scanner -html script.js

# Scan a Python file with default output format
./scanner script.py
```

## Output Formats

### HTML Output
When using the `-html` flag, the output will be in HTML format with a table structure:
```html
<table>
  <tr>
    <th>Severity</th>
    <th>Issue</th>
    <th>Location</th>
    <th>Description</th>
  </tr>
  <!-- Vulnerability rows -->
</table>
```

### Plain Text Output (Default)
The default output format is a plain text table with the following structure:
```
| Severity | Issue | Location | Description |
|----------|-------|----------|-------------|
| High     | ...   | ...      | ...         |
| Medium   | ...   | ...      | ...         |
| Low      | ...   | ...      | ...         |
```

## Error Handling

The utility will display appropriate error messages for:
- Missing API key
- Unsupported file types
- File access issues
- Analysis service errors

## Requirements

- For building: Go 1.21 or later
- For running: Linux server (x86_64)
- Supported file types as listed above 