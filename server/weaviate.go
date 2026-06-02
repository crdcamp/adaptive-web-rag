package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/crdcamp/charsplitter"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
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
func CreateCollection(client *weaviate.Client, className string, description string) { // You're probably gonna need more parameters for this later
	ctx := context.Background()
	fmt.Println("Checking existence for collection:", className)
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)

	// There's probably a more elegant way to do these if statements (switch maybe?)
	// Also gotta redo this to make the print statements a bit better
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Collection found:", className)
	}
	if exists == true {
		fmt.Println("Collection already exists:", className)
		return
	}

	// Class is missing a lot of parameters that show in the retrieval output
	fmt.Println("Creating class:", className)
	emptyClass := &models.Class{
		Class:           className,
		Description:     description,
		Vectorizer:      "text2vec-openai",
		VectorIndexType: "hnsw",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: schema.DataTypeText.PropString(),
			},
			{
				Name:     "body",
				DataType: schema.DataTypeText.PropString(),
			},
		},
	}

	err = client.Schema().ClassCreator().WithClass(emptyClass).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Class created:", className)
}

// MAKE SURE TO READ THE OUTPUT OF THIS
// THERE'S A FEW CONFIGURATIONS YOU NEED TO ADDRESS
// You also need to handle what happens if the collection doesn't exist
func GetCollection(client *weaviate.Client, className string) []byte {
	ctx := context.Background()

	fmt.Println("Retrieving collection: ", className)
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
	fmt.Println("Executing crawl.py")
	cmd := exec.Command("python3", "crawl.py")
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error when running command: ", err)
	} else {
		fmt.Println("crawl.py successfully executed")
	}
}

// Read `crawl_results.json` and upload results to your Weaviate vector database.
func ChunkEmbedAndUploadCrawlResults() {
	// Read `crawl_results.json`
	content, err := ioutil.ReadFile("crawl_data/crawl_results.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// Chunk results
	splitter := charsplitter.New(
		charsplitter.WithChunkSize(1024),
		charsplitter.WithChunkOverlap(150),
		charsplitter.WithKeepSeparator(False),
	)
	for url, pageContent := range payload {
		chunks, err := splitter.SplitText(pageContent)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Embed results

	// Upload to vector db with href as metadata
}
