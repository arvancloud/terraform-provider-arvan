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
			"abrak_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "uuid of abrak",
				ValidateFunc: validation.IsUUID,
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "size of abrak",
			},
		},
	}
}

func resourceAbrakChangeDiskSizeCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	uuid := data.Get("abrak_uuid").(string)

	size := data.Get("size").(int)
	err := c.Server.Actions.ChangeDiskSize(region, uuid, size)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not change size of server %v to %v", uuid, size),
		})
		return errors
	}

	data.SetId(uuid)
	return errors
}

func resourceAbrakChangeDiskSizeUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChange("size") {
		return resourceAbrakChangeDiskSizeCreate(ctx, data, meta)
	}
	return nil
}
