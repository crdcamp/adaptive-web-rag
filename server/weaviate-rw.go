package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/crdcamp/charsplitter"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

type HrefAndContent struct {
	Href    string `json:"href"`
	Content string `json:"content"`
}

func CreateCollectionRw(client *weaviate.Client, className string, description string) {
	ctx := context.Background()
	// Does ClassCreator() overwrite new classes with the same name?
	fmt.Printf("Checking for existence of collection %q\n", className)
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Printf("Collection %q already exists\n", className)
		return
	}

	fmt.Printf("Creating collection %q\n", className)
	emptyClass := &models.Class{
		Class:           className,
		Description:     description,
		Vectorizer:      "text2vec-openai",
		VectorIndexType: "hnsw",
		Properties: []*models.Property{
			{
				Name:     "source",
				DataType: schema.DataTypeText.PropString(),
			},
			{
				Name:     "content",
				DataType: schema.DataTypeText.PropString(),
			},
		},
		ModuleConfig: map[string]interface{}{
			"text2vec-openai": map[string]interface{}{
				"baseURL":            LlamaBaseURL,
				"model":              EmbedModel,
				"vectorizeClassName": true,
			},
		},
	}
	err = client.Schema().ClassCreator().WithClass(emptyClass).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created class %q\n", className)
}

func GetCollectionRw(client *weaviate.Client, className string) []byte {
	ctx := context.Background()

	fmt.Println("Retrieving collection:", className)
	class, err := client.Schema().ClassGetter().
		WithClassName(className).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Collection retrieved:", className)

	b, err := json.MarshalIndent(class, "", " ")

	return b
}

// Delete a collection from your vector database.
func DeleteCollectionRw(client *weaviate.Client, className string) {
	fmt.Println("Deleting collection:", className)
	if err := client.Schema().ClassDeleter().WithClassName(className).Do(context.Background()); err != nil {
		// Weaviate will return a 400 if the class does not exist, so this is allowed, only return an error if it's not a 400
		if status, ok := err.(*fault.WeaviateClientError); ok && status.StatusCode != http.StatusBadRequest {
			panic(err)
		}
	}
	fmt.Println("Collection deleted:", className)
}

func SplitCrawlResults(fileName string) []HrefAndContent {
	fmt.Printf("Reading file %v and splitting text content\n", fileName)
	contentBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var embedJSON []HrefAndContent
	json.Unmarshal([]byte(contentBytes), &embedJSON)

	splitter := charsplitter.New(
		charsplitter.WithChunkSize(2000), // Character length, not word
		charsplitter.WithChunkOverlap(400),
		charsplitter.WithKeepSeparator(false),
	)

	var results []HrefAndContent

	for i := range embedJSON {
		content := embedJSON[i].Content
		chunks, err := splitter.SplitText(content)
		if err != nil {
			panic(err)
		}
		for _, chunk := range chunks {
			results = append(results, HrefAndContent{
				Href:    embedJSON[i].Href,
				Content: chunk,
			})

		}
	}

	return results
}
