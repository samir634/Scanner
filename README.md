# Code Analysis Web Application

This application allows users to upload code files (zipped or individual files) for analysis using ChatGPT's API. The results are displayed in a formatted table on a unique page for each submission.

## Project Structure

- `/frontend` - NextJS frontend application
- `/backend` - Go backend server

## Setup Instructions

### Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Create a `.env.local` file with your environment variables:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```
4. Run the development server:
   ```bash
   npm run dev
   ```

### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install Go dependencies:
   ```bash
   go mod tidy
   ```
3. Create a `.env` file with your environment variables:
   ```
   OPENAI_API_KEY=your_api_key_here
   ```
4. Run the server:
   ```bash
   go run main.go
   ```

## Features
- File upload support for zip files and individual code files
- ChatGPT API integration for code analysis
- Unique page generation for each submission
- Formatted table display of analysis results 