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
