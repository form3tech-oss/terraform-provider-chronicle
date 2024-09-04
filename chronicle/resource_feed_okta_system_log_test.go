package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedOktaSystemLog_Basic(t *testing.T) {
	displayName := "test" + randString(40)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	key := randString(10)
	value := randString(10)

	rootRef := feedOktaSystemLogRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaSystemLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaSystemLog(displayName, enabled, namespace, labels, hostname, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaSystemLogExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.key",
					"details.0.authentication.0.value"},
			},
		},
	})
}

func TestAccChronicleFeedOktaSystemLog_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(40)
	displayName1 := "test" + randString(40)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	key := "XXX"
	value := "XXX"
	key1 := randString(10)
	value1 := randString(10)

	rootRef := feedOktaSystemLogRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaSystemLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaSystemLog(displayName, enabled, namespace, labels, hostname, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaSystemLogExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedOktaSystemLog(displayName1, enabled, namespace, labels, hostname, key1, value1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaSystemLogExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					testAccCheckChronicleFeedOktaSystemLogAuthUpdated(t, rootRef, key1, value1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.key",
					"details.0.authentication.0.value"},
			},
		},
	})
}

func TestAccChronicleFeedOktaSystemLog_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(40)
	displayName1 := "test" + randString(40)
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	key := randString(10)
	value := randString(10)

	rootRef := feedOktaSystemLogRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaSystemLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaSystemLog(displayName, enabled, namespace, labels, hostname, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaSystemLogExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedOktaSystemLog(displayName1, notEnabled, namespace, labels, hostname, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaSystemLogExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", notEnabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.key",
					"details.0.authentication.0.value"},
			},
		},
	})
}

func TestAccChronicleOktaSystemLog_UpdateHostname(t *testing.T) {
	displayName := "test" + randString(40)
	displayName1 := "test" + randString(40)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	hostname1 := randString(30) + ".okta.com"
	key := randString(10)
	value := randString(10)

	rootRef := feedOktaSystemLogRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaSystemLogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaSystemLog(displayName, enabled, namespace, labels, hostname, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaSystemLogExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedOktaSystemLog(displayName1, enabled, namespace, labels, hostname1, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaSystemLogExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"display_name", "state", "details.0.authentication.0.key",
					"details.0.authentication.0.value"},
			},
		},
	})
}

//nolint:unparam
func testAccCheckChronicleFeedOktaSystemLogAuthUpdated(t *testing.T, n, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		if rs.Primary.Attributes["details.0.authentication.0.key"] != key ||
			rs.Primary.Attributes["details.0.authentication.0.value"] != value {
			return fmt.Errorf("key or value differs")
		}

		return nil
	}
}

//nolint:unparam
func testAccCheckChronicleFeedOktaSystemLog(displayName, enabled, namespace, labels, hostname, key, value string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_okta_system_log" "test" {
			display_name = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				hostname = "%s"
				authentication {
					key = "%s"
					value = "%s"
				}
			}
			}`, displayName, enabled, namespace, labels, hostname, key, value)
}

func testAccCheckChronicleFeedOktaSystemLogExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedOktaSystemLogDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_okta_system_log.test" {
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
func feedOktaSystemLogRef(name string) string {
	return fmt.Sprintf("chronicle_feed_okta_system_log.%v", name)
}
