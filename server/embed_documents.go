package main

func EmbedDocuments() {
	embedServer := StartLLMServer(EmbedModelPath, EmbedModelPort, true)
	WaitForServer(EmbedModelPort)
	StopLLMServer(embedServer)
}
