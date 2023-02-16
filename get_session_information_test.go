package builder

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSessionInfo(t *testing.T) {
	userAgent := fmt.Sprintf("builder-go/%s", clientversion)

	sessionID := "c563cd9a979c46c18d8d892b122f5e38"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedURL := fmt.Sprintf("/v2/tenants/my_tenant_1312/executions/%s", sessionID)

		if r.URL.String() != expectedURL {
			t.Errorf("got [%s] want [%s]", r.URL.String(), expectedURL)
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

	response, err := client.GetSessionInformation(sessionID)
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
