package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

type Client struct {
	userAgent       string
	requestAttempts uint
	requestTimeout  time.Duration
	context         context.Context
	rateLimiters    ClientRateLimiters

	bigQueryAPIClient  *http.Client
	backstoryAPIClient *http.Client
	ingestionAPIClient *http.Client
	forwarderAPIClient *http.Client

	EventsBasePath         string
	AlertBasePath          string
	ArtifactBasePath       string
	AliasBasePath          string
	AssetBasePath          string
	IOCBasePath            string
	RuleBasePath           string
	FeedManagementBasePath string
	SubjectsBasePath       string
	ReferenceListsBasePath string
}

type Option func(*Client) error

const (
	defaultRequestAttempts = 5
	defaultRequestTimeout  = time.Second * 120
)

var defaultClientScopes = []string{
	"https://www.googleapis.com/auth/chronicle-backstory",
	"https://www.googleapis.com/auth/malachite-ingestion",
}

func NewClient(region string, userAgent string, ctx context.Context, opts ...Option) (*Client, error) {
	defaultBasePaths := GenerateDefaultBasePaths(region)

	client := &Client{
		userAgent:       userAgent,
		requestAttempts: defaultRequestAttempts,
		requestTimeout:  defaultRequestTimeout,
		context:         ctx,
		rateLimiters:    *NewClientRateLimiters(),

		EventsBasePath:         defaultBasePaths[EventsBasePathKey],
		AlertBasePath:          defaultBasePaths[AlertBasePathKey],
		ArtifactBasePath:       defaultBasePaths[ArtifactBasePathKey],
		AliasBasePath:          defaultBasePaths[AliasBasePathKey],
		AssetBasePath:          defaultBasePaths[AssetBasePathKey],
		IOCBasePath:            defaultBasePaths[IOCBasePathKey],
		RuleBasePath:           defaultBasePaths[RuleBasePathKey],
		FeedManagementBasePath: defaultBasePaths[FeedManagementBasePathKey],
		SubjectsBasePath:       defaultBasePaths[SubjectsBasePathKey],
		ReferenceListsBasePath: defaultBasePaths[ReferenceListsPathKey],
	}

	for _, opt := range opts {
		err := opt(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (cli *Client) WithEventsBasePath(uri string) *Client {
	cli.EventsBasePath = uri
	return cli
}

func (cli *Client) WithAlertBasePath(uri string) *Client {
	cli.AlertBasePath = uri
	return cli
}

func (cli *Client) WithArtifactBasePath(uri string) *Client {
	cli.ArtifactBasePath = uri
	return cli
}

func (cli *Client) WithAliasBasePath(uri string) *Client {
	cli.AliasBasePath = uri
	return cli
}

func (cli *Client) WithAssetBasePath(uri string) *Client {
	cli.AssetBasePath = uri
	return cli
}

func (cli *Client) WithIOCBasePath(uri string) *Client {
	cli.IOCBasePath = uri
	return cli
}

func (cli *Client) WithRuleBasePath(uri string) *Client {
	cli.RuleBasePath = uri
	return cli
}

func (cli *Client) WithFeedManagementBasePath(uri string) *Client {
	cli.FeedManagementBasePath = uri
	return cli
}

func (cli *Client) WithSubjectsBasePath(uri string) *Client {
	cli.SubjectsBasePath = uri
	return cli
}

func WithBigQueryAPICredentials(credentials string) Option {
	return func(cli *Client) error {
		var err error
		cli.bigQueryAPIClient, err = cli.initHTTPClient(defaultClientScopes, "", credentials, "")
		return err
	}
}

func WithBigQueryAPIAccessToken(accesstoken string) Option {
	return func(cli *Client) error {
		var err error
		cli.bigQueryAPIClient, err = cli.initHTTPClient(defaultClientScopes, accesstoken, "", "")
		return err
	}
}

func WithBigQueryAPIEnvVar() Option {
	return func(cli *Client) error {
		var err error
		cli.bigQueryAPIClient, err = cli.initHTTPClient(defaultClientScopes, "", "", BigQueryAPIEnvVar)
		return err
	}
}

func WithBackstoryAPICredentials(credentials string) Option {
	return func(cli *Client) error {
		var err error
		cli.backstoryAPIClient, err = cli.initHTTPClient(defaultClientScopes, "", credentials, "")
		return err
	}
}

func WithBackstoryAPIAccessToken(accesstoken string) Option {
	return func(cli *Client) error {
		var err error
		cli.backstoryAPIClient, err = cli.initHTTPClient(defaultClientScopes, accesstoken, "", "")
		return err
	}
}

func WithBackstoryAPIEnvVar() Option {
	return func(cli *Client) error {
		var err error
		cli.backstoryAPIClient, err = cli.initHTTPClient(defaultClientScopes, "", "", BackstoryAPIEnvVar)
		return err
	}
}

func WithIngestionAPICredentials(credentials string) Option {
	return func(cli *Client) error {
		var err error
		cli.ingestionAPIClient, err = cli.initHTTPClient(defaultClientScopes, "", credentials, "")
		return err
	}
}

func WithIngestionAPIAccessToken(accesstoken string) Option {
	return func(cli *Client) error {
		var err error
		cli.ingestionAPIClient, err = cli.initHTTPClient(defaultClientScopes, accesstoken, "", "")
		return err
	}
}

func WithIngestionAPIEnvVar() Option {
	return func(cli *Client) error {
		var err error
		cli.ingestionAPIClient, err = cli.initHTTPClient(defaultClientScopes, "", "", IngestionAPIEnvVar)
		return err
	}
}

func WithForwarderAPICredentials(credentials string) Option {
	return func(cli *Client) error {
		var err error
		cli.forwarderAPIClient, err = cli.initHTTPClient(defaultClientScopes, "", credentials, "")
		return err
	}
}

func WithForwarderAPIAccessToken(accesstoken string) Option {
	return func(cli *Client) error {
		var err error
		cli.forwarderAPIClient, err = cli.initHTTPClient(defaultClientScopes, accesstoken, "", "")
		return err
	}
}

func WithForwarderAPIEnvVar() Option {
	return func(cli *Client) error {
		var err error
		cli.forwarderAPIClient, err = cli.initHTTPClient(nil, "", "", ForwarderAPIEnvVar)
		return err
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(cli *Client) error {
		cli.requestTimeout = timeout
		return nil
	}
}

func WithRequestAttempts(attempts uint) Option {
	return func(cli *Client) error {
		cli.requestAttempts = attempts
		return nil
	}
}

func (cli *Client) initHTTPClient(scopes []string, accesstoken, credentials, envVariable string) (*http.Client, error) {
	tokenSource, err := cli.getTokenSource(scopes, accesstoken, credentials, envVariable)
	if err != nil {
		return nil, err
	}

	cleanCtx := context.WithValue(cli.context, oauth2.HTTPClient, cleanhttp.DefaultClient())

	client, _, err := transport.NewHTTPClient(cleanCtx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}

	client.Timeout = cli.requestTimeout
	loggingTransport := logging.NewSubsystemLoggingHTTPTransport("Chronicle", client.Transport)
	client.Transport = loggingTransport

	return client, nil
}

func (cli *Client) getTokenSource(scopes []string, accesstoken, credentials, envVariable string) (oauth2.TokenSource, error) {
	creds, err := cli.GetCredentials(scopes, accesstoken, credentials, envVariable)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return creds.TokenSource, nil
}

func (cli *Client) GetCredentials(clientScopes []string, accessToken, credentials, envVariable string) (*googleoauth.Credentials, error) {
	if accessToken != "" {
		contents, _, err := pathOrContents(accessToken)
		if err != nil {
			return &googleoauth.Credentials{}, fmt.Errorf("error loading access token: %s", err)
		}

		token := &oauth2.Token{AccessToken: contents}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return &googleoauth.Credentials{
			TokenSource: oauth2.StaticTokenSource(token),
		}, nil
	}

	if credentials != "" {
		contents, _, err := pathOrContents(credentials)
		if err != nil {
			return &googleoauth.Credentials{}, fmt.Errorf("error loading credentials: %s", err)
		}

		creds, err := googleoauth.CredentialsFromJSON(cli.context, []byte(contents), clientScopes...)
		if err != nil {
			return &googleoauth.Credentials{}, fmt.Errorf("unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return creds, nil
	}

	env := envSearch(envVariable)
	if env != "" {
		envDecoded, err := base64.StdEncoding.DecodeString(env)
		if err != nil {
			return &googleoauth.Credentials{}, fmt.Errorf("unable to base64 decode credentials from '%s': %s", envVariable, err)
		}

		log.Printf("[INFO] Authenticating using environmental variable")

		creds, err := googleoauth.CredentialsFromJSON(cli.context, envDecoded, clientScopes...)
		if err != nil {
			return &googleoauth.Credentials{}, fmt.Errorf("unable to parse credentials from '%s': %s", envVariable, err)
		}

		log.Printf("[INFO] Authenticating using environmental variable'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return creds, nil
	}

	return &googleoauth.Credentials{}, fmt.Errorf("error loading credentials: env variable not found: %s", env)
}
