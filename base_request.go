package builder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	headerSessionID = "X-Session-Id"
	headerRequestID = "X-Request-Id"
)

type builderResponse struct {
	TreeVersion  string       `json:"tree_version"`
	ResponseType string       `json:"response_type"`
	Data         ResponseData `json:"data"`
}

var errBuilderAPI = errors.New("builder api response different from 2xx")

func (a *API) builderBaseRequest(ctx context.Context, request *http.Request) (Response, error) {
	var res Response

	request.Header.Set("Content-Type", "application/json")

	authorizationValue := fmt.Sprintf("Bearer %s", a.apiKey)
	request.Header.Set("Authorization", authorizationValue)

	response, err := a.httpClient.Do(request.WithContext(ctx))
	if err != nil {
		return res, fmt.Errorf("%w", err)
	}

	res.SessionID = response.Header.Get(headerSessionID)
	res.RequestID = response.Header.Get(headerRequestID)

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("error closing body [%v]", err)
		}
	}()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return res, fmt.Errorf("%w", err)
	}

	var baseResponse builderResponse

	if err := json.Unmarshal(content, &baseResponse); err != nil {
		return res, fmt.Errorf("%w", err)
	}

	if unacceptableStatusCode := 399; response.StatusCode > unacceptableStatusCode {
		return res, errBuilderAPI
	}

	res.TreeVersion = baseResponse.TreeVersion
	res.Data = baseResponse.Data
	res.ResponseType = baseResponse.ResponseType

	return res, nil
}
