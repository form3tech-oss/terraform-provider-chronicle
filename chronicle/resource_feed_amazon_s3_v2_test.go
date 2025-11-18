package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedAmazonS3V2_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	region := "us-east-1"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonS3V2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3V2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3V2(displayName, logType, enabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays, region, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3V2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options",
					"details.0.authentication.0.access_key_id", "details.0.authentication.0.secret_access_key",
					"details.0.authentication.0.region"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonS3V2_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	sourceDeleteOptions := "ON_SUCCESS"
	maxLookbackDays := "90"
	region := "us-east-1"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	region1 := "us-west-2"
	accessKeyID1 := "XXXXXXXXXXXXXXXXXXX1"
	secretAccessKey1 := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX1"

	rootRef := feedAmazonS3V2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3V2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3V2(displayName, logType, enabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays, region, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3V2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonS3V2(displayName1, logType, enabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays, region1, accessKeyID1, secretAccessKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3V2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					testAccCheckChronicleFeedAmazonS3V2AuthUpdated(t, rootRef, region1, accessKeyID1, secretAccessKey1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options",
					"details.0.authentication.0.access_key_id", "details.0.authentication.0.secret_access_key",
					"details.0.authentication.0.region"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonS3V2_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	region := "us-east-1"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonS3V2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3V2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3V2(displayName, logType, enabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays, region, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3V2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonS3V2(displayName1, logType, notEnabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays, region, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3V2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options",
					"details.0.authentication.0.access_key_id", "details.0.authentication.0.secret_access_key",
					"details.0.authentication.0.region"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonS3V2_UpdateMaxLookbackDays(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	maxLookbackDays1 := "90"
	region := "us-east-1"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonS3V2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3V2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3V2(displayName, logType, enabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays, region, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3V2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonS3V2(displayName1, logType, enabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays1, region, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3V2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays1),
				),
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedAmazonS3V2AuthUpdated(t *testing.T, n, region, accessKeyID, secretAccessKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["details.0.authentication.0.access_key_id"] != accessKeyID ||
			rs.Primary.Attributes["details.0.authentication.0.secret_access_key"] != secretAccessKey ||
			rs.Primary.Attributes["details.0.authentication.0.region"] != region {
			return fmt.Errorf("accessKeyID or secretAccessKey or region differs")
		}

		return nil
	}
}

//nolint:unparam
func testAccCheckChronicleFeedAmazonS3V2(displayName, logType, enabled, namespace, labels, s3Uri,
	sourceDeleteOptions, maxLookbackDays, region, accessKeyID, secretAccessKey string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_amazon_s3_v2" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				s3_uri = "s3://%s/"
				source_delete_options = "%s"
				max_lookback_days = %s
				authentication {
					region = "%s"
					access_key_id = "%s"
					secret_access_key = "%s"
				}
			}
		}`, displayName, logType, enabled, namespace, labels, s3Uri, sourceDeleteOptions, maxLookbackDays, region, accessKeyID, secretAccessKey)
}

func testAccCheckChronicleFeedAmazonS3V2Exists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedAmazonS3V2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_amazon_s3_v2.test" {
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
func feedAmazonS3V2Ref(name string) string {
	return fmt.Sprintf("chronicle_feed_amazon_s3_v2.%v", name)
}
