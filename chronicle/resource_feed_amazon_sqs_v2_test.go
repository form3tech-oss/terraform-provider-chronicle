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
	region := "us-east-1"
	accountNumber := "123456789012"
	queueName := "test-queue"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	sqsAccessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(displayName, logType, enabled, namespace, labels, s3Uri, region, accountNumber, queueName, sourceDeleteOptions, maxLookbackDays, sqsAccessKeyID, sqsSecretAccessKey, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
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
					"details.0.authentication.0.sqs_access_key_id", "details.0.authentication.0.sqs_secret_access_key"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQSV2_WithS3Auth(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	region := "us-east-1"
	accountNumber := "123456789012"
	queueName := "test-queue"
	sourceDeleteOptions := "ON_SUCCESS"
	maxLookbackDays := "90"
	sqsAccessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	s3AccessKeyID := "XXXXXXXXXXXXXXXXXXX1"
	s3SecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX1"

	rootRef := feedAmazonSQSV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(displayName, logType, enabled, namespace, labels, s3Uri, region, accountNumber, queueName, sourceDeleteOptions, maxLookbackDays, sqsAccessKeyID, sqsSecretAccessKey, s3AccessKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.s3_authentication.0.access_key_id", s3AccessKeyID),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options",
					"details.0.authentication.0.sqs_access_key_id", "details.0.authentication.0.sqs_secret_access_key",
					"details.0.authentication.0.s3_authentication.0.access_key_id", "details.0.authentication.0.s3_authentication.0.secret_access_key"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQSV2_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	region := "us-east-1"
	accountNumber := "123456789012"
	queueName := "test-queue"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	sqsAccessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(displayName, logType, enabled, namespace, labels, s3Uri, region, accountNumber, queueName, sourceDeleteOptions, maxLookbackDays, sqsAccessKeyID, sqsSecretAccessKey, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(displayName1, logType, notEnabled, namespace, labels, s3Uri, region, accountNumber, queueName, sourceDeleteOptions, maxLookbackDays, sqsAccessKeyID, sqsSecretAccessKey, "", ""),
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
	logType := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test-bucket"
	region := "us-east-1"
	accountNumber := "123456789012"
	queueName := "test-queue"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	maxLookbackDays1 := "90"
	sqsAccessKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(displayName, logType, enabled, namespace, labels, s3Uri, region, accountNumber, queueName, sourceDeleteOptions, maxLookbackDays, sqsAccessKeyID, sqsSecretAccessKey, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSV2(displayName1, logType, enabled, namespace, labels, s3Uri, region, accountNumber, queueName, sourceDeleteOptions, maxLookbackDays1, sqsAccessKeyID, sqsSecretAccessKey, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays1),
				),
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedAmazonSQSV2(displayName, logType, enabled, namespace, labels, s3Uri, region, accountNumber, queueName,
	sourceDeleteOptions, maxLookbackDays, sqsAccessKeyID, sqsSecretAccessKey, s3AccessKeyID, s3SecretAccessKey string) string {

	s3AuthBlock := ""
	if s3AccessKeyID != "" && s3SecretAccessKey != "" {
		s3AuthBlock = fmt.Sprintf(`
				s3_authentication {
					access_key_id = "%s"
					secret_access_key = "%s"
				}`, s3AccessKeyID, s3SecretAccessKey)
	}

	return fmt.Sprintf(
		`resource "chronicle_feed_amazon_sqs_v2" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				s3_uri = "s3://%s/"
				region = "%s"
				account_number = "%s"
				queue_name = "%s"
				source_delete_options = "%s"
				max_lookback_days = %s
				authentication {
					sqs_access_key_id = "%s"
					sqs_secret_access_key = "%s"
					%s
				}
			}
		}`, displayName, logType, enabled, namespace, labels, s3Uri, region, accountNumber, queueName,
		sourceDeleteOptions, maxLookbackDays, sqsAccessKeyID, sqsSecretAccessKey, s3AuthBlock)
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
		if rs.Type != "chronicle_feed_amazon_sqs_v2.test" {
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
func feedAmazonSQSV2Ref(name string) string {
	return fmt.Sprintf("chronicle_feed_amazon_sqs_v2.%v", name)
}
