package chronicle

import (
	"fmt"
	"log"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	RBACSubjectTypeAnalyst  = "SUBJECT_TYPE_ANALYST"
	RBACSubjectTypeIDPGroup = "SUBJECT_TYPE_IDP_GROUP"
)

func resourceRBACSubject() *schema.Resource {
	return &schema.Resource{
		Create: resourceRBACSubjectCreate,
		Read:   resourceRBACSubjectRead,
		Update: resourceRBACSubjectUpdate,
		Delete: resourceRBACSubjectDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(FiveMinutesTimeout),
			Update: schema.DefaultTimeout(FiveMinutesTimeout),
			Delete: schema.DefaultTimeout(FiveMinutesTimeout),
			Read:   schema.DefaultTimeout(FiveMinutesTimeout),
		},

		Description: "Creates a subject and assigns the given role.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the subject.`,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateSubjectType,
				Description:      `The type of the subject, e.g., SUBJECT_TYPE_ANALYST (an analyst) or SUBJECT_TYPE_IDP_GROUP (a group).`,
			},
			"roles": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The role(s) the created subject must have.`,
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
		},
	}
}

func resourceRBACSubjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	subject := expandSubject(d)

	log.Printf("[DEBUG] Creating new Schema: %#v", subject)
	err := client.CreateSubject(*subject)
	if err != nil {
		return fmt.Errorf("error creating Schema: %s", err)
	}

	d.SetId(subject.Name)

	log.Printf("[DEBUG] Finished creating Subject %q: %#v", d.Id(), subject)

	return resourceRBACSubjectRead(d, meta)
}

func resourceRBACSubjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	subject, err := client.GetSubject(d.Id())
	if err != nil {
		return HandleNotFoundError(err, d, d.Id())
	}

	if err := d.Set("name", subject.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set("type", subject.Type); err != nil {
		return fmt.Errorf("error reading Type: %s", err)
	}

	roleNames := flattenRoleNames(subject.Roles)

	if err := d.Set("roles", roleNames); err != nil {
		return fmt.Errorf("error reading Roles: %s", err)
	}

	log.Printf("[DEBUG] Finished reading Subject %q: %#v", d.Id(), subject)

	return nil
}

func resourceRBACSubjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	if d.HasChange("roles") {
		roles := readStringSliceFromResource(d, "roles")
		subject := chronicle.Subject{
			Name:  d.Id(),
			Type:  readStringFromResource(d, "type"),
			Roles: expandRolesFromRolesNameStringSlice(roles),
		}

		err := client.UpdateSubject(subject)
		if err != nil {
			return fmt.Errorf("error updating Subject %q: %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Finished updating Subject %q: %#v", d.Id(), subject)
	}

	return resourceRBACSubjectRead(d, meta)
}

func resourceRBACSubjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	log.Printf("[DEBUG] Deleting Schema: %#v", d.Id())
	err := client.DeleteSubject(d.Id())
	if err != nil {
		return handleNotFoundError(err, d, "Subject")
	}

	log.Printf("[DEBUG] Finished deleting Subject %q", d.Id())

	return nil
}

func expandSubject(d *schema.ResourceData) *chronicle.Subject {
	name := readStringFromResource(d, "name")
	subjectType := readStringFromResource(d, "type")

	rolesNames := readStringSliceFromResource(d, "roles")
	roles := expandRolesFromRolesNameStringSlice(rolesNames)

	return &chronicle.Subject{
		Name:  name,
		Type:  subjectType,
		Roles: roles,
	}
}

func expandRolesFromRolesNameStringSlice(rolesNames []string) []chronicle.Role {
	roles := make([]chronicle.Role, 0, len(rolesNames))
	for _, name := range rolesNames {
		roles = append(roles, chronicle.Role{
			Name: name,
		})
	}

	return roles
}

func flattenRoleNames(roles []chronicle.Role) []string {
	if roles == nil {
		return nil
	}
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	return roleNames
}
