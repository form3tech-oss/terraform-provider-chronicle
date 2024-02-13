package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleRule_Basic(t *testing.T) {
	ruleText := `rule singleEventRule2{meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	liveDisabled := "false"
	alertingDisabled := "false"

	rootRef := rulePolicyRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRule(ruleText, liveDisabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRule_BasicTrailingNewLineRuleText(t *testing.T) {
	ruleText := `rule singleEventRule2{meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n" + "\n"
	liveDisabled := "false"
	alertingDisabled := "false"

	rootRef := rulePolicyRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRule(ruleText, liveDisabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRule_UpdateRuleText(t *testing.T) {
	ruleText := `rule singleEventRule2{    meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	ruleText1 := `rule singleEventRule2{    meta:      author = "newAuthor"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	liveDisabled := "false"
	alertingDisabled := "false"

	rootRef := rulePolicyRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRule(ruleText, liveDisabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				Config: testAccCheckChronicleRule(ruleText1, liveDisabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText1),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRule_UpdateAlerting(t *testing.T) {
	ruleText := `rule singleEventRule2{    meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	liveDisabled := "false"
	alertingEnabled := "true"
	alertingDisabled := "false"

	rootRef := rulePolicyRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRule(ruleText, liveDisabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				Config: testAccCheckChronicleRule(ruleText, liveDisabled, alertingEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingEnabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRule_UpdateLive(t *testing.T) {
	ruleText := `rule singleEventRule2{    meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	liveEnabled := "true"
	liveDisabled := "false"
	alertingDisabled := "false"

	rootRef := rulePolicyRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRule(ruleText, liveDisabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				Config: testAccCheckChronicleRule(ruleText, liveEnabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveEnabled),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRule_LiveAndAlertingEnabled(t *testing.T) {
	ruleText := `rule singleEventRule2{    meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	liveEnabled := "true"
	alertingEnabled := "true"

	rootRef := rulePolicyRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRule(ruleText, liveEnabled, alertingEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingEnabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveEnabled),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRule_UpdateRuleTextAndLiveAndAlerting(t *testing.T) {
	ruleText := `rule singleEventRule2{    meta:      author = "securityuser"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	ruleText1 := `rule singleEventRule2{    meta:      author = "newAuthor"      description = "single event rule that should generate detections TEST"
	    events:      $e.metadata.event_type = "NETWORK_DNS"    condition:       $e}` + "\n"
	liveDisabled := "false"
	alertingDisabled := "false"
	liveEnabled := "true"
	alertingEnabled := "true"

	rootRef := rulePolicyRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRule(ruleText, liveDisabled, alertingDisabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingDisabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveDisabled),
				),
			},
			{
				Config: testAccCheckChronicleRule(ruleText1, liveEnabled, alertingEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRuleExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "rule_text", ruleText1),
					resource.TestCheckResourceAttr(rootRef, "alerting_enabled", alertingEnabled),
					resource.TestCheckResourceAttr(rootRef, "live_enabled", liveEnabled),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckChronicleRule(ruleText string, liveEnabled string, alertingEnabled string) string {
	s := fmt.Sprintf(
		`resource "chronicle_rule" "test" {
			rule_text = %q
			alerting_enabled = "%s"
			live_enabled = "%s"
		}`, ruleText, alertingEnabled, liveEnabled)
	return s
}

func testAccCheckChronicleRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return NewNotFoundErrorf("%s in state", n)
		}

		if rs.Primary.ID == "" {
			return NewNotFoundErrorf("ID for %s in state", n)
		}
		return nil
	}
}

func testAccCheckChronicleRuleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_rule.test" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

//nolint:all
func rulePolicyRef(name string) string {
	return fmt.Sprintf("chronicle_rule.%v", name)
}
