package chronicle

import (
	"fmt"
	"log"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleCreate,
		Read:   resourceRuleRead,
		Update: resourceRuleUpdate,
		Delete: resourceRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(FiveMinutesTimeout),
			Update: schema.DefaultTimeout(FiveMinutesTimeout),
			Read:   schema.DefaultTimeout(FiveMinutesTimeout),
			Delete: schema.DefaultTimeout(FiveMinutesTimeout),
		},

		Description: `Creates a new rule and rule versions`,

		Schema: map[string]*schema.Schema{
			"rule_text": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      `Text of the new rule in YARA-L 2.0 format.`,
				ValidateDiagFunc: validateRuleText,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `Name of the rule as parsed from "rule_text".`,
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: false,
				Computed: true,
				Description: `Unique identifier for a specific version of a rule,
				 defined and returned by the server. You can specify exactly one version identifier.
				  Use the following format to specify the rule: {ruleId}@v_{int64}_{int64}.`,
			},
			"metadata": {
				Type:        schema.TypeMap,
				Elem:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `Metadata for the rule as parsed from "rule_text".`,
			},
			"rule_type": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `Type of rule. Can be either SINGLE_EVENT or MULTI_EVENT.`,
			},
			"live_enabled": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Computed:    true,
				Description: `Whether the rule is enabled to run as a Live Rule.`,
			},
			"alerting_enabled": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Computed:    true,
				Description: `Whether the rule is enabled to generate alerts.`,
			},
			"version_create_time": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `String representing the time in RFC 3339 format.`,
			},
			"compilation_state": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `Compilability of the rule ("SUCCEEDED" or "FAILED)".`,
			},
			"compilation_error": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `If "compilationState" is "FAILED", a compilation error for the rule is returned. Does not appear if "compilationState" is "SUCCEEDED".`,
			},
		},
	}
}

func resourceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	ruleRequest := chronicle.Rule{
		Text: readStringFromResource(d, "rule_text"),
	}

	if ok, err := client.VerifyYARARule(ruleRequest.Text); !ok {
		return fmt.Errorf("error verifying YARA-L 2.0 rule: %s", err)
	}

	log.Printf("[DEBUG] Creating new Schema: %#v", ruleRequest)

	id, err := client.CreateRule(ruleRequest)
	if err != nil {
		return fmt.Errorf("error creating Schema: %s", err)
	}

	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Rule %q", d.Id())

	if readBoolFromResource(d, "alerting_enabled") {
		log.Printf("[DEBUG] Enabling alerting for rule : %#v", id)

		err = client.ChangeAlertingRule(id, true)
		if err != nil {
			return fmt.Errorf("error enabling alerting: %s", err)
		}

		log.Printf("[DEBUG] Finished alerting for rule : %#v", id)
	}

	if readBoolFromResource(d, "live_enabled") {
		log.Printf("[DEBUG] Enabling live rule for rule : %#v", id)

		err = client.ChangeLiveRule(id, true)
		if err != nil {
			return fmt.Errorf("error enabling live rule: %s", err)
		}

		log.Printf("[DEBUG] Finished enabling live rule for rule : %#v", id)
	}

	return resourceRuleRead(d, meta)
}

func resourceRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	rule, err := client.GetRule(d.Id())
	if err != nil {
		return HandleNotFoundError(err, d, d.Id())
	}

	if err := d.Set("rule_text", rule.Text); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set("version_id", rule.VersionID); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set("rule_name", rule.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set("metadata", rule.Metadata); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set("rule_type", rule.Type); err != nil {
		return fmt.Errorf("error reading Type: %s", err)
	}
	if err := d.Set("live_enabled", rule.LiveEnabled); err != nil {
		return fmt.Errorf("error reading Roles: %s", err)
	}
	if err := d.Set("alerting_enabled", rule.AlertingEnabled); err != nil {
		return fmt.Errorf("error reading Roles: %s", err)
	}
	if err := d.Set("version_create_time", rule.VersionCreateTime); err != nil {
		return fmt.Errorf("error reading Roles: %s", err)
	}
	if err := d.Set("compilation_state", rule.CompilationState); err != nil {
		return fmt.Errorf("error reading Roles: %s", err)
	}
	if err := d.Set("compilation_error", rule.CompilationError); err != nil {
		return fmt.Errorf("error reading Roles: %s", err)
	}

	log.Printf("[DEBUG] Finished reading Rule %q: %#v", d.Id(), rule)

	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	// An update to rule_text creates a new version of the rule
	if d.HasChange("rule_text") {
		ruleText := readStringFromResource(d, "rule_text")
		ruleVersion := chronicle.Rule{
			ID:   d.Id(),
			Text: ruleText,
		}

		if ok, err := client.VerifyYARARule(ruleVersion.Text); !ok {
			return fmt.Errorf("error verifying YARA-L 2.0 rule: %s", err)
		}

		err := client.CreateRuleVersion(ruleVersion)
		if err != nil {
			return fmt.Errorf("error creating rule version for rule %q: %s", d.Id(), err)
		} else {
			log.Printf("[DEBUG] Finished creating new rule version for rule %q: %#v", d.Id(), ruleVersion)
		}
	}

	if d.HasChange("alerting_enabled") {
		log.Printf("[DEBUG] Enabling alerting for rule : %#v", d.Id())

		err := client.ChangeAlertingRule(d.Id(), readBoolFromResource(d, "alerting_enabled"))
		if err != nil {
			return fmt.Errorf("error enabling alerting: %s", err)
		} else {
			log.Printf("[DEBUG] Finished updating alerting for rule : %#v", d.Id())
		}
	}

	if d.HasChange("live_enabled") {
		log.Printf("[DEBUG] Enabling live rule for rule : %#v", d.Id())

		err := client.ChangeLiveRule(d.Id(), readBoolFromResource(d, "live_enabled"))
		if err != nil {
			return fmt.Errorf("error enabling live rule: %s", err)
		} else {
			log.Printf("[DEBUG] Finished enabling live rule for rule : %#v", d.Id())
		}
	}

	return resourceRuleRead(d, meta)
}

func resourceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	log.Printf("[DEBUG] Deleting Schema: %#v", d.Id())
	err := client.DeleteRule(d.Id())
	if err != nil {
		return handleNotFoundError(err, d, "Rule")
	}

	log.Printf("[DEBUG] Finished deleting Rule %q", d.Id())

	return nil
}
