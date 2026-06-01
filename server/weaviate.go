package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

// YA REALLY GOTTA FIGURE OUT WHAT THE HELL THE `context` PACKAGE IS FOR!
// IT'S NOT GONNA STOP APPEARING ANY TIME SOON!

// YOU ALSO REALLY GOTTA FIGURE OUT PROPER ERROR HANDLING BEFORE YOU
// HAVE MORE TO DO THAN YOU SHOULD LATER ON!

const VectorizationMethod string = "text2vec-openai" // Not sure if this variable is necessary... given that you're probably not gonna need a different vectorization method you can probably get rid of it

func CreateCollection(client *weaviate.Client, className string, description string) { // You're probably gonna need more parameters for this later
	ctx := context.Background()
	fmt.Println("Checking existence for collection: ", className)
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)

	// There's probably a more elegant way to do these if statements (switch maybe?)
	if err != nil {
		panic(err)
	}
	if exists == true {
		fmt.Println("Collection already exists: ", className)
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

//func DeleteCollection() {}

// func EmbedText() {}

// func UploadCollection() {}
