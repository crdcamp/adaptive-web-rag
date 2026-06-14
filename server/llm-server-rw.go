package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func UnloadModelRw(modelName string) {
	const unloadURL = LlamaBaseUrl + "/models/unload"
	// Need to research more into json encoding in Go. I have no idea how this works at the moment
	payload, err := json.Marshal(map[string]string{"model": modelName})
	if err != nil {
		panic(err)
	}
	fmt.Println("Unloading model:", modelName)
	resp, err := http.Post(unloadURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err) // Need a better error handling method here
	}
	defer resp.Body.Close()
	//fmt.Printf("Status: %s\n", resp.Status)
	fmt.Println("Model unloaded:", modelName)
}
