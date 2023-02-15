package builder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

type builderError struct {
	Error string `json:"error"`
}

var (
	errBuilderAPI      = errors.New("internal_builder_error")
	errTreeNotFound    = errors.New("tree_not_found")
	errReleaseNotFound = errors.New("release_not_found")
	errInvalidAPIKey   = errors.New("invalidApiKey")
	errTenantNotFound  = errors.New("tenant_not_found")
	errAPIKeyFormat    = errors.New("wrong_api_key_format")
	errPermissions     = errors.New("not_enough_privileges")
	errRateLimit       = errors.New("rate_limit_reached")
)

func proc404(err string) error {
	switch err {
	case "function_not_found":
		return errReleaseNotFound
	case "tree_not_found":
		return errTreeNotFound
	}

	return errBuilderAPI
}

func proc400(err string) error {
	switch err {
	case "authorization header format must be Bearer {token}":
		return errAPIKeyFormat
	case "tree_not_found":
		return errTreeNotFound
	}

	return errBuilderAPI
}

func procBuilderErrors(status int, err string) error {
	switch status {
	case http.StatusNotFound:
		return proc404(err)
	case http.StatusUnauthorized:
		return errInvalidAPIKey
	case http.StatusForbidden:
		return errPermissions
	case http.StatusBadRequest:
		return proc400(err)
	}

	return errBuilderAPI
}

func procBalancerError(status int) error {
	switch status {
	case http.StatusNotFound:
		return errTenantNotFound
	case http.StatusServiceUnavailable:
		return errRateLimit
	}

	return errBuilderAPI
}

func procErrors(response *http.Response, body []byte) error {
	contentType := response.Header.Get("Content-Type")

	if strings.Contains(contentType, "text/html") {
		return procBalancerError(response.StatusCode)
	}

	var res builderError

	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("%w", err)
	}

	return procBuilderErrors(response.StatusCode, res.Error)
}

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

	if unacceptableStatusCode := 399; response.StatusCode > unacceptableStatusCode {
		return Response{}, procErrors(response, content)
	}

	var baseResponse builderResponse

	if err := json.Unmarshal(content, &baseResponse); err != nil {
		return Response{}, fmt.Errorf("%w", err)
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
