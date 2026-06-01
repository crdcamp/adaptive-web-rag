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

func CreateCollection(client *weaviate.Client, className string, description string) { // You're prob gonna need more parameters for this later
	ctx := context.Background()
	fmt.Println("Checking existence for collection: ", className)
	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)

	// There's probably a more elegant way to do this error handling
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
		Vectorizer:      "text2vec-openai", // Double check this. Might be the wrong vectorizer
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
