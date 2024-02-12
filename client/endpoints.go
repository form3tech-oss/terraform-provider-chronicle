package client

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

type ClientRateLimiters struct {
	FeedManagementCreateFeed *rate.Limiter
	FeedManagementGetFeed    *rate.Limiter
	FeedManagementListFeeds  *rate.Limiter
	FeedManagementUpdateFeed *rate.Limiter
	FeedManagementDeleteFeed *rate.Limiter
	FeedManagementEnableFeed *rate.Limiter

	DetectionCreateRule         *rate.Limiter
	DetectionCreateRuleVersion  *rate.Limiter
	DetectionGetRule            *rate.Limiter
	DetectionUpdateRule         *rate.Limiter
	DetectionDeleteRule         *rate.Limiter
	DetectionEnableLiveRule     *rate.Limiter
	DetectionEnableAlertingRule *rate.Limiter
	DetectionVerifyYARARule     *rate.Limiter

	RBACCreateSubject *rate.Limiter
	RBACGetSubject    *rate.Limiter
	RBACUpdateSubject *rate.Limiter
	RBACDeleteSubject *rate.Limiter

	ReferenceListsCreateList *rate.Limiter
	ReferenceListsGetList    *rate.Limiter
	ReferenceListsUpdateList *rate.Limiter
}

func NewClientRateLimiters() *ClientRateLimiters {
	return &ClientRateLimiters{
		FeedManagementCreateFeed: rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementGetFeed:    rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementListFeeds:  rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementUpdateFeed: rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementDeleteFeed: rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementEnableFeed: rate.NewLimiter(rate.Every(time.Second), 1),

		DetectionCreateRule:         rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionCreateRuleVersion:  rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionGetRule:            rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionUpdateRule:         rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionDeleteRule:         rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionEnableLiveRule:     rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionEnableAlertingRule: rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionVerifyYARARule:     rate.NewLimiter(rate.Every(time.Second), 1),

		RBACCreateSubject: rate.NewLimiter(rate.Every(time.Second), 1),
		RBACGetSubject:    rate.NewLimiter(rate.Every(time.Second), 1),
		RBACUpdateSubject: rate.NewLimiter(rate.Every(time.Second), 1),
		RBACDeleteSubject: rate.NewLimiter(rate.Every(time.Second), 1),

		ReferenceListsCreateList: rate.NewLimiter(rate.Every(time.Second), 1),
		ReferenceListsGetList:    rate.NewLimiter(rate.Every(time.Second), 1),
		ReferenceListsUpdateList: rate.NewLimiter(rate.Every(time.Second), 1),
	}
}

const (
	BigQueryAPIEnvVar     = "CHRONICLE_BIGQUERY_CREDENTIALS"
	BackstoryAPIEnvVar    = "CHRONICLE_BACKSTORY_CREDENTIALS"
	IngestionAPIEnvVar    = "CHRONICLE_INGESTION_CREDENTIALS"
	ForwarderAPIEnvVar    = "CHRONICLE_FORWARDER_CREDENTIALS"
	ChronicleRegionEnvVar = "CHRONICLE_REGION"
)

var EnvAPICrendetialsVar = []string{BigQueryAPIEnvVar, BackstoryAPIEnvVar, IngestionAPIEnvVar, ForwarderAPIEnvVar}

const (
	RegionUS             = "us"
	RegionEurope         = "europe"
	RegionEuropeWest2    = "europe-west2"
	RegionAsiaSouthEast1 = "asia-southeast1"
)

var Regions = []string{RegionUS, RegionEurope, RegionEuropeWest2, RegionAsiaSouthEast1}

const APIDomain = "googleapis.com"

const (
	SearchAPIKey          = "SerachAPI"
	DetectionEngineAPIKey = "DetectionEngineAPI"
	FeedManagementAPIKey  = "FeedManagementAPI"
	IngestionAPIKey       = "IngestionAPI"
	GCTIAPIKey            = "GCTIAPI"
	RBACAPIKey            = "RBACAPI"
	ReferenceListsAPIKey  = "ReferenceListsAPI"
)

var RegionalSubDomains = map[string]map[string]string{
	SearchAPIKey: {
		RegionUS:             "backstory",
		RegionEurope:         "europe-backstory",
		RegionEuropeWest2:    "europe-west2-backstory",
		RegionAsiaSouthEast1: "asia-southeast1-backstory",
	},
	DetectionEngineAPIKey: {
		RegionUS:             "backstory",
		RegionEurope:         "europe-backstory",
		RegionEuropeWest2:    "europe-west2-backstory",
		RegionAsiaSouthEast1: "asia-southeast1-backstory",
	},
	FeedManagementAPIKey: {
		RegionUS:             "backstory",
		RegionEurope:         "europe-backstory",
		RegionEuropeWest2:    "europe-west2-backstory",
		RegionAsiaSouthEast1: "asia-southeast1-backstory",
	},
	IngestionAPIKey: {
		RegionUS:             "malachiteingestion-pa",
		RegionEurope:         "europe-malachiteingestion-pa",
		RegionEuropeWest2:    "europe-west2-malachiteingestion-pa",
		RegionAsiaSouthEast1: "asia-southeast1-malachiteingestion-pa",
	},
	GCTIAPIKey: {
		RegionUS:             "backstory",
		RegionEurope:         "backstory",
		RegionEuropeWest2:    "backstory",
		RegionAsiaSouthEast1: "backstory",
	},
	RBACAPIKey: {
		RegionUS:             "backstory",
		RegionEurope:         "europe-backstory",
		RegionEuropeWest2:    "europe-west2-backstory",
		RegionAsiaSouthEast1: "asia-southeast1-backstory",
	},
	ReferenceListsAPIKey: {
		RegionUS:             "backstory",
		RegionEurope:         "europe-backstory",
		RegionEuropeWest2:    "europe-west2-backstory",
		RegionAsiaSouthEast1: "asia-southeast1-backstory",
	},
}

const (
	EventsBasePathKey   = "Events"
	AlertBasePathKey    = "Alert"
	ArtifactBasePathKey = "Artifact"
	AliasBasePathKey    = "Alias"
	AssetBasePathKey    = "Asset"
	IOCBasePathKey      = "IOC"

	RuleBasePathKey           = "rules"
	FeedManagementBasePathKey = "Feed"

	SubjectsBasePathKey = "Subjects"

	ReferenceListsPathKey = "ReferenceLists"
)

func GenerateDefaultBasePaths(region string) map[string]string {
	var DefaultBasePaths = map[string]string{
		EventsBasePathKey:   getBasePathFromDomainsAndPath("/v1/events", RegionalSubDomains[SearchAPIKey][region]),
		AlertBasePathKey:    getBasePathFromDomainsAndPath("/v1/alert", RegionalSubDomains[SearchAPIKey][region]),
		ArtifactBasePathKey: getBasePathFromDomainsAndPath("/v1/artifact", RegionalSubDomains[SearchAPIKey][region]),
		AliasBasePathKey:    getBasePathFromDomainsAndPath("/v1/alias", RegionalSubDomains[SearchAPIKey][region]),
		AssetBasePathKey:    getBasePathFromDomainsAndPath("/v1/asset", RegionalSubDomains[SearchAPIKey][region]),
		IOCBasePathKey:      getBasePathFromDomainsAndPath("/v1/ioc", RegionalSubDomains[SearchAPIKey][region]),

		RuleBasePathKey:           getBasePathFromDomainsAndPath("/v2/detect/rules", RegionalSubDomains[SearchAPIKey][region]),
		FeedManagementBasePathKey: getBasePathFromDomainsAndPath("/v1/feeds", RegionalSubDomains[FeedManagementAPIKey][region]),

		SubjectsBasePathKey: getBasePathFromDomainsAndPath("/v1/subjects", RegionalSubDomains[RBACAPIKey][region]),

		ReferenceListsPathKey: getBasePathFromDomainsAndPath("/v2/lists", RegionalSubDomains[ReferenceListsAPIKey][region]),
	}

	return DefaultBasePaths
}

func getBasePathFromDomainsAndPath(basePath string, domain string) string {
	return fmt.Sprintf("https://%s.%s%s", domain, APIDomain, basePath)
}
