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

func ResourceAbrakChangeDiskSize() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakChangeDiskSizeCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakChangeDiskSizeUpdate,
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
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "size of abrak",
			},
		},
	}
}

func resourceAbrakChangeDiskSizeCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	size := data.Get("size").(int)
	err := c.Server.Actions.ChangeDiskSize(region, id, size)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not rename server %v to %v", id, size),
		})
		return errors
	}

	data.SetId(id)
	return errors
}

func resourceAbrakChangeDiskSizeUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if data.HasChange("size") {
		return resourceAbrakChangeDiskSizeCreate(ctx, data, meta)
	}
	return nil
}
