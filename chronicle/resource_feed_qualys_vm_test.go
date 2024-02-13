package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedQualysVM_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "qualysapi.qualys.com/api/2.0/fo/asset/host/?action=list"
	user := randString(10)
	secret := randString(10)

	rootRef := feedQualysVMRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedQualysVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedQualysVM(displayName, enabled, namespace, labels, hostname, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedQualysVMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.user",
					"details.0.authentication.0.secret"},
			},
		},
	})
}

func TestAccChronicleFeedQualysVM_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "qualysapi.qualys.com/api/2.0/fo/asset/host/?action=list"
	user := "XXXXXXXXXXXXXXXXXXXX"
	secret := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	user1 := randString(10)
	secret1 := randString(10)

	rootRef := feedQualysVMRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedQualysVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedQualysVM(displayName, enabled, namespace, labels, hostname, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedQualysVMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedQualysVM(displayName, enabled, namespace, labels, hostname, user1, secret1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedQualysVMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
					testAccCheckChronicleFeedQualysVMAuthUpdated(t, rootRef, user1, secret1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.user",
					"details.0.authentication.0.secret"},
			},
		},
	})
}

func TestAccChronicleFeedQualysVM_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "qualysapi.qualys.com/api/2.0/fo/asset/host/?action=list"
	user := randString(10)
	secret := randString(10)

	rootRef := feedQualysVMRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedQualysVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedQualysVM(displayName, enabled, namespace, labels, hostname, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedQualysVMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedQualysVM(displayName, notEnabled, namespace, labels, hostname, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedQualysVMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.user",
					"details.0.authentication.0.secret"},
			},
		},
	})
}

func TestAccChronicleFeedQualysVM_UpdateHostname(t *testing.T) {
	displayName := "test" + randString(10)
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	hostname := "qualysapi.qualys.com/api/2.0/fo/asset/host/?action=list"
	hostname1 := "qualysapi.qualys.eu/api/2.0/fo/asset/host/?action=list"
	user := randString(10)
	secret := randString(10)

	rootRef := feedQualysVMRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedQualysVMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedQualysVM(displayName, notEnabled, namespace, labels, hostname, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedQualysVMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedQualysVM(displayName, notEnabled, namespace, labels, hostname1, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedQualysVMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.user",
					"details.0.authentication.0.secret"},
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedQualysVMAuthUpdated(t *testing.T, n, user, secret string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["details.0.authentication.0.user"] != user ||
			rs.Primary.Attributes["details.0.authentication.0.secret"] != secret {
			return fmt.Errorf("user or secret differs")
		}

		return nil
	}
}

//nolint:unparam
func testAccCheckChronicleFeedQualysVM(displayName, enabled, namespace, labels, hostname, user, secret string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_qualys_vm" "test" {
			display_name = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				hostname = "%s"
				authentication {
					user = "%s"
					secret = "%s"
				}
			}
			}`, displayName, enabled, namespace, labels, hostname, user, secret)
}

func testAccCheckChronicleFeedQualysVMExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedQualysVMDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_qualys_vm.test" {
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
func feedQualysVMRef(name string) string {
	return fmt.Sprintf("chronicle_feed_qualys_vm.%v", name)
}
