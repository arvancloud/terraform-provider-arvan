package iaas

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

func ResourceAbrakRebuild() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakRebuildCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakRebuildUpdate,
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
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "uuid of abrak",
			},
			"image_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "uuid of image",
			},
		},
	}
}

func resourceAbrakRebuildCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var errors diag.Diagnostics
	c := meta.(*client.Client).Iaas

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	id := data.Get("uuid").(string)

	imageUuid := data.Get("image_uuid").(string)
	err := c.Server.Actions.Rebuild(region, id, imageUuid)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not rebuild server %v", id),
		})
		return errors
	}

	data.SetId(id)
	return errors
}

func resourceAbrakRebuildUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if data.HasChange("image_uuid") {
		return resourceAbrakRebuildCreate(ctx, data, meta)
	}
	return nil
}
