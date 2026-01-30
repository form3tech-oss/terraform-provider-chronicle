package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedGoogleCloudStorageEventDriven_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	pubsubSubscription := "projects/test-project/subscriptions/test-subscription"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"

	rootRef := feedGoogleCloudStorageEventDrivenRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageEventDrivenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName, logType, enabled, namespace, labels, bucketUri, pubsubSubscription, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.pubsub_subscription", pubsubSubscription),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options"},
			},
		},
	})
}

func TestAccChronicleFeedGoogleCloudStorageEventDriven_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	pubsubSubscription := "projects/test-project/subscriptions/test-subscription"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"

	rootRef := feedGoogleCloudStorageEventDrivenRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageEventDrivenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName, logType, enabled, namespace, labels, bucketUri, pubsubSubscription, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName1, logType, notEnabled, namespace, labels, bucketUri, pubsubSubscription, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options"},
			},
		},
	})
}

func TestAccChronicleFeedGoogleCloudStorageEventDriven_UpdateSubscription(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	pubsubSubscription := "projects/test-project/subscriptions/test-subscription"
	pubsubSubscription1 := "projects/test-project/subscriptions/test-subscription-updated"
	sourceDeleteOptions := "ON_SUCCESS"
	maxLookbackDays := "90"

	rootRef := feedGoogleCloudStorageEventDrivenRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageEventDrivenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName, logType, enabled, namespace, labels, bucketUri, pubsubSubscription, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.pubsub_subscription", pubsubSubscription),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName1, logType, enabled, namespace, labels, bucketUri, pubsubSubscription1, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.pubsub_subscription", pubsubSubscription1),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options"},
			},
		},
	})
}

func TestAccChronicleFeedGoogleCloudStorageEventDriven_UpdateMaxLookbackDays(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	pubsubSubscription := "projects/test-project/subscriptions/test-subscription"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	maxLookbackDays1 := "90"

	rootRef := feedGoogleCloudStorageEventDrivenRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageEventDrivenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName, logType, enabled, namespace, labels, bucketUri, pubsubSubscription, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName1, logType, enabled, namespace, labels, bucketUri, pubsubSubscription, sourceDeleteOptions, maxLookbackDays1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays1),
				),
			},
		},
	})
}

func testAccCheckChronicleFeedGoogleCloudStorageEventDriven(displayName, logType, enabled, namespace, labels, bucketUri,
	pubsubSubscription, sourceDeleteOptions, maxLookbackDays string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_google_cloud_storage_event_driven" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				bucket_uri = "gs://%s/"
				pubsub_subscription = "%s"
				source_delete_options = "%s"
				max_lookback_days = %s
			}
		}`, displayName, logType, enabled, namespace, labels, bucketUri, pubsubSubscription, sourceDeleteOptions, maxLookbackDays)
}

func testAccCheckChronicleFeedGoogleCloudStorageEventDrivenExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedGoogleCloudStorageEventDrivenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_google_cloud_storage_event_driven.test" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func feedGoogleCloudStorageEventDrivenRef(name string) string {
	return fmt.Sprintf("chronicle_feed_google_cloud_storage_event_driven.%v", name)
}
