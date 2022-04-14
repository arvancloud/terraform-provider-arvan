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

func ResourceAbrakRescue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakRescueCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakRescueUpdate,
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
			"enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "enable rescue (false means un-rescue)",
			},
		},
	}
}

func resourceAbrakRescueCreate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	var errors diag.Diagnostics
	var err error

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

	enable := data.Get("enable").(bool)
	if enable {
		err = c.Server.Actions.Rescue(region, id)
	} else {
		err = c.Server.Actions.UnRescue(region, id)
	}

	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not rescue/un-rescue server %v", id),
		})
		return errors
	}

	data.SetId(id)
	return errors
}

func resourceAbrakRescueUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChange("enable") {
		return resourceAbrakRescueCreate(ctx, data, meta)
	}
	return nil
}
