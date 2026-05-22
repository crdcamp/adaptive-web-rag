# Project Structure

## Go Files

* main.go         - DB setup, document ingestion (replaces main.py)
* query.go        - Collection querying (replaces query.py)  
* chat.go         - RAG chat loop (replaces chat_demo.py)
* embeddings.go   - HTTP client wrapper for /v1/embeddings

# Endpoints

```go
// Embedding
POST localhost:8081/v1/embeddings
{"model": "...", "input": "your text"}

// Chat
POST localhost:8082/v1/chat/completions
{"model": "...", "messages": [...], "stream": true}
```

# Main Go Script Claude Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Messages []ChatMessage `json:"messages"`
    Stream   bool          `json:"stream"`
}

type ChatResponse struct {
    Choices []struct {
        Message ChatMessage `json:"message"`
    } `json:"choices"`
}

func main() {
    payload := ChatRequest{
        Messages: []ChatMessage{
            {Role: "user", Content: "What is 2 + 2?"},
        },
        Stream: false,
    }

    body, err := json.Marshal(payload)
    if err != nil {
        panic(err)
    }

    resp, err := http.Post(
        "http://localhost:8080/v1/chat/completions",
        "application/json",
        bytes.NewReader(body),
    )
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var result ChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        panic(err)
    }

    if len(result.Choices) == 0 {
        panic("no choices returned")
    }

    fmt.Println(result.Choices[0].Message.Content)
}
```
