// The simple example shows how to add an execution.
package main

import (
	"log"
	"os"

	"github.com/reevolute/builder-go"
)

func main() {
	apiKey := os.Getenv("API_KEY")

	tenantID := os.Getenv("TENANT_ID")

	treeID := os.Getenv("TREE_ID")

	client := builder.New(apiKey, tenantID)

	params := map[string]interface{}{
		"color": "red",
	}

	result, err := client.AddExecution(treeID, "production", params)
	if err != nil {
		log.Printf("Error: %v\n", err)

		return
	}

	log.Printf("The result is %v\n", result)
}
