package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleFeedOktaUsers_Basic(t *testing.T) {
	displayName := "test" + randString(40)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	key := randString(10)
	value := randString(10)
	manager_id := "fooId"

	rootRef := feedOktaUsersRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaUsers(displayName, enabled, namespace, labels, hostname, manager_id, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaUsersExists(rootRef),
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

func TestAccChronicleFeedOktaUsers_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(40)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	manager_id := "fooId"
	key := "XXX"
	value := "XXX"
	key1 := randString(10)
	value1 := randString(10)

	rootRef := feedOktaUsersRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaUsers(displayName, enabled, namespace, labels, hostname, manager_id, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaUsersExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedOktaUsers(displayName, enabled, namespace, labels, hostname, manager_id, key1, value1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaUsersExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					testAccCheckChronicleFeedOktaUsersAuthUpdated(t, rootRef, key1, value1),
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

func TestAccChronicleFeedOktaUsers_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(40)
	enabled := "true"
	notEnabled := "false"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	manager_id := "fooId"
	key := randString(10)
	value := randString(10)

	rootRef := feedOktaUsersRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaUsers(displayName, enabled, namespace, labels, hostname, manager_id, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaUsersExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedOktaUsers(displayName, notEnabled, namespace, labels, hostname, manager_id, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaUsersExists(rootRef),
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

func TestAccChronicleOktaUsers_UpdateHostname(t *testing.T) {
	displayName := "test" + randString(40)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	hostname := randString(30) + ".okta.com"
	hostname1 := randString(30) + ".okta.com"
	manager_id := "fooId"
	key := randString(10)
	value := randString(10)

	rootRef := feedOktaUsersRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleFeedOktaUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleFeedOktaUsers(displayName, enabled, namespace, labels, hostname, manager_id, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaUsersExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					resource.TestCheckResourceAttr(rootRef, "details.0.hostname", hostname),
				),
			},
			{
				Config: testAccCheckChronicleFeedOktaUsers(displayName, enabled, namespace, labels, hostname1, manager_id, key, value),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleFeedOktaUsersExists(rootRef),
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
func testAccCheckChronicleFeedOktaUsersAuthUpdated(t *testing.T, n, key, value string) resource.TestCheckFunc {
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
func testAccCheckChronicleFeedOktaUsers(displayName, enabled, namespace, labels, hostname, managerID, key, value string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_okta_users" "test" {
			display_name = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				hostname = "%s"
				manager_id = "%s"
				authentication {
					key = "%s"
					value = "%s"
				}
			}
			}`, displayName, enabled, namespace, labels, hostname, managerID, key, value)
}

func testAccCheckChronicleFeedOktaUsersExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleFeedOktaUsersDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_okta_users.test" {
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
func feedOktaUsersRef(name string) string {
	return fmt.Sprintf("chronicle_feed_okta_users.%v", name)
}
