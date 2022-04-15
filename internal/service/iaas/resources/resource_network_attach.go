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

func ResourceNetworkAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkAttachCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: resourceNetworkAttachDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of network",
			},
			"abrak_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of abrak",
				ValidateFunc: validation.IsUUID,
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ip address in this network which we want to assign to abrak",
				ValidateFunc: validation.IsIPv4Address,
			},
			"enable_port_security": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "enable port security",
			},
		},
	}
}

func resourceNetworkAttachCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	name := data.Get("name").(string)

	// looking for network id
	network, err := c.Network.Find(region, name)
	if err != nil {
		return diag.FromErr(err)
	}

	// networkAttach Options
	networkAttachOpts := &iaas.NetworkAttachOpts{
		ServerId: data.Get("abrak_id").(string),
	}

	err = c.Network.Attach(region, network.ID, networkAttachOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(network.ID)
	return errors
}

func resourceNetworkAttachDelete(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	// networkDetach Options
	networkDetachOpts := &iaas.NetworkDetachOpts{
		ServerId: data.Get("abrak_id").(string),
	}

	err := c.Network.Detach(region, data.Id(), networkDetachOpts)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not detach network %v", data.Id()),
		})
		return errors
	}

	return nil
}
