package main

func EmbedDocuments() {
	startServer := StartLLMServer(EmbedModelPath, EmbedModelPort, true)
	WaitForServer(EmbedModelPort)
	StopLLMServer(startServer)
}
