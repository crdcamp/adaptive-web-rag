package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/crdcamp/charsplitter"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

// YA REALLY GOTTA FIGURE OUT WHAT THE HELL THE `context` PACKAGE IS FOR!
// IT'S NOT GONNA STOP APPEARING ANY TIME SOON!

// YOU ALSO REALLY GOTTA FIGURE OUT PROPER ERROR HANDLING BEFORE YOU
// HAVE MORE TO DO THAN YOU SHOULD LATER ON!

// Not sure if this variable is necessary... given that you're probably not gonna need a different vectorization method you can probably get rid of it
// But.. what if you want to easily change the vectorization method?
// const VectorizationMethod string = "text2vec-openai"

// Create a Weaviate vector database collection. Note: `className` must be camelcase.
func CreateCollection(client *weaviate.Client, className string, description string) {
	ctx := context.Background()

	fmt.Printf("Checking if collection %q exists\n", className)
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		panic(err)
	}
	if exists {
		fmt.Printf("Collection %q already exists\n", className)
		return
	}
	fmt.Printf("Collection %q does not exist. Creating collection %q\n", className, className)
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
				"baseURL":            LlamaBaseUrl,
				"model":              EmbedModel,
				"vectorizeClassName": true,
			},
		},
	}
	err = client.Schema().ClassCreator().WithClass(emptyClass).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Class %q created", className)
}

// MAKE SURE TO READ THE OUTPUT OF THIS
// THERE'S A FEW CONFIGURATIONS YOU NEED TO ADDRESS
// You also need to handle what happens if the collection doesn't exist
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

// Delete a collection from your vector database
func DeleteCollection(client *weaviate.Client, className string) {
	fmt.Println("Deleting collection:", className)
	if err := client.Schema().ClassDeleter().WithClassName(className).Do(context.Background()); err != nil {
		// Weaviate will return a 400 if the class does not exist, so this is allowed, only return an error if it's not a 400
		if status, ok := err.(*fault.WeaviateClientError); ok && status.StatusCode != http.StatusBadRequest {
			panic(err)
		}
	}
	fmt.Println("Collection deleted:", className)
}

// Calls crawl.py to conduct web search. Results are saved to `server/crawl_data/crawl_results.json`.
func CallCrawlScript() {
	// Run crawl.py
	fmt.Println("Executing: crawl.py")
	cmd := exec.Command("python3", "crawl.py")
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error when running command: ", err)
	} else {
		fmt.Println("Successfully executed: crawl.py")
	}
}

// HORRENDOUS FUNCTION! BREAK THIS DOWN AND GET IT TOGETHER MAN!
// Read `crawl_results.json` and upload results to your Weaviate vector database.
// I'm probably creating way too many structs here.
// Need to come back to this when I have more understanding of Go
// Regardless, I think there's a lot of redundancy here. This function is a mess cause idk anything about handling JSON in Go
func SplitEmbedAndUploadCrawlResults(client *weaviate.Client, targetCollection string) {
	ctx := context.Background()
	// Read `crawl_results.json`
	content, err := os.ReadFile("crawl_data/crawl_results.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Create a struct for each key's values
	type HrefContent struct {
		Href    string `json:"href"`
		Content string `json:"content"`
	}

	// Unmarshal the data into HrefContent
	jsonMap := map[string]HrefContent{}
	err = json.Unmarshal(content, &jsonMap)
	if err != nil {
		panic(err)
	}

	// Text splitter definition
	splitter := charsplitter.New(
		charsplitter.WithChunkSize(512),
		charsplitter.WithChunkOverlap(100),
		charsplitter.WithKeepSeparator(false),
	)

	// Create chunk object
	type WeaviateChunkObject struct {
		Href  string `json:"href"`
		Chunk string `json:"chunk"`
	}

	var docSeq int = 1
	var totalDocs int = len(jsonMap)
	// Split text
	for _, hrefAndContent := range jsonMap {
		fmt.Printf("Splitting and embedding text for href(%v/%v): %v\n", docSeq, totalDocs, hrefAndContent.Href)
		chunks, err := splitter.SplitText(hrefAndContent.Content)
		if err != nil {
			log.Printf("Failed to split text for %s, %v", hrefAndContent.Href, hrefAndContent.Content)
		}
		// Loop through each individual chunk
		for _, chunkText := range chunks {
			chunkPayload := WeaviateChunkObject{
				Href:  hrefAndContent.Href,
				Chunk: chunkText,
			}

			// Convert chunks into a slice of models.Object
			objects := []models.PropertySchema{}
			objects = append(objects, map[string]interface{}{
				"source": chunkPayload.Href,
				"body":   chunkPayload.Chunk,
			})

			// Batch write items
			batcher := client.Batch().ObjectsBatcher()
			for _, dataObj := range objects {
				batcher.WithObjects(&models.Object{
					Class:      targetCollection,
					Properties: dataObj,
				})
			}

			batchRes, err := batcher.Do(ctx)

			if err != nil {
				panic(err)
			}
			for _, res := range batchRes {
				if res.Result.Errors != nil {
					for _, err := range res.Result.Errors.Error {
						fmt.Printf("Error details: %v\n", *err)
						panic(err.Message)
					}
				}
			}
		}
	}
	docSeq = docSeq + 1
	UnloadModel(EmbedModel)
}

func NearTextSearch(client *weaviate.Client, className string, limit int, query string) {
	ctx := context.Background()

	fmt.Println("NearTextSearch() called for query:", query, "\nQuerying...")
	nearTextResponse, err := client.GraphQL().Get().
		WithClassName(className).
		WithFields(
			graphql.Field{Name: "source"},
			graphql.Field{Name: "body"},
		).
		WithNearText(client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{query})).
		WithLimit(limit).
		Do(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Println("NearTextSearch() complete")
	UnloadModel(EmbedModel)
	fmt.Println("NearTextSearch() result:\n", nearTextResponse)
}
