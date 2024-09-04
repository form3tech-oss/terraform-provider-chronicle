package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedAmazonS3_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "GITHUB"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test"
	s3SourceType := "FILES"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	region := "EU_WEST_1"
	accesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonS3Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3(displayName, logType, enabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region, accesKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.access_key_id",
					"details.0.authentication.0.secret_access_key", "details.0.authentication.0.region"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonS3_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test"
	s3SourceType := "FILES"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	region := "EU_WEST_1"
	accesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	region1 := "EU_WEST_2"
	accesKeyID1 := "XXXXXXXXXXXXXXXXXXX1"
	secretAccessKey1 := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX1"

	rootRef := feedAmazonS3Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3(displayName, logType, enabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region, accesKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonS3(displayName1, logType, enabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region1, accesKeyID1, secretAccessKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					testAccCheckChronicleFeedAmazonS3AuthUpdated(t, rootRef, region1, accesKeyID1, secretAccessKey1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.access_key_id",
					"details.0.authentication.0.secret_access_key", "details.0.authentication.0.region"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonS3_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test"
	s3SourceType := "FILES"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	region := "EU_WEST_1"
	accesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonS3Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3(displayName, logType, enabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region, accesKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonS3(displayName1, logType, notEnabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region, accesKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.access_key_id",
					"details.0.authentication.0.secret_access_key", "details.0.authentication.0.region"},
			},
		},
	})
}

func TestAccChronicleFeedAmazonS3_UpdateLogType(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "GITHUB"
	logType1 := "AWS_CLOUDTRAIL"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	s3Uri := "test"
	s3SourceType := "FILES"
	sourceDeleteOptions := "SOURCE_DELETION_NEVER"
	region := "EU_WEST_1"
	accesKeyID := "XXXXXXXXXXXXXXXXXXXX"
	secretAccessKey := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAmazonS3Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAmazonS3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAmazonS3(displayName, logType, notEnabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region, accesKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAmazonS3(displayName1, logType1, notEnabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region, accesKeyID, secretAccessKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAmazonS3Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType1),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.source_delete_options", "details.0.authentication.0.access_key_id",
					"details.0.authentication.0.secret_access_key", "details.0.authentication.0.region"},
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedAmazonS3AuthUpdated(t *testing.T, n, region, accessKeyID, secretAccessKey string) resource.TestCheckFunc {
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
			return fmt.Errorf("accessKeyID or secretAccessKey differs")
		}

		return nil
	}
}

//nolint:unparam
func testAccCheckChronicleFeedAmazonS3(displayName, logType, enabled, namespace, labels, s3Uri, s3SourceType,
	sourceDeleteOptions, region, accesKeyID, secretAccessKey string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_amazon_s3" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				s3_uri = "s3://%s/"
				s3_source_type = "%s"
				source_delete_options = "%s"
				authentication {
					region = "%s"
					access_key_id = "%s"	
					secret_access_key = "%s"
				}
			}
			}`, displayName, logType, enabled, namespace, labels, s3Uri, s3SourceType, sourceDeleteOptions, region, accesKeyID, secretAccessKey)
}

func testAccCheckChronicleFeedAmazonS3Exists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedAmazonS3Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_amazon_s3.test" {
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
func feedAmazonS3Ref(name string) string {
	return fmt.Sprintf("chronicle_feed_amazon_s3.%v", name)
}
