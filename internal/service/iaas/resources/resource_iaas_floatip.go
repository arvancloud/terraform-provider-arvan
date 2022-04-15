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

func ResourceFloatIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFloatIPCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: resourceFloatIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "description",
			},
		},
	}
}

func resourceFloatIPCreate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	var errors diag.Diagnostics

	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	// FloatIP Options
	FloatIP := &iaas.FloatIPOpts{
		Description: data.Get("description").(string),
	}

	response, err := c.FloatIP.Create(region, FloatIP)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprint(response.ID))
	return errors
}

func resourceFloatIPDelete(_ context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	var errors diag.Diagnostics

	c := meta.(*client.Client).IaaS
	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	err := c.FloatIP.Delete(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete FloatIP %v", data.Id()),
		})
		return errors
	}

	return nil
}
