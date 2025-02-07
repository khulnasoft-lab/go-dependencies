package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPError is raised when status code is not 2xx
type HTTPError struct {
	HTTPStatusCode int
	HTTPStatus     string
	Message        string      `json:"message,omitempty"`
	Code           string      `json:"code,omitempty"`
	MoreInfo       string      `json:"moreInfo,omitempty"`
	Details        interface{} `json:"details,omitempty"`
}

// This allows HTTPError to satisfy the error interface
func (he *HTTPError) Error() string {
	jsonErrMsg, err := json.Marshal(he)
	if err != nil {
		return fmt.Sprintf("Http Error - HTTPStatusCode: [%v], HTTPStatus: %v, Error Message: %v, Error Code: %v",
			he.HTTPStatusCode, he.HTTPStatus, he.Message, he.Code)
	}
	return string(jsonErrMsg)
}

// ParseHTTPStatusCodeInResponse returns http response and HTTPError struct based on response status code
func ParseHTTPStatusCodeInResponse(response *http.Response) (*http.Response, error) {
	if response != nil && (response.StatusCode < 200 || response.StatusCode >= 400) {
		httpErr := HTTPError{
			HTTPStatusCode: response.StatusCode,
			HTTPStatus:     response.Status,
		}
		if response.Body != nil {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return response, &httpErr
			}
			err = json.Unmarshal(body, &httpErr)
			if err != nil {
				return nil, &httpErr
			}
		}
		return response, &httpErr
	}
	return response, nil
}
