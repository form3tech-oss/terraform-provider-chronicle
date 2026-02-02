package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedAmazonSQSV2_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	queue := "arn:aws:sqs:us-east-1:123456789012:test-queue"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSV2Ref()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(
					displayName, enabled, namespace, labels, s3Uri, queue,
					sourceDeleteOptions, maxLookbackDays, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.queue", queue),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options",
					"details.0.authentication.0.access_key_id", "details.0.authentication.0.secret_access_key"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQSV2_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	queue := "arn:aws:sqs:us-east-1:123456789012:test-queue"
	sourceDeleteOptions := "ON_SUCCESS"
	maxLookbackDays := "90"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	accessKeyID1 := "XXXXXXXXXXXXXXXXXXX1"
	secretAccessKey1 := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX1"

	rootRef := feedAmazonSQSV2Ref()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(
					displayName, enabled, namespace, labels, s3Uri, queue,
					sourceDeleteOptions, maxLookbackDays, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(
					displayName1, enabled, namespace, labels, s3Uri, queue,
					sourceDeleteOptions, maxLookbackDays, accessKeyID1, secretAccessKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options",
					"details.0.authentication.0.access_key_id", "details.0.authentication.0.secret_access_key"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQSV2_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	queue := "arn:aws:sqs:us-east-1:123456789012:test-queue"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSV2Ref()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(
					displayName, enabled, namespace, labels, s3Uri, queue,
					sourceDeleteOptions, maxLookbackDays, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(
					displayName1, notEnabled, namespace, labels, s3Uri, queue,
					sourceDeleteOptions, maxLookbackDays, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
				),
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQSV2_UpdateMaxLookbackDays(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	queue := "arn:aws:sqs:us-east-1:123456789012:test-queue"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	maxLookbackDays1 := "90"
	accessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSV2Ref()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(
					displayName, enabled, namespace, labels, s3Uri, queue,
					sourceDeleteOptions, maxLookbackDays, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(
					displayName1, enabled, namespace, labels, s3Uri, queue,
					sourceDeleteOptions, maxLookbackDays1, accessKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays1),
				),
			},
		},
	})
}

func testAccCheckChronicleFeedAmazonSQSV2(displayName, enabled, labels, s3Uri, queue,
	sourceDeleteOptions, maxLookbackDays, accessKeyID, secretAccessKey string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_amazon_sqs_v2" "test" {
			display_name = "%s"
			log_type = "AWS_CLOUDTRAIL"
			enabled = %s
			namespace = "test"
			labels = {
				%s
			}
			details {
				queue = "%s"
				s3_uri = "s3://%s/"
				source_delete_options = "%s"
				max_lookback_days = %s
				authentication {
					access_key_id = "%s"
					secret_access_key = "%s"
				}
			}
		}`, displayName, enabled, labels, queue, s3Uri,
		sourceDeleteOptions, maxLookbackDays, accessKeyID, secretAccessKey)
}

func testAccCheckChronicleFeedAmazonSQSV2Exists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedAmazonSQSV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_amazon_sqs_v2" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func feedAmazonSQSV2Ref() string {
	return "chronicle_feed_amazon_sqs_v2.test"
}
