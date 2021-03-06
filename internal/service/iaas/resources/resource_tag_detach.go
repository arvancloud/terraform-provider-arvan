package resources

import (
	"context"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTagDetach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagDetachCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: helper.DummyResourceAction,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"tag_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of tag",
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of instance",
				ValidateFunc: validation.IsUUID,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "type of instance",
				ValidateFunc: validation.StringInSlice(iaas.SupportedTagTypes, false),
			},
		},
	}
}

func resourceTagDetachCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	tagId := data.Get("tag_id").(string)

	// Tag options
	tagOpts := &iaas.TagAttachmentOpts{
		InstanceId:   data.Get("instance_id").(string),
		InstanceType: data.Get("instance_type").(string),
	}

	err := c.Tag.Detach(region, tagId, tagOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(tagId)
	return errors
}
