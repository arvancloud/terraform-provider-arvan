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

func ResourceNetworkDetach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkDetachCreate,
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
				Description:  "Region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"network_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of network",
				ValidateFunc: validation.IsUUID,
			},
			"abrak_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of abrak",
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceNetworkDetachCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	networkId := data.Get("network_id").(string)

	// networkAttach Options
	networkAttachOpts := &iaas.NetworkDetachOpts{
		ServerId: data.Get("abrak_id").(string),
	}

	err := c.Network.Detach(region, networkId, networkAttachOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(networkId)
	return errors
}
