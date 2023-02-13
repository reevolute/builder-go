package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddExection adds single execution to Builder.
func (a *API) AddExecution(treeID, deploymentID string, params map[string]interface{}) (Response, error) {
	baseURL := fmt.Sprintf("%s/v2/tenants/%s/trees/%s/releases/%s/executions",
		APIURL, a.TenantID, treeID, deploymentID)

	var requestBody struct {
		Parameters      map[string]interface{} `json:"parameters"`
		InteractionType string                 `json:"type"`
	}

	requestBody.Parameters = params
	requestBody.InteractionType = "sync"

	body, err := json.Marshal(requestBody)
	if err != nil {
		return Response{}, err
	}

	request, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, baseURL, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}

	return a.builderBaseRequest(context.TODO(), request)
}
