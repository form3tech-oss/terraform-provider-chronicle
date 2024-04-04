package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/avast/retry-go"
	"google.golang.org/api/googleapi"
)

func sendRequest(client *Client, httpClient *http.Client, method, userAgent string, rawurl string, body interface{}) ([]byte, error) {
	reqHeaders := make(http.Header)
	reqHeaders.Set("Content-Type", "application/json")
	reqHeaders.Set("User-Agent", userAgent)

	var res *http.Response
	err := retry.Do(
		func() error {
			var buf bytes.Buffer
			if body != nil {
				err := json.NewEncoder(&buf).Encode(body)
				if err != nil {
					return err
				}
			}

			u, err := addQueryParams(rawurl, map[string]string{"alt": "json"})
			if err != nil {
				return err
			}

			//nolint:all
			req, err := http.NewRequest(method, u, &buf)
			if err != nil {
				return err
			}

			req.Header = reqHeaders
			//nolint:all
			res, err = httpClient.Do(req)
			if err != nil {
				return err
			}

			if err := googleapi.CheckResponse(res); err != nil {
				googleapi.CloseBody(res)
				return errorForStatusCode(res, err)
			}

			return nil
		}, retry.Attempts(client.requestAttempts), retry.DelayType(retry.BackOffDelay), retry.OnRetry(func(n uint, err error) {
			log.Printf("[DEBUG] Retrying request after error: %v", err)
		}),
	)
	if err != nil {
		// Get error from last attempt
		e := err.(retry.Error)
		return nil, e[len(e)-1]
	}

	if res == nil {
		return nil, fmt.Errorf(`unable to parse server response. This is most likely a terraform problem,
		 please file a bug at https://github.com/form3tech-oss/terraform-provider-chronicle/issues`)
	}

	defer googleapi.CloseBody(res)

	if res.StatusCode == 204 {
		return nil, nil
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func addQueryParams(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
