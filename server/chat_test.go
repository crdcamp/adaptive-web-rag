package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRoot(t *testing.T) {
	w := httptest.NewRecorder()

	handleRoot(w, nil)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v, but got: %v\nbody:%s\n",
			desiredCode, w.Code, w.Body.String())
	}

	expectedMessage := []byte("Welcome to the root page. Hit the chat endpoint instead please.\n")
	if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
		t.Errorf("bad return, got: %q, expected: %q", w.Body.String(), expectedMessage)
	}
}

func TestHandleChat(t *testing.T) {
	testRequest := ChatPost{
		UserPrompt: "Tell me about how inconceivable the size of the universe is.",
	}

	marshalledRequestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("error marshalling test data: %v", err)
	}

	LlamaClient = CreateLlamaClient(AppConfig.LlamaBaseURL, AppConfig.LlamaAPIKey)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(marshalledRequestBody))

	handleChat(w, req)

	desiredCode := http.StatusOK
	if w.Code != desiredCode {
		t.Errorf("bad response code, expected: %v, but got: %v\nbody: %s\n",
			desiredCode, w.Code, w.Body.String())
	}
}
