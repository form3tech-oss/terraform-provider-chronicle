package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedMicrosoftOffice365ManagementActivity_Basic(t *testing.T) {
	displayName := "testtf" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "manage.office.com/api/v1.0"
	tenantID := "50352502-a347-11ed-a8fc-0242ac120001"
	contentType := "AUDIT_AZURE_ACTIVE_DIRECTORY"
	clientID := "50352701-a307-11ed-a8fc-0242ac120001"
	clientSecret := "000"

	rootRef := feedMicrosoftOffice365ManagementActivityRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivity(displayName, enabled, namespace, labels, hostname, tenantID, contentType, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "details.0.content_type", contentType),
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

func TestAccChronicleFeedMicrosoftOffice365ManagementActivity_UpdateAuth(t *testing.T) {
	displayName := "testtf" + randString(10)
	displayName1 := "testf" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "manage.office.com/api/v1.0"
	tenantID := "50352702-a347-11ed-a8fc-0242ac120001"
	contentType := "AUDIT_AZURE_ACTIVE_DIRECTORY"
	clientID := "50352701-a307-11ed-a8fc-0242ac120001"
	clientID1 := "50352702-a307-11ed-a8fc-0242ac120001"
	clientSecret := "000"
	clientSecret1 := "001"

	rootRef := feedMicrosoftOffice365ManagementActivityRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivity(displayName, enabled, namespace, labels, hostname, tenantID, contentType, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "details.0.content_type", contentType),
				),
			},
			{
				Config: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivity(displayName1, enabled, namespace, labels, hostname, tenantID, contentType, clientID1, clientSecret1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "details.0.content_type", contentType),
					testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityAuthUpdated(t, rootRef, clientID1, clientSecret1),
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

func TestAccChronicleFeedMicrosoftOffice365ManagementActivity_UpdateEnabled(t *testing.T) {
	displayName := "testtf" + randString(10)
	displayName1 := "testf" + randString(10)
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "manage.office.com/api/v1.0"
	tenantID := "50352702-a347-11ed-a8fc-0242ac120001"
	contentType := "AUDIT_AZURE_ACTIVE_DIRECTORY"
	clientID := "50352701-a307-11ed-a8fc-0242ac120001"
	clientSecret := "000"

	rootRef := feedMicrosoftOffice365ManagementActivityRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivity(displayName, enabled, namespace, labels, hostname, tenantID, contentType, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "details.0.content_type", contentType),
				),
			},
			{
				Config: testAccCheckChronicleFeedMicrosoftOffice365ManagementActivity(displayName1, notEnabled, namespace, labels, hostname, tenantID, contentType, clientID, clientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					resource.TestCheckResourceAttr(rootRef, "details.0.tenant_id", tenantID),
					resource.TestCheckResourceAttr(rootRef, "details.0.content_type", contentType),
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
func testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityAuthUpdated(t *testing.T, n, clientID, clientSecret string) resource.TestCheckFunc {
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
func testAccCheckChronicleFeedMicrosoftOffice365ManagementActivity(displayName, enabled, namespace, labels, hostname, tenantID, contentType, clientID, clientSecret string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_microsoft_office_365_management_activity" "test" {
			display_name = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				hostname = "%s"
				tenant_id = "%s"
				content_type = "%s"
				authentication {
					client_id = "%s"
					client_secret = "%s"
				}
			}
			}`, displayName, enabled, namespace, labels, hostname, tenantID, contentType, clientID, clientSecret)
}

func testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedMicrosoftOffice365ManagementActivityDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_microsoft_office_365_management_activity.test" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func feedMicrosoftOffice365ManagementActivityRef(name string) string {
	return fmt.Sprintf("chronicle_feed_microsoft_office_365_management_activity.%v", name)
}
