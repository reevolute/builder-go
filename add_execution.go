package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddExecution adds single execution to Builder.
func (a *API) AddExecution(treeID, deploymentID string, params map[string]interface{}) (Response, error) {
	baseURL := fmt.Sprintf("%s/v2/tenants/%s/trees/%s/releases/%s/executions",
		a.apiURL, a.TenantID, treeID, deploymentID)

	var requestBody struct {
		Parameters      map[string]interface{} `json:"parameters"`
		InteractionType string                 `json:"type"`
	}

	requestBody.Parameters = params
	requestBody.InteractionType = "sync"

	body, err := json.Marshal(requestBody)
	if err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}

	request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}

	return a.builderBaseSyncRequest(context.TODO(), request)
}

// AddAsyncExecution adds single execution to Builder.
func (a *API) AddAsyncExecution(treeID, deploymentID string, params map[string]interface{}) (string, error) {
	baseURL := fmt.Sprintf("%s/v2/tenants/%s/trees/%s/releases/%s/executions",
		a.apiURL, a.TenantID, treeID, deploymentID)

	var requestBody struct {
		Parameters      map[string]interface{} `json:"parameters"`
		InteractionType string                 `json:"type"`
	}

	requestBody.Parameters = params
	requestBody.InteractionType = "async"

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return a.builderBaseAsyncRequest(context.TODO(), request)
}
