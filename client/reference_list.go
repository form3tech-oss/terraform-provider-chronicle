package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type ReferenceList struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Lines       []string                 `json:"lines,omitempty"`
	ContentType ReferenceListContentType `json:"content_type,omitempty"`
	CreateTime  string                   `json:"create_time,omitempty"`
}
type ReferenceListResponseContenType struct {
	ContentType ReferenceListContentType `json:"contentType,omitempty"`
}

type ReferenceListResponseCreateTime struct {
	CreateTime string `json:"createTime,omitempty"`
}

type ReferenceListContentType string

const ReferenceListContentTypeDefault ReferenceListContentType = "CONTENT_TYPE_DEFAULT_STRING"
const ReferenceListContentTypeREGEX ReferenceListContentType = "REGEX"
const ReferenceListContentTypeCIDR ReferenceListContentType = "CIDR"

func (cli *Client) GetReferenceList(name string) (*ReferenceList, error) {
	url := fmt.Sprintf("%s/%s", cli.ReferenceListsBasePath, name)

	err := cli.rateLimiters.ReferenceListsGetList.Wait(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while getting reference list %s", name))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "GET", cli.userAgent, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting reference list")
	}

	var referenceList ReferenceList
	err = json.Unmarshal(res, &referenceList)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal reference list response")
	}

	// Set contentType
	var contentTypeResponse ReferenceListResponseContenType
	err = json.Unmarshal(res, &contentTypeResponse)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal contentType response")
	}

	referenceList.ContentType = contentTypeResponse.ContentType
	if referenceList.ContentType == "" {
		referenceList.ContentType = ReferenceListContentTypeDefault
	}

	// Set time
	var createTimeResponse ReferenceListResponseCreateTime
	err = json.Unmarshal(res, &createTimeResponse)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal createTime response")
	}

	referenceList.CreateTime = createTimeResponse.CreateTime

	return &referenceList, nil
}

func (cli *Client) CreateReferenceList(referenceList ReferenceList) (string, error) {
	url := cli.ReferenceListsBasePath

	err := cli.rateLimiters.ReferenceListsCreateList.Wait(context.Background())
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while creating reference list %v", referenceList))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, referenceList)
	if err != nil {
		return "", errors.Wrap(err, "failed creating reference list")
	}

	var referenceListRes ReferenceList
	err = json.Unmarshal(res, &referenceListRes)
	if err != nil {
		return "", errors.Wrap(err, "could not unmarshal reference list response")
	}

	return referenceListRes.Name, nil
}

func (cli *Client) UpdateReferenceList(referenceList ReferenceList, updateLines, updateDescription bool) (*ReferenceList, error) {
	url := fmt.Sprintf("%s?update_mask=%s", cli.ReferenceListsBasePath, CreateReferenceListUpdateMask(updateLines, updateDescription))

	err := cli.rateLimiters.ReferenceListsUpdateList.Wait(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while updating reference list %v", referenceList))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "PATCH", cli.userAgent, url, referenceList)
	if err != nil {
		return nil, errors.Wrap(err, "failed updating reference list")
	}

	var referenceListRes ReferenceList
	err = json.Unmarshal(res, &referenceListRes)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal reference list response")
	}

	return &referenceListRes, nil
}

func CreateReferenceListUpdateMask(updateLines, updateDescription bool) string {
	mask := make([]string, 0)
	if updateLines {
		mask = append(mask, "list.lines")
	}
	if updateDescription {
		mask = append(mask, "list.description")
	}

	return strings.Join(mask, ",")
}
