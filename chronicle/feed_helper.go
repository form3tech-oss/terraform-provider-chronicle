package chronicle

import (
	"fmt"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Feed State as returned from Chronicle Feed API.
const (
	FeedStateActive     = "ACTIVE"
	FeedStateInactive   = "INACTIVE"
	FeedStateInProgress = "IN_PROGRESS"
	FeedStateCompleted  = "COMPLETED"
	FeedStateFailed     = "FAILED"
)

// A ConcreteFeedFlattenFunc flatten a ConcreteFeedConfiguration using an original and a read value.
type ConcreteFeedFlattenFunc func(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{}

// A ConcreteFeedExpandFunc retunrs a ConcreteFeedConfiguration from ResourceData attributes.
type ConcreteFeedExpandFunc func(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration

// ConcreteFeedResource interface contains methods to be implemented by specific feed resources such as S3.
type ConcreteFeedResource interface {
	// flattenDetailsFromReadOperation is the implementation of ConcreteFeedFlattenFunc.
	flattenDetailsFromReadOperation(originalConf chronicle.ConcreteFeedConfiguration, readConf chronicle.ConcreteFeedConfiguration) []map[string]interface{}
	// expandConcreteFeedConfiguration is the implementation of ConcreteFeedExpandFunc.
	expandConcreteFeedConfiguration(d *schema.ResourceData) chronicle.ConcreteFeedConfiguration
	// getLogType returns feed's log type if it's static otherwise an empty string
	getLogType() string
}

func newFeedResourceSchema(details *schema.Resource, concreteFeed ConcreteFeedResource, description string, withLogType bool) *schema.Resource {
	resource := &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			return resourceFeedCreate(d, meta, concreteFeed.expandConcreteFeedConfiguration, concreteFeed.flattenDetailsFromReadOperation, concreteFeed.getLogType())
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return resourceFeedRead(d, meta, concreteFeed.expandConcreteFeedConfiguration, concreteFeed.flattenDetailsFromReadOperation)
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			return resourceFeedUpdate(d, meta, concreteFeed.expandConcreteFeedConfiguration, concreteFeed.flattenDetailsFromReadOperation)
		},
		Delete: resourceFeedDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(FiveMinutesTimeout),
			Update: schema.DefaultTimeout(FiveMinutesTimeout),
			Delete: schema.DefaultTimeout(FiveMinutesTimeout),
			Read:   schema.DefaultTimeout(FiveMinutesTimeout),
		},

		Description: description,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Name to be displayed.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Enabled specifies whether a feed is allowed to be executed.`,
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `State gives some insight into the current state of a feed.`,
			},

			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Log Type is a label which describes the nature of the data being ingested.`,
			},
			"feed_source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Feed Source Type describes how data is collected.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The namespace the feed will be associated with.`,
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `All of the events that result from this feed will have this label applied.`,
			},

			"details": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Each feed type has its own requirements for which this field must fulfil.`,
				Elem:        details,
			},
		},
	}

	if !withLogType {
		resource.Schema["log_type"].Required = false
		resource.Schema["log_type"].Computed = true
	}

	return resource
}

func resourceFeedRead(d *schema.ResourceData, meta interface{}, expandFunc ConcreteFeedExpandFunc, flattenDetailsFromConcreteConfiguration ConcreteFeedFlattenFunc) error {
	client := meta.(*chronicle.Client)

	baseFeed, concreteFeed, err := client.ReadFeed(d.Id())
	if err != nil {
		return HandleNotFoundError(err, d, d.Id())
	}

	err = setBaseFeedProperties(d, *baseFeed)
	if err != nil {
		return err
	}

	details := flattenDetailsFromConcreteConfiguration(expandFunc(d), *concreteFeed)
	if err := d.Set("details", details); err != nil {
		return fmt.Errorf("error setting Details: %s", err)
	}

	return nil
}

func resourceFeedCreate(d *schema.ResourceData, meta interface{}, expandFunc ConcreteFeedExpandFunc,
	flattenDetailsFromConcreteConfiguration ConcreteFeedFlattenFunc, staticLogType string) error {
	client := meta.(*chronicle.Client)

	concreteFeed := expandFunc(d)

	var logType string
	if staticLogType != "" {
		logType = staticLogType
	} else {
		logType = readStringFromResource(d, "log_type")
	}

	id, err := client.CreateFeed(readStringFromResource(d, "display_name"), logType,
		readStringFromResource(d, "namespace"), extractLabelsFromFeedResource(d), concreteFeed)
	if err != nil {
		return err
	}

	d.SetId(id)

	if !readBoolFromResource(d, "enabled") {
		err = client.ChangeEnableFeed(id, false)
		if err != nil {
			return err
		}
	}

	return resourceFeedRead(d, meta, expandFunc, flattenDetailsFromConcreteConfiguration)
}

func resourceFeedUpdate(d *schema.ResourceData, meta interface{}, expandFunc ConcreteFeedExpandFunc, flattenDetailsFromConcreteConfiguration ConcreteFeedFlattenFunc) error {
	client := meta.(*chronicle.Client)

	concreteFeed := expandFunc(d)

	if d.HasChange("details") || d.HasChange("display_name") || d.HasChange("log_type") ||
		d.HasChange("namespace") || d.HasChange("labels") {
		err := client.UpdateFeed(d.Id(), readStringFromResource(d, "display_name"), readStringFromResource(d, "log_type"), readStringFromResource(d, "namespace"),
			extractLabelsFromFeedResource(d), concreteFeed)
		if err != nil {
			return err
		}
	}

	if d.HasChange("enabled") {
		err := client.ChangeEnableFeed(d.Id(), readBoolFromResource(d, "enabled"))
		if err != nil {
			return err
		}
	}
	return resourceFeedRead(d, meta, expandFunc, flattenDetailsFromConcreteConfiguration)
}

func resourceFeedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	if err := client.DestroyFeed(d.Id()); err != nil {
		return err
	}

	return nil
}

func setBaseFeedProperties(d *schema.ResourceData, feed chronicle.BaseFeed) error {
	if err := d.Set("state", feed.State); err != nil {
		return fmt.Errorf("error reading State: %s", err)
	}
	if err := d.Set("feed_source_type", feed.Details.SourceType); err != nil {
		return fmt.Errorf("error reading feed source type: %s", err)
	}
	if err := d.Set("log_type", feed.Details.LogType); err != nil {
		return fmt.Errorf("error reading log type: %s", err)
	}
	if err := d.Set("namespace", feed.Details.Namespace); err != nil {
		return fmt.Errorf("error reading log type: %s", err)
	}

	if err := d.Set("labels", extractLabelMapFromFeedLabels(feed.Details.Labels)); err != nil {
		return fmt.Errorf("error reading log type: %s", err)
	}

	if err := d.Set("enabled", d.Get("state").(string) != FeedStateInactive); err != nil {
		return fmt.Errorf("error setting enabled: %s", err)
	}

	return nil
}

func extractLabelsFromFeedResource(d *schema.ResourceData) []chronicle.Label {
	labelsRaw := readMapFromResource(d, "labels")

	if labelsRaw == nil {
		return nil
	}
	labels := make([]chronicle.Label, 0, len(labelsRaw))
	for k, v := range labelsRaw {
		labels = append(labels, chronicle.Label{Key: k, Value: v.(string)})
	}

	return labels
}

func extractLabelMapFromFeedLabels(labels []chronicle.Label) map[string]string {
	if labels == nil {
		return nil
	}

	labelMap := map[string]string{}

	for _, label := range labels {
		labelMap[label.Key] = label.Value
	}

	return labelMap
}
