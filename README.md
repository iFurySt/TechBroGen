# TechBro Generator

A satirical TechBro post generator powered by Google Gemini AI. Generate cringe-worthy, overly motivational tech bro social media posts in different categories.

## Features

- 🤖 AI-powered post generation using Google Gemini
- 📱 Three categories: SaaS, AI, and Growth Marketing
- 🎨 Beautiful UI inspired by the original TechBro Generator
- 🔄 Random profile names and handles
- 📤 Direct sharing to X (Twitter)

## Setup

1. Clone this repository
2. Install Go dependencies: `go mod tidy`
3. Run the server with environment variables:
   ```bash
   GEMINI_API_KEY=your_api_key_here go run main.go
   ```
   or specify port:
   ```bash
   PORT=8000 GEMINI_API_KEY=your_api_key_here go run main.go
   ```
4. Open http://localhost:8080 in your browser

## Environment Variables

- `GEMINI_API_KEY` - Your Google Gemini API key (required)
- `PORT` - Server port (default: 8080)

## API Endpoints

- `GET /` - Serve the frontend
- `POST /api/generate` - Generate a new TechBro post
- `GET /health` - Health check
