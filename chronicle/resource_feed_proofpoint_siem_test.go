package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleProofpointSIEM_Basic(t *testing.T) {
	displayName := "test" + randString(10)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	user := randString(10)
	secret := randString(10)

	rootRef := ProofpointSIEMRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleProofpointSIEMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleProofpointSIEM(displayName, enabled, namespace, labels, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleProofpointSIEMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
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

func TestAccChronicleProofpointSIEM_UpdateAuth(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(40)
	enabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	user := "XXXXXXXXXXXXXXXXXXXX"
	secret := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	user1 := randString(10)
	secret1 := randString(10)

	rootRef := ProofpointSIEMRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleProofpointSIEMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleProofpointSIEM(displayName, enabled, namespace, labels, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleProofpointSIEMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleProofpointSIEM(displayName1, enabled, namespace, labels, user1, secret1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleProofpointSIEMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
					testAccCheckChronicleProofpointSIEMAuthUpdated(t, rootRef, user1, secret1),
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

func TestAccChronicleProofpointSIEM_UpdateEnabled(t *testing.T) {
	displayName := "test" + randString(10)
	displayName1 := "test" + randString(40)
	enabled := "true"
	notEnabled := "true"
	namespace := "test"
	labels := `"test"="test"`
	user := randString(10)
	secret := randString(10)

	rootRef := ProofpointSIEMRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleProofpointSIEMDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleProofpointSIEM(displayName, enabled, namespace, labels, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleProofpointSIEMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
				),
			},
			{
				Config: testAccCheckChronicleProofpointSIEM(displayName1, notEnabled, namespace, labels, user, secret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleProofpointSIEMExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "enabled", enabled),
					resource.TestCheckResourceAttr(rootRef, "namespace", namespace),
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
func testAccCheckChronicleProofpointSIEMAuthUpdated(t *testing.T, n, user, secret string) resource.TestCheckFunc {
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
func testAccCheckChronicleProofpointSIEM(displayName, enabled, namespace, labels, user, secret string) string {
	return fmt.Sprintf(
		`resource "chronicle_feed_proofpoint_siem" "test" {
			display_name = "%s"
			enabled = %s
			namespace = "%s"
			labels = {
				%s
			}
			details {
				authentication {
					user = "%s"
					secret = "%s"
				}
			}
			}`, displayName, enabled, namespace, labels, user, secret)
}

func testAccCheckChronicleProofpointSIEMExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleProofpointSIEMDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_feed_proofpoint_siem.test" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

func ProofpointSIEMRef(name string) string {
	return fmt.Sprintf("chronicle_feed_proofpoint_siem.%v", name)
}
