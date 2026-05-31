package main

// Reference: https://docs.weaviate.io/weaviate/manage-collections/collection-operations

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func GetSchema() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	result, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nSCHEMA: %v", result)
}

func Test() {
	GetSchema()
}

func CreateTestCollection(name string, description string, client weaviate.Client) {
	ctx := context.Background()
	className := name

	emptyClass := &models.Class{
		Class:           className,
		Description:     description,
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
	// Add error handling
	err := client.Schema().ClassCreator().
		WithClass(emptyClass).
		Do(ctx)

	if err != nil {
		panic(err)
	}
}
