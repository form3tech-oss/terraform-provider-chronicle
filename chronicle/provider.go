package chronicle

import (
	"context"
	"fmt"
	"time"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/form3tech-oss/terraform-provider-chronicle/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateRegion,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_REGION",
				}, chronicle.RegionEurope),
				Description: fmt.Sprintf(`Region to which send requests, available regions are: %v. It may be replaced by CHRONICLE_REGION environment variable.`, chronicle.Regions),
			},

			"bigqueryapi_credentials": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateCredentials,
				ConflictsWith:    []string{"bigqueryapi_access_token"},
				Description: `BigQuery API crendential. Local file path or content.
				 It may be replaced by CHRONICLE_BIGQUERY_CREDENTIALS environment variable, which expects base64 encoded credential.`,
			},

			"bigqueryapi_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bigqueryapi_credentials"},
				Description:   `BigQuery API access token. Local file path or content.`,
			},
			"backstoryapi_credentials": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateCredentials,
				ConflictsWith:    []string{"backstoryapi_access_token"},
				Description: `Backstory API credential. Local file path or content.
				 It may be replaced by CHRONICLE_BACKSTORY_CREDENTIALS environment variable, which expects base64 encoded credential.`,
			},

			"backstoryapi_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"backstoryapi_credentials"},
				Description:   `Backstory API access token. Local file path or content.`,
			},
			"ingestionapi_credentials": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateCredentials,
				ConflictsWith:    []string{"ingestionapi_access_token"},
				Description: `Ingestion API crendential. Local file path or content.
				 It may be replaced by CHRONICLE_INGESTION_CREDENTIALS environment variable, which expects base64 encoded credential.`,
			},

			"ingestionapi_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ingestionapi_credentials"},
				Description:   `Ingestion API access token. Local file path or content.`,
			},
			"forwarderapi_credentials": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateCredentials,
				ConflictsWith:    []string{"forwarderapi_access_token"},
				Description: `Forwarder API crendential. Local file path or content.
				 It may be replaced by CHRONICLE_FORWARDER_CREDENTIALS environment variable, which expects base64 encoded credential.`,
			},
			"forwarderapi_access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"forwarderapi_credentials"},
				Description:   `Forwarder API Access token. Local file path or content.`,
			},

			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Request timeout in seconds. Defaults to 120 (s).`,
				Default:     120,
			},
			"request_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Number of attempts per request. Attempts follow exponential back-off strategy. Defaults to 5 attempts.`,
				Default:     5,
			},

			"events_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to events endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_EVENTS_CUSTOM_ENDPOINT",
				}, nil),
			},

			"alert_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to alert endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_ALERT_CUSTOM_ENDPOINT",
				}, nil),
			},

			"artifact_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to artifact endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_ARTIFACT_CUSTOM_ENDPOINT",
				}, nil),
			},

			"alias_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to alias endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_ALIAS_CUSTOM_ENDPOINT",
				}, nil),
			},

			"asset_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to asset endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_ASSET_CUSTOM_ENDPOINT",
				}, nil),
			},

			"ioc_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to ioc endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_IOC_CUSTOM_ENDPOINT",
				}, nil),
			},

			"rule_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to rule endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_RULE_CUSTOM_ENDPOINT",
				}, nil),
			},

			"feed_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to feed endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_FEED_CUSTOM_ENDPOINT",
				}, nil),
			},

			"subjects_custom_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `Custom URL to subjects endpoint.`,
				ValidateDiagFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHRONICLE_SUBJECTS_CUSTOM_ENDPOINT",
				}, nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"chronicle_rbac_subject":                                  resourceRBACSubject(),
			"chronicle_rule":                                          resourceRule(),
			"chronicle_reference_list":                                resourceReferenceList(),
			"chronicle_feed_amazon_s3":                                NewResourceFeedAmazonS3().TerraformResource,
			"chronicle_feed_amazon_s3_v2":                             NewResourceFeedAmazonS3V2().TerraformResource,
			"chronicle_feed_amazon_sqs":                               NewResourceFeedAmazonSQS().TerraformResource,
			"chronicle_feed_amazon_sqs_v2":                            NewResourceFeedAmazonSQSV2().TerraformResource,
			"chronicle_feed_qualys_vm":                                NewResourceFeedQualysVM().TerraformResource,
			"chronicle_feed_microsoft_office_365_management_activity": NewResourceFeedMicrosoftOffice365ManagementActivity().TerraformResource,
			"chronicle_feed_okta_system_log":                          NewResourceFeedOktaSystemLog().TerraformResource,
			"chronicle_feed_okta_users":                               NewResourceFeedOktaUsers().TerraformResource,
			"chronicle_feed_proofpoint_siem":                          NewResourceFeedProofpointSIEM().TerraformResource,
			"chronicle_feed_google_cloud_storage_bucket":              NewResourceFeedGoogleCloudStorageBucket().TerraformResource,
			"chronicle_feed_google_cloud_storage_v2":                  NewResourceFeedGoogleCloudStorageV2().TerraformResource,
			"chronicle_feed_azure_blobstore":                          NewResourceFeedAzureBlobStore().TerraformResource,
			"chronicle_feed_azure_blobstore_v2":                       NewResourceFeedAzureBlobStoreV2().TerraformResource,
			"chronicle_feed_thinkst_canary":                           NewResourceFeedThinkstCanary().TerraformResource,
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, provider)
	}

	return provider
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	var region string
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	} else {
		region = envSearch(chronicle.ChronicleRegionEnvVar)
		if region == "" {
			region = chronicle.RegionEurope
		} else if !isValidRegion(region) {
			return nil, diag.FromErr(fmt.Errorf("region %s is not valid", region))
		}
	}

	opts := getAPIAuthOpts(d)

	if v, ok := d.GetOk("request_timeout"); ok {
		opts = append(opts, chronicle.WithRequestTimeout(time.Duration(v.(int))*time.Second))
	}
	if v, ok := d.GetOk("request_attempts"); ok {
		attempts := v.(int)
		if attempts < 0 {
			return nil, diag.FromErr(fmt.Errorf("request_attempts must be non-negative"))
		}
		opts = append(opts, chronicle.WithRequestAttempts(uint(attempts)))
	}

	//nolint:all
	stopCtx, ok := schema.StopContext(ctx)
	if !ok {
		stopCtx = ctx
	}

	client, err := chronicle.NewClient(region, p.UserAgent("terraform-provider-chronicle", version.ProviderVersion), stopCtx, opts...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if endpoint, isCustom := customEndpoint(d, "events_custom_endpoint"); isCustom {
		client.WithEventsBasePath(endpoint)
	}
	if endpoint, isCustom := customEndpoint(d, "alert_custom_endpoint"); isCustom {
		client.WithAlertBasePath(endpoint)
	}
	if endpoint, isCustom := customEndpoint(d, "artifact_custom_endpoint"); isCustom {
		client.WithArtifactBasePath(endpoint)
	}
	if endpoint, isCustom := customEndpoint(d, "alias_custom_endpoint"); isCustom {
		client.WithAliasBasePath(endpoint)
	}
	if endpoint, isCustom := customEndpoint(d, "asset_custom_endpoint"); isCustom {
		client.WithAssetBasePath(endpoint)
	}
	if endpoint, isCustom := customEndpoint(d, "ioc_custom_endpoint"); isCustom {
		client.WithIOCBasePath(endpoint)
	}
	if endpoint, isCustom := customEndpoint(d, "rule_custom_endpoint"); isCustom {
		client.WithRuleBasePath(endpoint)
	}
	if endpoint, isCustom := customEndpoint(d, "subjects_custom_endpoint"); isCustom {
		client.WithSubjectsBasePath(endpoint)
	}

	return client, nil
}

func getAPIAuthOpts(d *schema.ResourceData) []chronicle.Option {
	opts := make([]chronicle.Option, 0)

	if v, ok := d.GetOk("bigqueryapi_credentials"); ok {
		opts = append(opts, chronicle.WithBigQueryAPICredentials(v.(string)))
	} else if v, ok := d.GetOk("bigqueryapi_access_token"); ok {
		opts = append(opts, chronicle.WithBigQueryAPIAccessToken(v.(string)))
	} else {
		env := envSearch(chronicle.BigQueryAPIEnvVar)
		if env != "" {
			opts = append(opts, chronicle.WithBigQueryAPIEnvVar())
		}
	}

	if v, ok := d.GetOk("backstoryapi_credentials"); ok {
		opts = append(opts, chronicle.WithBackstoryAPICredentials(v.(string)))
	} else if v, ok := d.GetOk("backstoryapi_credentials"); ok {
		opts = append(opts, chronicle.WithBackstoryAPIAccessToken(v.(string)))
	} else {
		env := envSearch(chronicle.BackstoryAPIEnvVar)
		if env != "" {
			opts = append(opts, chronicle.WithBackstoryAPIEnvVar())
		}
	}

	if v, ok := d.GetOk("ingestionapi_credentials"); ok {
		opts = append(opts, chronicle.WithIngestionAPICredentials(v.(string)))
	} else if v, ok := d.GetOk("ingestionapi_credentials"); ok {
		opts = append(opts, chronicle.WithIngestionAPIAccessToken(v.(string)))
	} else {
		env := envSearch(chronicle.IngestionAPIEnvVar)
		if env != "" {
			opts = append(opts, chronicle.WithIngestionAPIEnvVar())
		}
	}

	if v, ok := d.GetOk("forwarderapi_credentials"); ok {
		opts = append(opts, chronicle.WithForwarderAPICredentials(v.(string)))
	} else if v, ok := d.GetOk("forwarderapi_access_token"); ok {
		opts = append(opts, chronicle.WithForwarderAPIAccessToken(v.(string)))
	} else {
		env := envSearch(chronicle.ForwarderAPIEnvVar)
		if env != "" {
			opts = append(opts, chronicle.WithBigQueryAPIEnvVar())
		}
	}

	return opts
}

func customEndpoint(d *schema.ResourceData, endpoint string) (string, bool) {
	custom, ok := d.GetOk(endpoint)
	if ok {
		return custom.(string), true
	}
	return "", false
}
