package chronicle

import (
	"fmt"
	"testing"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccChronicleReferenceList_Basic(t *testing.T) {
	name := fmt.Sprintf("test%s", randString(5))
	description := "Acceptance test"
	lines := "Hello"
	contentType := string(chronicle.ReferenceListContentTypeREGEX)

	rootRef := referenceListRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleReferenceListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleReferenceList(name, description, contentType, lines),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleReferenceListExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "lines.0", lines),
					resource.TestCheckResourceAttr(rootRef, "description", description),
					resource.TestCheckResourceAttr(rootRef, "content_type", contentType),
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

func TestAccChronicleReferenceList_Update(t *testing.T) {
	name := fmt.Sprintf("test%s", randString(5))
	description := "Acceptance test"
	description1 := "Acceptance test1"
	lines := "Hello"
	lines1 := "Hello1"
	contentType := string(chronicle.ReferenceListContentTypeDefault)

	rootRef := referenceListRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleReferenceListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleReferenceList(name, description, contentType, lines),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleReferenceListExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "lines.0", lines),
					resource.TestCheckResourceAttr(rootRef, "description", description),
					resource.TestCheckResourceAttr(rootRef, "content_type", contentType),
				),
			},
			{
				Config: testAccCheckChronicleReferenceList(name, description1, contentType, lines1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleReferenceListExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "lines.0", lines1),
					resource.TestCheckResourceAttr(rootRef, "description", description1),
					resource.TestCheckResourceAttr(rootRef, "content_type", contentType),
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

func TestAccChronicleReferenceList_UpdateLines(t *testing.T) {
	name := fmt.Sprintf("test%s", randString(5))
	description := "Acceptance test"
	lines := "test"
	lines1 := "test1"
	contentType := string(chronicle.ReferenceListContentTypeDefault)

	rootRef := referenceListRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleReferenceListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleReferenceList(name, description, contentType, lines),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleReferenceListExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "lines.0", lines),
					resource.TestCheckResourceAttr(rootRef, "description", description),
					resource.TestCheckResourceAttr(rootRef, "content_type", contentType),
				),
			},
			{
				Config: testAccCheckChronicleReferenceList(name, description, contentType, lines1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleReferenceListExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "lines.0", lines1),
					resource.TestCheckResourceAttr(rootRef, "description", description),
					resource.TestCheckResourceAttr(rootRef, "content_type", contentType),
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

func TestAccChronicleReferenceList_UpdateDescription(t *testing.T) {
	name := fmt.Sprintf("test%s", randString(5))
	description := "Acceptance test"
	description1 := "Acceptance test 1"
	lines := "test"
	contentType := string(chronicle.ReferenceListContentTypeDefault)

	rootRef := referenceListRef("test")
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckChronicleReferenceListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckChronicleReferenceList(name, description, contentType, lines),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleReferenceListExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "lines.0", lines),
					resource.TestCheckResourceAttr(rootRef, "description", description),
					resource.TestCheckResourceAttr(rootRef, "content_type", contentType),
				),
			},
			{
				Config: testAccCheckChronicleReferenceList(name, description1, contentType, lines),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckChronicleReferenceListExists(rootRef),
					resource.TestCheckResourceAttr(rootRef, "name", name),
					resource.TestCheckResourceAttr(rootRef, "lines.0", lines),
					resource.TestCheckResourceAttr(rootRef, "description", description1),
					resource.TestCheckResourceAttr(rootRef, "content_type", contentType),
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

func testAccCheckChronicleReferenceList(name, description, contentType, lines string) string {
	s := fmt.Sprintf(
		`resource "chronicle_reference_list" "test" {
			name = %q
			description = "%s"
			content_type = "%s"
			lines = ["%s"]
		}`, name, description, contentType, lines)

	return s
}

func testAccCheckChronicleReferenceListExists(n string) resource.TestCheckFunc {
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

func testAccCheckChronicleReferenceListDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "chronicle_reference_list.test" {
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
func referenceListRef(name string) string {
	return fmt.Sprintf("chronicle_reference_list.%v", name)
}
