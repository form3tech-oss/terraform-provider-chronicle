package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedGoogleCloudStorageV2_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"

	rootRef := feedGoogleCloudStorageV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageV2(displayName, logType, enabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
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

func TestAccChronicleFeedGoogleCloudStorageV2_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"

	rootRef := feedGoogleCloudStorageV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageV2(displayName, logType, enabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageV2(displayName1, logType, notEnabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageV2Exists(rootRef),
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

func TestAccChronicleFeedGoogleCloudStorageV2_UpdateLogType(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	logType1 := "GCP_DNS"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	sourceDeleteOptions := "ON_SUCCESS"
	maxLookbackDays := "90"

	rootRef := feedGoogleCloudStorageV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageV2(displayName, logType, enabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageV2(displayName1, logType1, enabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType1),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
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

func TestAccChronicleFeedGoogleCloudStorageV2_UpdateMaxLookbackDays(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GCP_CLOUDAUDIT"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "test-bucket/path"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	maxLookbackDays1 := "90"

	rootRef := feedGoogleCloudStorageV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageV2(displayName, logType, enabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageV2(displayName1, logType, enabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays1),
				),
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedGoogleCloudStorageV2(displayName, logType, enabled, namespace, labels, bucketUri,
	sourceDeleteOptions, maxLookbackDays string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_google_cloud_storage_v2" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				bucket_uri = "gs://%s/"
				source_delete_options = "%s"
				max_lookback_days = %s
			}
		}`, displayName, logType, enabled, namespace, labels, bucketUri, sourceDeleteOptions, maxLookbackDays)
}

func testAccCheckChronicleFeedGoogleCloudStorageV2Exists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedGoogleCloudStorageV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_google_cloud_storage_v2.test" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

//nolint:unparam
func feedGoogleCloudStorageV2Ref(name string) string {
	return fmt.Sprintf("chronicle_feed_google_cloud_storage_v2.%v", name)
}
