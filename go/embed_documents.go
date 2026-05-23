package main

func EmbedDocuments() {
	startServer := StartLLMServer(EmbedModelPath, EmbedModelPort, true)
	StopLLMServer(startServer)
}
