package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AddInteraction adds an interaction for a session.
func (a *API) AddInteraction(sessionID, interactionType string, params map[string]interface{}) (Response, error) {
	baseURL := fmt.Sprintf("%s/v2/tenants/%s/executions/%s/interactions",
		a.apiURL, a.TenantID, sessionID)

	var requestBody struct {
		Parameters      map[string]interface{} `json:"parameters"`
		InteractionType string                 `json:"type"`
	}

	requestBody.Parameters = params
	requestBody.InteractionType = interactionType

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
