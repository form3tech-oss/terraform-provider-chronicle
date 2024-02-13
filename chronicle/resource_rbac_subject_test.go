package chronicle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleRBACSubject_Basic(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test%s", randString(5))
	subjectType := "SUBJECT_TYPE_ANALYST"
	role := "Editor"

	rootRef := rbacSubjectPolicyRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRBACSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRBACSubject(name, subjectType, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRBACSubjectExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "type", subjectType),
					resource.TestCheckResourceAttr(rootRef, "roles.0", role),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRBACSubject_UpdateRole(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test%s", randString(5))
	subjectType := "SUBJECT_TYPE_ANALYST"
	role := "Editor"
	role1 := "Editor"

	rootRef := rbacSubjectPolicyRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRBACSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRBACSubject(name, subjectType, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRBACSubjectExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "type", subjectType),
					resource.TestCheckResourceAttr(rootRef, "roles.0", role),
				),
			},
			{
				Config: testAccCheckChronicleRBACSubject(name, subjectType, role1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRBACSubjectExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "type", subjectType),
					resource.TestCheckResourceAttr(rootRef, "roles.0", role1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRBACSubject_UpdateType(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test%s", randString(5))
	subjectType := "SUBJECT_TYPE_ANALYST"
	subjectType1 := "SUBJECT_TYPE_IDP_GROUP"
	role := "Editor"

	rootRef := rbacSubjectPolicyRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRBACSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRBACSubject(name, subjectType, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRBACSubjectExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "type", subjectType),
					resource.TestCheckResourceAttr(rootRef, "roles.0", role),
				),
			},
			{
				Config: testAccCheckChronicleRBACSubject(name, subjectType1, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRBACSubjectExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "type", subjectType1),
					resource.TestCheckResourceAttr(rootRef, "roles.0", role),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccChronicleRBACSubject_UpdateTypeAndRole(t *testing.T) {
	t.Parallel()
	name := fmt.Sprintf("test%s", randString(5))
	subjectType := "SUBJECT_TYPE_ANALYST"
	subjectType1 := "SUBJECT_TYPE_IDP_GROUP"
	role := "Editor"
	role1 := "Viewer"

	rootRef := rbacSubjectPolicyRef("test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleRBACSubjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleRBACSubject(name, subjectType, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRBACSubjectExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "type", subjectType),
					resource.TestCheckResourceAttr(rootRef, "roles.0", role),
				),
			},
			{
				Config: testAccCheckChronicleRBACSubject(name, subjectType1, role1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleRBACSubjectExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "type", subjectType1),
					resource.TestCheckResourceAttr(rootRef, "roles.0", role1),
				),
			},
			{
				ResourceName:      rootRef,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckChronicleRBACSubject(name string, subjectType string, roles string) string {
	return fmt.Sprintf(
		`resource "chronicle_rbac_subject" "test" {
			name = "%s"
			type = "%s"
			roles = ["%s"]
		}`, name, subjectType, roles)
}

func testAccCheckChronicleRBACSubjectExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleRBACSubjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_rbac_subject.test" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Object %q still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}

//nolint:all
func rbacSubjectPolicyRef(name string) string {
	return fmt.Sprintf("chronicle_rbac_subject.%v", name)
}
