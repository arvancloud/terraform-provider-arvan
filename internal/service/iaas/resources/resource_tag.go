package resources

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceTagUpdate,
		DeleteContext: resourceTagDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of tag",
			},
		},
	}
}

func resourceTagCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	var err error
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	// Tag options
	tagOpts := &iaas.TagOpts{
		TagName: data.Get("name").(string),
	}

	tag, err := c.Tag.Create(region, tagOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(tag.ID.(string))
	return errors
}

func resourceTagUpdate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	if data.HasChange("name") {
		// Tag options
		tagOpts := &iaas.TagUpdateOpts{
			TagName: data.Get("name").(string),
		}

		tag, err := c.Tag.Update(region, data.Id(), tagOpts)
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("can not delete FloatIP %v", data.Id()),
			})
			return errors
		}

		data.Set("name", tag.Name)

		return nil
	}
	return nil
}

func resourceTagDelete(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	// TODO: do we need to detach before of delete the tag ?

	err := c.Tag.Delete(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete FloatIP %v", data.Id()),
		})
		return errors
	}

	return nil
}
