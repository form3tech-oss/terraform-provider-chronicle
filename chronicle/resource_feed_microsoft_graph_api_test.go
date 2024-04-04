package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedMicrosoftGraphAPI_Basic(t *testing.T) {
	displayName := "testtf" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "graph.microsoft.com/v1.0/auditLogs/directoryAudits"
	tenantID := "50352502-a347-11ed-a8fc-0242ac120001"
	logType := "AZURE_AD_AUDIT"
	clientID := "50352701-a307-11ed-a8fc-0242ac120001"
	clientSecret := "000"

	rootRef := feedMicrosoftGraphAPIRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedMicrosoftGraphAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedMicrosoftGraphAPI(displayName, enabled, namespace, labels, hostname, logType, tenantID, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftGraphAPIExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.hostname", "details.0.authentication.0.client_id",
					"details.0.authentication.0.client_secret"},
			},
		},
	})
}

func TestAccChronicleFeedMicrosoftGraphAPI_UpdateAuth(t *testing.T) {
	displayName := "testtf" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "graph.microsoft.com/v1.0/auditLogs/signIns"
	tenantID := "50352702-a347-11ed-a8fc-0242ac120001"
	logType := "AZURE_AD"
	clientID := "50352701-a307-11ed-a8fc-0242ac120001"
	clientID1 := "50352702-a307-11ed-a8fc-0242ac120001"
	clientSecret := "000"
	clientSecret1 := "001"

	rootRef := feedMicrosoftGraphAPIRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedMicrosoftGraphAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedMicrosoftGraphAPI(displayName, enabled, namespace, labels, logType, hostname, tenantID, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftGraphAPIExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
				),
			},
			{
				Config: testAccCheckChronicleFeedMicrosoftGraphAPI(displayName, enabled, namespace, labels, logType, hostname, tenantID, clientID1, clientSecret1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftGraphAPIExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					testAccCheckChronicleFeedMicrosoftGraphAPIAuthUpdated(t, rootRef, clientID1, clientSecret1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.hostname", "details.0.authentication.0.client_id",
					"details.0.authentication.0.client_secret"},
			},
		},
	})
}

func TestAccChronicleFeedMicrosoftGraphAPI_UpdateEnabled(t *testing.T) {
	displayName := "testtf" + randString(10)
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "graph.microsoft.com/beta"
	tenantID := "50352702-a347-11ed-a8fc-0242ac120001"
	logType := "AZURE_AD_CONTEXT"
	clientID := "50352701-a307-11ed-a8fc-0242ac120001"
	clientSecret := "000"

	rootRef := feedMicrosoftGraphAPIRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedMicrosoftGraphAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedMicrosoftGraphAPI(displayName, enabled, namespace, labels, logType, hostname, tenantID, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftGraphAPIExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
				),
			},
			{
				Config: testAccCheckChronicleFeedMicrosoftGraphAPI(displayName, notEnabled, namespace, labels, logType, hostname, tenantID, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftGraphAPIExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.hostname", "details.0.authentication.0.client_id",
					"details.0.authentication.0.client_secret"},
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedMicrosoftGraphAPIAuthUpdated(t *testing.T, n, clientID, clientSecret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["details.0.authentication.0.client_id"] != clientID ||
			rs.Primary.Attributes["details.0.authentication.0.client_secret"] != clientSecret {
			return fmt.Errorf("clientID or clientSecret differs")
		}

		return nil
	}
}

//nolint:unparam
func testAccCheckChronicleFeedMicrosoftGraphAPI(displayName, enabled, namespace, labels, logType, hostname, tenantID, clientID, clientSecret string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_microsoft_graph_api" "test" {
			display_name = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			log_type = "%s"
			details {
				hostname = "%s"
				tenant_id = "%s"
				authentication {
					client_id = "%s"
					client_secret = "%s"
				}
			}
			}`, displayName, enabled, namespace, labels, logType, hostname, tenantID, clientID, clientSecret)
}

func testAccCheckChronicleFeedMicrosoftGraphAPIExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedMicrosoftGraphAPIDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_microsoft_graph_api.test" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func feedMicrosoftGraphAPIRef(name string) string {
	return fmt.Sprintf("chronicle_feed_microsoft_graph_api.%v", name)
}
