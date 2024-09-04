package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedGoogleCloudStorageBucket_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "ONEPASSWORD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "gs://f3test/test/"
	bucketSourceType := FeedGoogleCloudStorageBucketSourceTypeFoldersRecursive
	sourceDeleteOptions := FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionNever

	rootRef := feedGoogleCloudStorageBucketRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName, logType, enabled, namespace, labels, bucketUri, bucketSourceType, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state"},
			},
		},
	})
}

func TestAccChronicleFeedGoogleCloudStorageBucket_UpdateBucketSourceType(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "ONEPASSWORD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "gs://f3test"
	bucketSourceType := FeedGoogleCloudStorageBucketSourceTypeFoldersRecursive
	bucketSourceType1 := FeedGoogleCloudStorageBucketSourceTypeFolders
	sourceDeleteOptions := FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionNever

	rootRef := feedGoogleCloudStorageBucketRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName, logType, enabled, namespace, labels, bucketUri, bucketSourceType, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName1, logType, enabled, namespace, labels, bucketUri, bucketSourceType1, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType1),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state"},
			},
		},
	})
}

func TestAccChronicleFeedGoogleCloudStorageBucket_UpdateSourceDeletionOptions(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "ONEPASSWORD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "gs://f3test"
	bucketSourceType := FeedGoogleCloudStorageBucketSourceTypeFoldersRecursive
	sourceDeleteOptions := FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionNever
	sourceDeleteOptions1 := FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionOnSuccess

	rootRef := feedGoogleCloudStorageBucketRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName, logType, enabled, namespace, labels, bucketUri, bucketSourceType, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName1, logType, enabled, namespace, labels, bucketUri, bucketSourceType, sourceDeleteOptions1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions1),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state"},
			},
		},
	})
}

func TestAccChronicleFeedGoogleCloudStorageBucket_UpdateLogType(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "ONEPASSWORD"
	logType1 := "AWS_CLOUDTRAIL"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	bucketUri := "gs://f3test"
	bucketSourceType := FeedGoogleCloudStorageBucketSourceTypeFoldersRecursive
	sourceDeleteOptions := FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionNever

	rootRef := feedGoogleCloudStorageBucketRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName, logType, enabled, namespace, labels, bucketUri, bucketSourceType, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName1, logType1, enabled, namespace, labels, bucketUri, bucketSourceType, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType1),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state"},
			},
		},
	})
}

func TestAccChronicleFeedGoogleCloudStorageBucket_UpdateNamespace(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "ONEPASSWORD"
	enabled := "true"
	namespace := "test"
	namespace1 := "test"
	labels := `"test"="test"`
	bucketUri := "gs://f3test"
	bucketSourceType := FeedGoogleCloudStorageBucketSourceTypeFoldersRecursive
	sourceDeleteOptions := FeedGoogleCloudStorageBucketSourceDeleteOptionDeletionNever

	rootRef := feedGoogleCloudStorageBucketRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedGoogleCloudStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName, logType, enabled, namespace, labels, bucketUri, bucketSourceType, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				Config: testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName1, logType, enabled, namespace1, labels, bucketUri, bucketSourceType, sourceDeleteOptions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedGoogleCloudStorageBucketExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace1),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_uri", bucketUri),
					resource.TestCheckResourceAttr(rootRef, "details.0.bucket_source_type", bucketSourceType),
					resource.TestCheckResourceAttr(rootRef, "details.0.source_delete_options", sourceDeleteOptions),
				),
			},
			{
				ResourceName:            rootRef,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "state"},
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedGoogleCloudStorageBucket(displayName, logType, enabled, namespace, labels, bucket_uri,
	bucket_source_type, sourceDeleteOptions string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_google_cloud_storage_bucket" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				bucket_uri = "%s"
				bucket_source_type = "%s"
				source_delete_options = "%s"
			}
		}`, displayName, logType, enabled, namespace, labels, bucket_uri,
		bucket_source_type, sourceDeleteOptions)
}

func testAccCheckChronicleFeedGoogleCloudStorageBucketExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedGoogleCloudStorageBucketDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_google_cloud_storage_bucket.test" {
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
func feedGoogleCloudStorageBucketRef(name string) string {
	return fmt.Sprintf("chronicle_feed_google_cloud_storage_bucket.%v", name)
}
