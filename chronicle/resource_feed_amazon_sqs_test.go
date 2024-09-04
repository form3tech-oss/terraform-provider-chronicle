package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedAmazonSQS_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "ONEPASSWORD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	queue := "test"
	region := "US_EAST_1"
	accountNumber := "111111111111"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	sqsAccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX1"

	rootRef := feedAmazonSQSRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQS(displayName, logType, enabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.region", region),
					resource.TestCheckResourceAttr(rootRef, "details.0.account_number", accountNumber),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.sqs_access_key_id", sqsAccesKeyID),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.sqs_secret_access_key", sqsSecretAccessKey),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.sqs_access_key_id",
					"details.0.authentication.0.sqs_secret_access_key", "details.0.authentication.0.s3_authentication"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQS_BasicWithS3Auth(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "ONEPASSWORD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	queue := "test"
	region := "US_EAST_1"
	accountNumber := "111111111111"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	sqsAccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX1"
	s3AccesKeyID := "XXXXXXXXXXXXXXXXXXX1"
	s3SecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX1"

	rootRef := feedAmazonSQSRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName, logType, enabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.region", region),
					resource.TestCheckResourceAttr(rootRef, "details.0.account_number", accountNumber),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.sqs_access_key_id", sqsAccesKeyID),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.sqs_secret_access_key", sqsSecretAccessKey),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.s3_authentication.0.access_key_id", s3AccesKeyID),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.s3_authentication.0.secret_access_key", s3SecretAccessKey),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.sqs_access_key_id",
					"details.0.authentication.0.sqs_secret_access_key", "details.0.authentication.0.s3_authentication"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQS_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	queue := "test"
	region := "EU_WEST_1"
	region1 := "EU_WEST_2"
	accountNumber := "111111111111"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	sqsAccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	s3AccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	s3SecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	sqsAccesKeyID1 := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey1 := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	s3AccesKeyID1 := "XXXXXXXXXXXXXXXXXXXX"
	s3SecretAccessKey1 := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName, logType, enabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName1, logType, enabled, namespace, labels, queue, region1, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID1, sqsSecretAccessKey1, s3AccesKeyID1, s3SecretAccessKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					testAccCheckChronicleFeedAmazonSQSAuthUpdated(t, rootRef, region1, sqsAccesKeyID1, sqsSecretAccessKey1, s3AccesKeyID1, s3SecretAccessKey1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.sqs_access_key_id",
					"details.0.authentication.0.sqs_secret_access_key", "details.0.authentication.0.s3_authentication"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQS_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	queue := "test"
	region := "EU_WEST_1"
	accountNumber := "111111111111"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	sqsAccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	s3AccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	s3SecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName, logType, enabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName1, logType, notEnabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.sqs_access_key_id",
					"details.0.authentication.0.sqs_secret_access_key", "details.0.authentication.0.s3_authentication"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQS_UpdateLogType(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	logType1 := "AWS_CLOUDTRAIL"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	queue := "test"
	region := "EU_WEST_1"
	accountNumber := "111111111111"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	sqsAccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	s3AccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	s3SecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName, logType, notEnabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName1, logType1, notEnabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType1),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.sqs_access_key_id",
					"details.0.authentication.0.sqs_secret_access_key", "details.0.authentication.0.s3_authentication"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQS_UpdateAccountNumber(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	queue := "test"
	region := "EU_WEST_1"
	accountNumber := "111111111111"
	accountNumber1 := "111111111112"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	sqsAccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	s3AccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	s3SecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName, logType, notEnabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.account_number", accountNumber),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName1, logType, notEnabled, namespace, labels, queue, region, accountNumber1,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.account_number", accountNumber1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.sqs_access_key_id",
					"details.0.authentication.0.sqs_secret_access_key", "details.0.authentication.0.s3_authentication"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonSQS_UpdateRegion(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	queue := "test"
	region := "EU_WEST_1"
	region1 := "EU_WEST_2"
	accountNumber := "111111111111"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	sqsAccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	sqsSecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	s3AccesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	s3SecretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonSQSRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonSQSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName, logType, notEnabled, namespace, labels, queue, region, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.region", region),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName1, logType, notEnabled, namespace, labels, queue, region1, accountNumber,
					sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonSQSExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.region", region1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.sqs_access_key_id",
					"details.0.authentication.0.sqs_secret_access_key", "details.0.authentication.0.s3_authentication"},
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedAmazonSQSAuthUpdated(t *testing.T, n, region, sqsAccessKeyID,
	sqsSecretAccessKey, s3AccessKeyID, s3SecretAccessKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["details.0.authentication.0.sqs_access_key_id"] != sqsAccessKeyID ||
			rs.Primary.Attributes["details.0.authentication.0.sqs_secret_access_key"] != sqsSecretAccessKey ||
			rs.Primary.Attributes["details.0.authentication.0.s3_authentication.0.access_key_id"] != s3AccessKeyID ||
			rs.Primary.Attributes["details.0.authentication.0.s3_authentication.0.secret_access_key"] != s3SecretAccessKey ||
			rs.Primary.Attributes["details.0.region"] != region {
			return fmt.Errorf("accessKeyID, secretAccessKey, or region differs")
		}

		return nil
	}
}

func testAccCheckChronicleFeedAmazonSQS(displayName, logType, enabled, namespace, labels, queue, region,
	accountNumber, sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_amazon_sqs" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				queue = "%s"
				region = "%s"
				account_number = "%s"
				source_delete_options = "%s"
				authentication {
					sqs_access_key_id = "%s"	
					sqs_secret_access_key = "%s"
				}
			}
		}`, displayName, logType, enabled, namespace, labels, queue, region,
		accountNumber, sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey)
}

//nolint:unparam
func testAccCheckChronicleFeedAmazonSQSWithS3Auth(displayName, logType, enabled, namespace, labels, queue, region,
	accountNumber, sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_amazon_sqs" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				queue = "%s"
				region = "%s"
				account_number = "%s"
				source_delete_options = "%s"
				authentication {
					sqs_access_key_id = "%s"	
					sqs_secret_access_key = "%s"
					s3_authentication {
						access_key_id = "%s"
						secret_access_key = "%s"
					}	
				}
			}
		}`, displayName, logType, enabled, namespace, labels, queue, region,
		accountNumber, sourceDeleteOptions, sqsAccesKeyID, sqsSecretAccessKey, s3AccesKeyID, s3SecretAccessKey)
}

func testAccCheckChronicleFeedAmazonSQSExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedAmazonSQSDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_amazon_sqs.test" {
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
func feedAmazonSQSRef(name string) string {
	return fmt.Sprintf("chronicle_feed_amazon_sqs.%v", name)
}
