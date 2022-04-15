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

func ResourceAbrakReboot() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceAbrakRebootCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: ResourceAbrakRebootUpdate,
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
			"hard_reboot": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "hard reboot",
			},
		},
	}
}

func ResourceAbrakRebootCreate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	var errors diag.Diagnostics
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

	id := data.Get("uuid").(string)

	hardReboot := data.Get("hard_reboot").(bool)
	if hardReboot {
		err = c.Server.Actions.HardReboot(region, id)
	} else {
		err = c.Server.Actions.SoftReboot(region, id)
	}

	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not turn on server %v", id),
		})
		return errors
	}

	data.SetId(id)
	return errors
}

func ResourceAbrakRebootUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChange("hard_reboot") {
		return ResourceAbrakRebootCreate(ctx, data, meta)
	}
	return nil
}
