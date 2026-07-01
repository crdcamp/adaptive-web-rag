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
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

type HrefAndContent struct {
	Href    string `json:"href"`
	Content string `json:"content"`
}

type CollectionNames struct {
	Classes []struct {
		Class string `json:"class"`
	}
}

func CreateCollection(client *weaviate.Client, className string, description string) {
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
				"baseURL":            AppConfig.LlamaBaseURL,
				"model":              AppConfig.EmbedModel,
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

func GetCollection(client *weaviate.Client, className string) []byte {
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
func DeleteCollection(client *weaviate.Client, className string) {
	ctx := context.Background()
	fmt.Println("Deleting collection:", className)
	if err := client.Schema().ClassDeleter().WithClassName(className).Do(ctx); err != nil {
		// Weaviate will return a 400 if the class does not exist, so this is allowed, only return an error if it's not a 400
		if status, ok := err.(*fault.WeaviateClientError); ok && status.StatusCode != http.StatusBadRequest {
			panic(err)
		}
	}
	fmt.Println("Collection deleted:", className)
}

func ReadAllCollectionDefinitions(client *weaviate.Client) []byte {
	ctx := context.Background()
	schema, err := client.Schema().Getter().
		Do(ctx)

	b, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}

func ReadAllCollectionNames(client *weaviate.Client) CollectionNames {
	// Needs error handling
	str := string(ReadAllCollectionDefinitions(client))
	res := CollectionNames{}
	_ = json.Unmarshal([]byte(str), &res)

	return res
}

func SplitCrawlResults(fileName string) []models.PropertySchema {
	fmt.Printf("Reading file %v and splitting text content\n", fileName)
	contentBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var embedJSON []HrefAndContent
	json.Unmarshal(contentBytes, &embedJSON)

	splitter := charsplitter.New(
		charsplitter.WithChunkSize(2000),
		charsplitter.WithChunkOverlap(400),
		charsplitter.WithKeepSeparator(false),
	)

	var results []models.PropertySchema

	for i := range embedJSON {
		content := embedJSON[i].Content
		chunks, err := splitter.SplitText(content)
		if err != nil {
			panic(err)
		}
		for _, chunk := range chunks {
			// Construct the Weaviate property schema map directly
			props := map[string]interface{}{
				"source":  embedJSON[i].Href,
				"content": chunk,
			}
			results = append(results, props)
		}
	}
	fmt.Printf("Done splitting text for file %v\n", fileName)

	return results
}

func EmbedText(client *weaviate.Client, className string, splitText []models.PropertySchema) {
	fmt.Println("Beginning text embeddings")
	ctx := context.Background()
	totalIterations := len(splitText)
	for i, text := range splitText {
		fmt.Printf("Embedding text (%v/%v)\n", i+1, totalIterations)
		_, err := client.Batch().ObjectsBatcher().WithObjects(&models.Object{
			Class:      className,
			Properties: text,
		}).Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Text embeddings complete")
}

// Needs a class existence check
func NearTextSearch(client *weaviate.Client, className string, limit int, query string) {
	ctx := context.Background()

	nearText := client.GraphQL().NearTextArgBuilder().
		WithConcepts([]string{query})

	response, err := client.GraphQL().Get().
		WithClassName(className).
		WithFields(
			graphql.Field{Name: "source"},
			graphql.Field{Name: "content"},
		).
		WithNearText(nearText).
		WithLimit(limit).
		Do(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Println(response)
}
