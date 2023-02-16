package builder

import (
	"context"
	"fmt"
	"net/http"
)

// GetSessionInformation adds an interaction for a session.
func (a *API) GetSessionInformation(sessionID string) (Response, error) {
	baseURL := fmt.Sprintf("%s/v2/tenants/%s/executions/%s",
		a.apiURL, a.TenantID, sessionID)

	request, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, baseURL, nil)
	if err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}

	return a.builderBaseSyncRequest(context.TODO(), request)
}
