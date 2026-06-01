package main

// Reference: https://docs.weaviate.io/weaviate/manage-collections/collection-operations

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func CreateCollection(client *weaviate.Client, className string, description string) {
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

	fmt.Println("Creating class:", className)
	emptyClass := &models.Class{
		Class:           className,
		Description:     description,
		Vectorizer:      "text2vec-openai", // Double check this. Probably wrong one
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
