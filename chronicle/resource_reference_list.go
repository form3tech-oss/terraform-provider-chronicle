package chronicle

import (
	"fmt"
	"log"

	chronicle "github.com/form3tech-oss/terraform-provider-chronicle/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceReferenceList() *schema.Resource {
	return &schema.Resource{
		Create: resourceReferenceListCreate,
		Read:   resourceReferenceListRead,
		Update: resourceReferenceListUpdate,
		Delete: resourceReferenceListDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(FiveMinutesTimeout),
			Update: schema.DefaultTimeout(FiveMinutesTimeout),
			Read:   schema.DefaultTimeout(FiveMinutesTimeout),
			Delete: schema.DefaultTimeout(FiveMinutesTimeout),
		},

		Description: `Creates a reference list.`,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Unique name for the list.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: `Description of the list.`,
			},
			"content_type": {
				Type:             schema.TypeString,
				Required:         false,
				Optional:         true,
				ForceNew:         true,
				Default:          string(chronicle.ReferenceListContentTypeDefault),
				ValidateDiagFunc: validateReferenceListContentType,
				Description:      `Type of list content: "CONTENT_TYPE_DEFAULT_STRING", "REGEX", "CIDR". If omitted, defaults to "CONTENT_TYPE_DEFAULT_STRING".`,
			},
			"lines": {
				Type:        schema.TypeList,
				Required:    true,
				Optional:    false,
				Description: `List of line items.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"create_time": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Description: `Create time.`,
			},
		},
	}
}

func resourceReferenceListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	referenceListRequest := chronicle.ReferenceList{
		Name:        readStringFromResource(d, "name"),
		Description: readStringFromResource(d, "description"),
		ContentType: chronicle.ReferenceListContentType(readStringFromResource(d, "content_type")),
		Lines:       readStringSliceFromResource(d, "lines"),
	}

	id, err := client.CreateReferenceList(referenceListRequest)
	if err != nil {
		return fmt.Errorf("error creating Schema: %s", err)
	}

	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Reference List %q", d.Id())

	return resourceReferenceListRead(d, meta)
}

func resourceReferenceListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	referenceList, err := client.GetReferenceList(d.Id())
	if err != nil {
		return fmt.Errorf("error reading Schema: %s", err)
	}

	if err := d.Set("name", referenceList.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set("description", referenceList.Description); err != nil {
		return fmt.Errorf("error reading Description: %s", err)
	}
	if err := d.Set("content_type", referenceList.ContentType); err != nil {
		return fmt.Errorf("error reading ContentType: %s", err)
	}
	if err := d.Set("lines", referenceList.Lines); err != nil {
		return fmt.Errorf("error reading Lines: %s", err)
	}
	if err := d.Set("create_time", referenceList.CreateTime); err != nil {
		return fmt.Errorf("error reading create time: %s", err)
	}

	log.Printf("[DEBUG] Finished reading Reference List %q: %#v", d.Id(), referenceList)

	return nil
}

func resourceReferenceListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chronicle.Client)

	referenceList := chronicle.ReferenceList{
		Name:        readStringFromResource(d, "name"),
		Description: readStringFromResource(d, "description"),
		ContentType: chronicle.ReferenceListContentType(readStringFromResource(d, "content_type")),
		Lines:       readStringSliceFromResource(d, "lines"),
	}

	linesHasChange, descriptionHasChange := true, false
	if d.HasChange("description") {
		descriptionHasChange = true
	} else {
		referenceList.Description = ""
	}

	_, err := client.UpdateReferenceList(referenceList, linesHasChange, descriptionHasChange)
	if err != nil {
		return fmt.Errorf("error updating reference list: %s", err)
	}
	return resourceReferenceListRead(d, meta)
}

func resourceReferenceListDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting Schema: %#v", d.Id())
	// Delete method hasn't been implemented by Google yet.
	log.Printf("[DEBUG] Finished deleting Reference List %q", d.Id())

	return nil
}
