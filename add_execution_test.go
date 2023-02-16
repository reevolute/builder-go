package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddExecutions200(t *testing.T) {
	userAgent := fmt.Sprintf("builder-go/%s", clientversion)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedURL := "/v2/tenants/my_tenant_1312/trees/color_pick/releases/production/executions"
		if r.URL.String() != expectedURL {
			t.Errorf("got [%s] want [%s]", r.URL.String(), expectedURL)
		}

		var requestBody struct {
			Parameters      map[string]interface{} `json:"parameters"`
			InteractionType string                 `json:"type"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Errorf("Error request body %v", err)
		}

		if requestBody.Parameters["color"] != "red" {
			t.Errorf("got [%s] expected [red]", requestBody.Parameters["color"])
		}

		expectedInteractionType := "sync"
		if requestBody.InteractionType != expectedInteractionType {
			t.Errorf("got [%s] expected [%s]", requestBody.InteractionType, expectedInteractionType)
		}

		reqAgent := r.Header.Get("User-Agent")
		if reqAgent != userAgent {
			t.Errorf("want [%s] got [%s]", userAgent, reqAgent)
		}

		serverResponse := []byte(`
		{
		  "tree_version": "3",
		  "response_type": "COMMON",
		  "data": {
		    "description": "function evaluation",
		    "error_code": "0",
		    "vars": {
		      "child_response": "red",
		      "concat_response": "COLOR: rojo"
		    }
		  }
		}
		`)

		w.Header().Set(headerSessionID, "c563cd9a979c46c18d8d892b122f5e38")
		w.Header().Set(headerRequestID, "c563cd9a979c46c18d8d892b122f5e39")
		w.Header().Set("X-Trace-Id", "c563cd9a979c46c18d8d892b122f5e40")

		n, err := w.Write([]byte(serverResponse))
		if err != nil {
			t.Errorf("Error writing response httptest Server [%v][%d]", err, n)
		}
	}))

	defer server.Close()

	apiKEY := "aabbcc"
	tenantID := "my_tenant_1312"

	client := New(apiKEY, tenantID)
	client.apiURL = server.URL

	parameters := map[string]interface{}{
		"color": "red",
	}

	response, err := client.AddExecution("color_pick", "production", parameters)
	if err != nil {
		t.Error(err)
	}

	payloadResponse := Response{
		SessionID:    "c563cd9a979c46c18d8d892b122f5e38",
		RequestID:    "c563cd9a979c46c18d8d892b122f5e39",
		TreeVersion:  "3",
		ResponseType: "COMMON",
		Data: ResponseData{
			Description: "function evaluation",
			ErrorCode:   "0",
			Vars: map[string]interface{}{
				"child_response":  "red",
				"concat_response": "COLOR: rojo",
			},
		},
	}

	if diff := cmp.Diff(payloadResponse, response); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestAddSyncExecutionsError(t *testing.T) {
	userAgent := fmt.Sprintf("builder-go/%s", clientversion)

	cases := []struct {
		name        string
		contentType string
		errMsg      string
		errResult   string
		status      int
	}{
		{
			"tenant not found",
			"text/html",
			"not found",
			errTenantNotFound.Error(),
			http.StatusNotFound,
		},
		{
			"rate limit",
			"text/html",
			"rate limit",
			errRateLimit.Error(),
			http.StatusServiceUnavailable,
		},
		{
			"tree not found",
			"application/json",
			"tree_not_found",
			errTreeNotFound.Error(),
			http.StatusNotFound,
		},
		{
			"release not found",
			"application/json",
			"function_not_found",
			errReleaseNotFound.Error(),
			http.StatusNotFound,
		},
		{
			"api key format",
			"application/json",
			"authorization header format must be Bearer {token}",
			errAPIKeyFormat.Error(),
			http.StatusBadRequest,
		},
		{
			"permissions",
			"application/json",
			"not_allowd",
			errPermissions.Error(),
			http.StatusForbidden,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqAgent := r.Header.Get("User-Agent")
				if reqAgent != userAgent {
					t.Errorf("want [%s] got [%s]", userAgent, reqAgent)
				}

				payload := map[string]string{
					"error": tt.errMsg,
				}

				serverResponse, err := json.Marshal(payload)
				if err != nil {
					t.Errorf("error marshalling response [%v]", err)
				}

				w.Header().Set("Content-Type", tt.contentType)
				w.Header().Set(headerSessionID, "c563cd9a979c46c18d8d892b122f5e38")
				w.Header().Set(headerRequestID, "c563cd9a979c46c18d8d892b122f5e39")
				w.Header().Set("X-Trace-Id", "c563cd9a979c46c18d8d892b122f5e40")
				w.WriteHeader(tt.status)

				n, err := w.Write(serverResponse)
				if err != nil {
					t.Errorf("Error writing response httptest Server [%v][%d]", err, n)
				}
			}))

			defer server.Close()

			apiKEY := "aabbcc"
			tenantID := "play_ground_1234"

			client := New(apiKEY, tenantID)
			client.apiURL = server.URL

			parameters := map[string]interface{}{
				"color": "red",
			}

			treeID := "01GS8E0S"

			response, err := client.AddExecution(treeID, "test", parameters)
			if err == nil {
				t.Error("Tesing error, result must have a non nil error")
			}

			if err.Error() != tt.errResult {
				t.Errorf("want [%s] got [%s]", tt.errResult, err.Error())
			}

			payloadResponse := Response{}

			if diff := cmp.Diff(payloadResponse, response); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAddAsyncExecution201(t *testing.T) {
	userAgent := fmt.Sprintf("builder-go/%s", clientversion)

	expectedRequestID := "c563cd9a979c46c18d8d892b122f5e39"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedURL := "/v2/tenants/my_tenant_1312/trees/color_pick/releases/production/executions"
		if r.URL.String() != expectedURL {
			t.Errorf("got [%s] want [%s]", r.URL.String(), expectedURL)
		}

		var requestBody struct {
			Parameters      map[string]interface{} `json:"parameters"`
			InteractionType string                 `json:"type"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Errorf("Error request body %v", err)
		}

		if requestBody.Parameters["color"] != "red" {
			t.Errorf("got [%s] expected [red]", requestBody.Parameters["color"])
		}

		expectedInteractionType := "async"
		if requestBody.InteractionType != expectedInteractionType {
			t.Errorf("got [%s] expected [%s]", requestBody.InteractionType, expectedInteractionType)
		}

		reqAgent := r.Header.Get("User-Agent")
		if reqAgent != userAgent {
			t.Errorf("want [%s] got [%s]", userAgent, reqAgent)
		}

		w.Header().Set(headerSessionID, "c563cd9a979c46c18d8d892b122f5e38")
		w.Header().Set(headerRequestID, expectedRequestID)
		w.Header().Set("X-Trace-Id", "c563cd9a979c46c18d8d892b122f5e40")
		w.WriteHeader(http.StatusCreated)
	}))

	defer server.Close()

	apiKEY := "aabbcc"
	tenantID := "my_tenant_1312"

	client := New(apiKEY, tenantID)
	client.apiURL = server.URL

	parameters := map[string]interface{}{
		"color": "red",
	}

	txID, err := client.AddAsyncExecution("color_pick", "production", parameters)
	if err != nil {
		t.Error(err)
	}

	if txID != expectedRequestID {
		t.Errorf("got[%s] want[%s]", txID, expectedRequestID)
	}

}

func TestAddAsyncExecutionsError(t *testing.T) {
	userAgent := fmt.Sprintf("builder-go/%s", clientversion)

	cases := []struct {
		name        string
		contentType string
		errMsg      string
		errResult   string
		status      int
	}{
		{
			"tenant not found",
			"text/html",
			"not found",
			errTenantNotFound.Error(),
			http.StatusNotFound,
		},
		{
			"rate limit",
			"text/html",
			"rate limit",
			errRateLimit.Error(),
			http.StatusServiceUnavailable,
		},
		{
			"tree not found",
			"application/json",
			"tree_not_found",
			errTreeNotFound.Error(),
			http.StatusNotFound,
		},
		{
			"release not found",
			"application/json",
			"function_not_found",
			errReleaseNotFound.Error(),
			http.StatusNotFound,
		},
		{
			"api key format",
			"application/json",
			"authorization header format must be Bearer {token}",
			errAPIKeyFormat.Error(),
			http.StatusBadRequest,
		},
		{
			"permissions",
			"application/json",
			"not_allowd",
			errPermissions.Error(),
			http.StatusForbidden,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqAgent := r.Header.Get("User-Agent")
				if reqAgent != userAgent {
					t.Errorf("want [%s] got [%s]", userAgent, reqAgent)
				}

				payload := map[string]string{
					"error": tt.errMsg,
				}

				serverResponse, err := json.Marshal(payload)
				if err != nil {
					t.Errorf("error marshalling response [%v]", err)
				}

				w.Header().Set("Content-Type", tt.contentType)
				w.Header().Set(headerSessionID, "c563cd9a979c46c18d8d892b122f5e38")
				w.Header().Set(headerRequestID, "c563cd9a979c46c18d8d892b122f5e39")
				w.Header().Set("X-Trace-Id", "c563cd9a979c46c18d8d892b122f5e40")
				w.WriteHeader(tt.status)

				n, err := w.Write(serverResponse)
				if err != nil {
					t.Errorf("Error writing response httptest Server [%v][%d]", err, n)
				}
			}))

			defer server.Close()

			apiKEY := "aabbcc"
			tenantID := "play_ground_1234"

			client := New(apiKEY, tenantID)
			client.apiURL = server.URL

			parameters := map[string]interface{}{
				"color": "red",
			}

			treeID := "01GS8E0S"

			response, err := client.AddAsyncExecution(treeID, "test", parameters)
			if err == nil {
				t.Error("Tesing error, result must have a non nil error")
			}

			if err.Error() != tt.errResult {
				t.Errorf("want [%s] got [%s]", tt.errResult, err.Error())
			}
			if response != "" {
				t.Errorf("on error response must be empty")
			}

		})
	}
}
