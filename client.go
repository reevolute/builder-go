package builder

import (
	"net/http"
	"time"
)

// defaultHTTPTimeout is the default timeout on the http.Client used by the library.
const defaultHTTPTimeout = 120 * time.Second

// The default HTTP client used for communications with builder.
var httpClient = &http.Client{
	Timeout: defaultHTTPTimeout,
}

type Response struct {
}

// Client interface.
type Client interface {
	AddExecution(treeID, releaseID string, params map[string]interface{}) (Response, error)
	AddAsyncExecution(treeID, releaseID string, params map[string]interface{}) (string, error)
}

// API is the builder client implementation.
type API struct {
	httpClient *http.Client
	apiKey     string
	TenantID   string
}

// New creates a new Builder client with the appropriate secret key
// and the tenantID associated.
func New(key string, tenantID string) *API {
	api := API{
		httpClient: httpClient,
		apiKey:     key,
		TenantID:   tenantID,
	}
	return &api
}
