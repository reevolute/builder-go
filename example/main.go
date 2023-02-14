// The simple example shows how to add an execution.
package main

import (
	"fmt"
	"os"

	"github.com/reevolute/builder-go"
)

func main() {

	var apiKey = os.Getenv("API_KEY")
	var tenantID = os.Getenv("TENANT_ID")
	var treeID = os.Getenv("TREE_ID")

	client := builder.New(apiKey, tenantID)

	params := map[string]interface{}{
		"color": "red",
	}

	result, err := client.AddExecution(treeID, "production", params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("The result is %v\n", result)
}
