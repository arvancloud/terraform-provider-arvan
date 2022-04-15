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

func ResourceAbrakChangeFlavor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakChangeFlavorCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakChangeFlavorUpdate,
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
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "new flavor",
			},
		},
	}
}

func resourceAbrakChangeFlavorCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
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

	flavor := data.Get("flavor").(string)
	err := c.Server.Actions.ChangeFlavor(region, uuid, flavor)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not change flavor of server %v to %v", uuid, flavor),
		})
		return errors
	}

	data.SetId(uuid)
	return errors
}

func resourceAbrakChangeFlavorUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChange("flavor") {
		return resourceAbrakChangeFlavorCreate(ctx, data, meta)
	}
	return nil
}
