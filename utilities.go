package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
)

// Using `http.ResponseWriter`, write a string to
// bytes on ze chat interface.
func WriteBytes(w http.ResponseWriter, input string) {
	_, err := w.Write([]byte(input))
	if err != nil {
		slog.Error("error writing response body", "err", err)
	}
}

// Read a markdown file by inputting a file path.
// (Needs a file extension check)
func ReadMDFile(filePath string) string {
	resultBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	resultString := string(resultBytes)
	return resultString
}

// Callx crawl.py to conduct web search. Results are saved to `server/crawl_data/crawl_results.json`.
func CallCrawlScript() {
	cmd := exec.Command("python3", "crawl.py")

	// Output to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
