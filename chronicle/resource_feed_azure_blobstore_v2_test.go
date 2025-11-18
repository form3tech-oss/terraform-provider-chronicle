package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedAzureBlobStoreV2_BasicWithSharedKey(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "AZURE_AD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	azureUri := "myaccount.blob.core.windows.net/logs"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	sharedKey := "dGVzdHNoYXJlZGtleXRlc3RzaGFyZWRrZXk="

	rootRef := feedAzureBlobStoreV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStoreV2WithSharedKey(displayName, logType, enabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreV2Exists(rootRef),
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
					"details.0.authentication.0.shared_key"},
			},
		},
	})
}

func TestAccChronicleFeedAzureBlobStoreV2_BasicWithSASToken(t *testing.T) {
	displayName := "test" + randString(10)
	logType := "AZURE_AD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	azureUri := "myaccount.blob.core.windows.net/logs"
	sourceDeleteOptions := "ON_SUCCESS"
	maxLookbackDays := "90"
	sasToken := "sv=2021-06-08&ss=bfqt&srt=sco&sp=rwdlacupiytfx&se=2025-01-01T00:00:00Z&st=2024-01-01T00:00:00Z&spr=https&sig=test"

	rootRef := feedAzureBlobStoreV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStoreV2WithSASToken(displayName, logType, enabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays, sasToken),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreV2Exists(rootRef),
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
					"details.0.authentication.0.sas_token"},
			},
		},
	})
}

func TestAccChronicleFeedAzureBlobStoreV2_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "AZURE_AD"
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	azureUri := "myaccount.blob.core.windows.net/logs"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	sharedKey := "dGVzdHNoYXJlZGtleXRlc3RzaGFyZWRrZXk="

	rootRef := feedAzureBlobStoreV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStoreV2WithSharedKey(displayName, logType, enabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "log_type", logType),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleFeedAzureBlobStoreV2WithSharedKey(displayName1, logType, notEnabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreV2Exists(rootRef),
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
					"details.0.authentication.0.shared_key"},
			},
		},
	})
}

func TestAccChronicleFeedAzureBlobStoreV2_UpdateMaxLookbackDays(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(10)
	logType := "AZURE_AD"
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	azureUri := "myaccount.blob.core.windows.net/logs"
	sourceDeleteOptions := "NEVER"
	maxLookbackDays := "180"
	maxLookbackDays1 := "90"
	sharedKey := "dGVzdHNoYXJlZGtleXRlc3RzaGFyZWRrZXk="

	rootRef := feedAzureBlobStoreV2Ref("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedAzureBlobStoreV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedAzureBlobStoreV2WithSharedKey(displayName, logType, enabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays),
				),
			},
			{
				Config: testAccCheckChronicleFeedAzureBlobStoreV2WithSharedKey(displayName1, logType, enabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays1, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedAzureBlobStoreV2Exists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "details.0.max_lookback_days", maxLookbackDays1),
				),
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedAzureBlobStoreV2WithSharedKey(displayName, logType, enabled, namespace, labels, azureUri,
	sourceDeleteOptions, maxLookbackDays, sharedKey string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_azure_blobstore_v2" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				azure_uri = "https://%s/"
				source_delete_options = "%s"
				max_lookback_days = %s
				authentication {
					shared_key = "%s"
				}
			}
		}`, displayName, logType, enabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays, sharedKey)
}

//nolint:unparam
func testAccCheckChronicleFeedAzureBlobStoreV2WithSASToken(displayName, logType, enabled, namespace, labels, azureUri,
	sourceDeleteOptions, maxLookbackDays, sasToken string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_azure_blobstore_v2" "test" {
			display_name = "%s"
			log_type = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				azure_uri = "https://%s/"
				source_delete_options = "%s"
				max_lookback_days = %s
				authentication {
					sas_token = "%s"
				}
			}
		}`, displayName, logType, enabled, namespace, labels, azureUri, sourceDeleteOptions, maxLookbackDays, sasToken)
}

func testAccCheckChronicleFeedAzureBlobStoreV2Exists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedAzureBlobStoreV2Destroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_azure_blobstore_v2.test" {
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
func feedAzureBlobStoreV2Ref(name string) string {
	return fmt.Sprintf("chronicle_feed_azure_blobstore_v2.%v", name)
}
