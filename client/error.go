package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/api/googleapi"
)

type ChronicleAPIError struct {
	Message        string `json:"message"`
	Result         string `json:"result"`
	HTTPStatusCode int
}

func (c *ChronicleAPIError) Error() string {
	return fmt.Sprintf("%s: %s, HTTP status code: %d", c.Result, c.Message, c.HTTPStatusCode)
}

func errorForStatusCode(r *http.Response, err error) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	message := ""
	gError, ok := err.(*googleapi.Error)
	if ok {
		message = gError.Message
	}

	apiError := &ChronicleAPIError{
		HTTPStatusCode: r.StatusCode,
		Message:        message,
	}
	res, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(res, &apiError)

	return apiError
}
