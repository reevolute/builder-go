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
	request.Header.Set("Content-Type", "application/json")

	userAgent := fmt.Sprintf("builder-go/%s", clientversion)

	request.Header.Set("User-Agent", userAgent)

	authorizationValue := fmt.Sprintf("Bearer %s", a.apiKey)
	request.Header.Set("Authorization", authorizationValue)

	response, err := a.httpClient.Do(request.WithContext(ctx))
	if err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("error closing body [%v]", err)
		}
	}()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}

	var baseResponse builderResponse

	if err := json.Unmarshal(content, &baseResponse); err != nil {
		return Response{}, fmt.Errorf("%w", err)
	}

	if unacceptableStatusCode := 399; response.StatusCode > unacceptableStatusCode {
		return Response{}, errBuilderAPI
	}

	res := Response{
		TreeVersion:  baseResponse.TreeVersion,
		Data:         baseResponse.Data,
		ResponseType: baseResponse.ResponseType,
		SessionID:    response.Header.Get(headerSessionID),
		RequestID:    response.Header.Get(headerRequestID),
	}

	return res, nil
}
