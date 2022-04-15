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

func ResourcePtr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePtrCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: resourcePtrDelete,
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
			"ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ip address",
				ValidateFunc: validation.IsIPv4Address,
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "domain",
			},
		},
	}
}

func resourcePtrCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	// Ptr Options
	Ptr := &iaas.PtrOpts{
		IP:     data.Get("ip").(string),
		Domain: data.Get("domain").(string),
	}

	_, err := c.Ptr.Create(region, Ptr)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprint(Ptr.IP))
	return errors
}

func resourcePtrDelete(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	err := c.Ptr.Delete(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete Ptr %v", data.Id()),
		})
		return errors
	}

	return nil
}
