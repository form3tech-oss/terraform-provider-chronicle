package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedAzureBlobStore_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "GITHUB"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	uri := "https://myaccount.blob.core.windows.net/logging"
	sourceType := "FILES"
	sharedKey := "XXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAzureBlobStoreRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStore(displayName, logType, enabled, namespace, labels, uri, sourceType, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.shared_key", sharedKey),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.shared_key",
					"details.0.authentication.0.sas_token"},
			},
		},
	})
}

func TestAccChronicleFeedAzureBlobStore_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "GITHUB"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	uri := "https://myaccount.blob.core.windows.net/logging"
	sourceType := "FILES"
	sharedKey := "XXXXXXXXXXXXXXXXXXXX"
	sharedKey1 := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAzureBlobStoreRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStore(displayName, logType, enabled, namespace, labels, uri, sourceType, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.shared_key", sharedKey),
				),
			},
			{
				Config: testAccCheckChronicleFeedAzureBlobStore(displayName, logType, enabled, namespace, labels, uri, sourceType, sharedKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.shared_key", sharedKey1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.shared_key",
					"details.0.authentication.0.sas_token"},
			},
		},
	})
}

func TestAccChronicleFeedAzureBlobStore_UpdateEnabled(t *testing.T) {
	enabled := "true"
	notEnabled := "false"
	displayName := "test" + randString(10)
	logType := "GITHUB"
	namespace := "test"
	labels := `"test"="test"`
	uri := "https://myaccount.blob.core.windows.net/logging"
	sourceType := "FILES"
	sharedKey := "XXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAzureBlobStoreRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStore(displayName, logType, enabled, namespace, labels, uri, sourceType, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.shared_key", sharedKey),
				),
			},
			{
				Config: testAccCheckChronicleFeedAzureBlobStore(displayName, logType, notEnabled, namespace, labels, uri, sourceType, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.authentication.0.shared_key", sharedKey),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.shared_key",
					"details.0.authentication.0.sas_token"},
			},
		},
	})
}

func TestAccChronicleFeedAzureBlobStore_UpdateLogType(t *testing.T) {
	displayName := "test" + randString(10)
	notEnabled := "false"
	logType := "GITHUB"
	logType1 := "AWS_CLOUDTRAIL"
	namespace := "test"
	labels := `"test"="test"`
	uri := "https://myaccount.blob.core.windows.net/logging"
	sourceType := "FILES"
	sharedKey := "XXXXXXXXXXXXXXXXXXXX"

	rootRef := feedAzureBlobStoreRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStore(displayName, logType, notEnabled, namespace, labels, uri, sourceType, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAzureBlobStore(displayName, logType1, notEnabled, namespace, labels, uri, sourceType, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType1),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.shared_key",
					"details.0.authentication.0.sas_token"},
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedAzureBlobStore(displayName, logType, enabled, namespace, labels, uri, sourceType, shared_key string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_azure_blobstore" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				uri = "%s"
				source_type = "%s"
				authentication {
					shared_key = "%s"	
				}
			}
			}`, displayName, logType, enabled, namespace, labels, uri, sourceType, shared_key)
}

func testAccCheckChronicleFeedAzureBlobStoreExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedAzureBlobStoreDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_azure_blobstore.test" {
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
func feedAzureBlobStoreRef(name string) string {
	return fmt.Sprintf("chronicle_feed_azure_blobstore.%v", name)
}
