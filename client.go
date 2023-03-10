package builder

import (
	"net/http"
	"time"
)

// defaultHTTPTimeout is the default timeout on the http.Client used by the library.
const defaultHTTPTimeout = 120 * time.Second

// APIURL is the base URL of the Builder API.
const APIURL string = "https://builder.api.reevolute.com"

// clientversion should be the same as tag, used for seeting user-agent.
const clientversion = "0.2.0"

// ResponseData data component from Builder response.
type ResponseData struct {
	Description string                 `json:"description"`
	ErrorCode   string                 `json:"error_code"`
	Vars        map[string]interface{} `json:"vars"`
}

// Response result of Builder execution.
type Response struct {
	SessionID    string
	RequestID    string
	TreeVersion  string
	ResponseType string
	Data         ResponseData
}

// Client interface.
type Client interface {
	AddExecution(treeID, releaseID string, params map[string]interface{}) (Response, error)
	AddAsyncExecution(treeID, releaseID string, params map[string]interface{}) (string, error)
	AddInteraction(sessionID, interactionType string, params map[string]interface{}) (Response, error)
	GetSessionInformation(sessionID string) (Response, error)
}

// API is the builder client implementation.
type API struct {
	httpClient *http.Client
	apiKey     string
	apiURL     string
	TenantID   string
}

// The default HTTP client used for communications with builder.
func getDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultHTTPTimeout,
	}
}

// New creates a new Builder client with the appropriate secret key
// and the tenantID associated.
func New(key string, tenantID string) *API {
	api := API{
		httpClient: getDefaultHTTPClient(),
		apiKey:     key,
		apiURL:     APIURL,
		TenantID:   tenantID,
	}

	return &api
}
