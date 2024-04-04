package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	FeedSourceTypeAPI            = "API"
	FeedSourceTypeAzureBlobStore = "AZURE_BLOBSTORE"
	FeedSourceTypeGCS            = "GOOGLE_CLOUD_STORAGE"
	FeedSourceTypeS3             = "AMAZON_S3"
	FeedSourceTypeSQS            = "AMAZON_SQS"
	FeedSourceTypeHTTP           = "HTTP"
)

type ConcreteFeedConfiguration interface {
	getConfigurationPropertyKey() string
	getFeedSourceType() string
}

type BaseFeed struct {
	Name        string      `json:"name,omitempty"`
	DisplayName string      `json:"display_name,omitempty"`
	Details     FeedDetails `json:"details,omitempty"`
	State       string      `json:"feedState,omitempty"`
}

type FeedDetails struct {
	SourceType string  `json:"feedSourceType,omitempty"`
	LogType    string  `json:"logType,omitempty"`
	Namespace  string  `json:"namespace,omitempty"`
	Labels     []Label `json:"labels,omitempty"`
}

type Label struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func concreteFeedConfigurationToMap(c ConcreteFeedConfiguration) (map[string]interface{}, error) {
	return toMapWithJSONTags(c)
}
func expandFeedFromFeedMap(configurationPropertyKey string, feedMap map[string]interface{}) (*BaseFeed, ConcreteFeedConfiguration, error) {
	baseFeed, concreteFeed, err := fromFeedMapToBaseFeedAndConcreteConfiguration(configurationPropertyKey, feedMap)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed generating feed")
	}

	return baseFeed, concreteFeed, nil
}

func newFeedAsMapFromConcreteFeed(name, displayName, logType, namespace string,
	labels []Label, concreteFeed ConcreteFeedConfiguration) (map[string]interface{}, error) {
	feed := BaseFeed{
		Name:        name,
		DisplayName: displayName,
		Details: FeedDetails{
			LogType:    logType,
			SourceType: concreteFeed.getFeedSourceType(),
			Namespace:  namespace,
			Labels:     labels,
		},
	}

	feedMap, err := toMapWithJSONTags(feed)
	if err != nil {
		return nil, err
	}
	concreteFeedMap, err := concreteFeedConfigurationToMap(concreteFeed)
	if err != nil {
		return nil, err
	}

	return joinBaseFeedMapAndConcreteFeedMap(concreteFeed.getConfigurationPropertyKey(), feedMap, concreteFeedMap), nil
}

func joinBaseFeedMapAndConcreteFeedMap(configurationPropertyKey string, baseFeedMap, concreteFeedMap map[string]interface{}) map[string]interface{} {
	details := baseFeedMap["details"].(map[string]interface{})
	var err error
	details[configurationPropertyKey], err = toMapWithJSONTags(concreteFeedMap)
	if err != nil {
		return nil
	}
	baseFeedMap["details"] = details

	return baseFeedMap
}

func (cli *Client) CreateFeed(displayName, logType, namespace string, labels []Label, concreteFeedConfiguration ConcreteFeedConfiguration) (string, error) {
	feed, err := newFeedAsMapFromConcreteFeed("", displayName, logType, namespace, labels, concreteFeedConfiguration)
	if err != nil {
		return "", errors.Wrap(err, "failed generating feed")
	}

	url := cli.FeedManagementBasePath

	err = cli.rateLimiters.FeedManagementCreateFeed.Wait(context.Background())
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while creating feed %s", displayName))
	}

	res, err := sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, feed)
	if err != nil {
		return "", errors.Wrap(err, "failed creating feed")
	}

	reader := bytes.NewReader(res)
	result := make(map[string]interface{})
	if err := json.NewDecoder(reader).Decode(&result); err != nil {
		return "", errors.Wrap(err, "failed decoding feed")
	}

	name := parseFeedID(result["name"].(string))

	return name, nil
}

func (cli *Client) UpdateFeed(name, displayName, logType, namespace string, labels []Label, conf ConcreteFeedConfiguration) error {
	feed, err := newFeedAsMapFromConcreteFeed(name, displayName, logType, namespace, labels, conf)
	if err != nil {
		return errors.Wrap(err, "failed updating feed")
	}

	name, ok := feed["name"].(string)
	if !ok {
		return fmt.Errorf("cannot get property name from feed")
	}

	url := fmt.Sprintf("%s/%s", cli.FeedManagementBasePath, name)

	err = cli.rateLimiters.FeedManagementUpdateFeed.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error Waiting for rateLimiter while updating feed %s", name))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "PATCH", cli.userAgent, url, feed)
	if err != nil {
		return errors.Wrap(err, "failed creating feed")
	}

	return nil
}

func (cli *Client) ReadFeed(name string) (*BaseFeed, *ConcreteFeedConfiguration, error) {
	url := fmt.Sprintf("%s/%s", cli.FeedManagementBasePath, name)

	err := cli.rateLimiters.FeedManagementGetFeed.Wait(context.Background())
	if err != nil {
		return nil, nil, errors.Wrap(err, fmt.Sprintf("Error Waiting for rateLimiter while reading feed %s", name))
	}
	res, err := sendRequest(cli, cli.backstoryAPIClient, "GET", cli.userAgent, url, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed reading feed")
	}

	reader := bytes.NewReader(res)
	result := make(map[string]interface{})
	if err := json.NewDecoder(reader).Decode(&result); err != nil {
		return nil, nil, errors.Wrap(err, "failed decoding feed")
	}
	details := result["details"].(map[string]interface{})

	feedSourceType := extractFeedSourceTypeFromDetails(details)
	logType := extractLogTypeFromDetails(details)
	concreteFeed := newConcreteFeedConfiguration(feedSourceType, logType)

	var baseFeed *BaseFeed
	baseFeed, concreteFeed, err = expandFeedFromFeedMap(concreteFeed.getConfigurationPropertyKey(), result)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed generating feed")
	}

	baseFeed.Name = parseFeedID(baseFeed.Name)

	return baseFeed, &concreteFeed, nil
}

func (cli *Client) DestroyFeed(name string) error {
	url := fmt.Sprintf("%s/%s", cli.FeedManagementBasePath, name)

	err := cli.rateLimiters.FeedManagementDeleteFeed.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while destroying feed %s", name))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "DELETE", cli.userAgent, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed deleting feed")
	}

	return nil
}

func (cli *Client) ChangeEnableFeed(id string, enabled bool) error {
	var operation string
	if enabled {
		operation = "enable"
	} else {
		operation = "disable"
	}

	url := fmt.Sprintf("%s/%s:%s", cli.FeedManagementBasePath, id, operation)

	err := cli.rateLimiters.FeedManagementEnableFeed.Wait(context.Background())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error waiting for rateLimiter while change enable feed %s", id))
	}

	_, err = sendRequest(cli, cli.backstoryAPIClient, "POST", cli.userAgent, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed creating feed")
	}

	return nil
}

func fromFeedMapToBaseFeedAndConcreteConfiguration(configurationPropertyKey string, feedMap map[string]interface{}) (*BaseFeed, ConcreteFeedConfiguration, error) {
	details := feedMap["details"].(map[string]interface{})

	var baseFeed *BaseFeed
	baseFeedBytes, err := json.Marshal(feedMap)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed generating feed")
	}
	err = json.Unmarshal(baseFeedBytes, &baseFeed)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed generating feed")
	}

	concreteFeedBytes, err := json.Marshal(details[configurationPropertyKey])
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed generating feed")
	}

	feedSourceType := extractFeedSourceTypeFromDetails(details)
	logType := extractLogTypeFromDetails(details)
	concreteFeed := newConcreteFeedConfiguration(feedSourceType, logType)
	err = json.Unmarshal(concreteFeedBytes, &concreteFeed)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed generating feed")
	}

	return baseFeed, concreteFeed, nil
}

func parseFeedID(idRaw string) string {
	return strings.ReplaceAll(idRaw, "feeds/", "")
}

func newConcreteFeedConfiguration(feedSourceType, logType string) ConcreteFeedConfiguration {
	switch feedSourceType {
	case FeedSourceTypeS3:
		return &S3FeedConfiguration{}
	case FeedSourceTypeSQS:
		return &SQSFeedConfiguration{}
	case FeedSourceTypeAPI:
		return newAPIConcreteFeedConfigurationFromLogType(logType)
	case FeedSourceTypeGCS:
		return &GCPBucketFeedConfiguration{}
	case FeedSourceTypeAzureBlobStore:
		return &AzureBlobStoreFeedConfiguration{}

	default:
		return nil
	}
}

func newAPIConcreteFeedConfigurationFromLogType(logType string) ConcreteFeedConfiguration {
	switch logType {

	case "AZURE_AD_AUDIT", "AZURE_AD_CONTEXT", "AZURE_AD", "AZURE_MDM_INTUNE", "MICROSOFT_GRAPH_ALERT":
		return &MicrosoftGraphAPIFeedConfiguration{ContentType: logType}

	case MicrosoftOffice365ManagementActivityFeedLogType:
		return &MicrosoftOffice365ManagementActivityFeedConfiguration{}

	case OktaUsersFeedLogType:
		return &OktaUsersFeedConfiguration{}

	case QualysVMFeedLogType:
		return &QualysVMFeedConfiguration{}

	case OktaSystemLogFeedLogType:
		return &OktaSystemLogFeedConfiguration{}

	case ProofpointSIEMFeedLogType:
		return &ProofpointSIEMFeedConfiguration{}

	case ThinkstCanaryFeedLogType:
		return &ThinkstCanaryFeedConfiguration{}

	default:
		return nil
	}
}

func extractFeedSourceTypeFromDetails(details map[string]interface{}) string {
	return details["feedSourceType"].(string)
}
func extractLogTypeFromDetails(details map[string]interface{}) string {
	return details["logType"].(string)
}
