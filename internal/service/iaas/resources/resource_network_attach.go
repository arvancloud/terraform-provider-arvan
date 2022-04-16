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

	networkId := data.Get("network_id").(string)

	// networkAttach Options
	networkAttachOpts := &iaas.NetworkAttachOpts{
		ServerId: data.Get("abrak_id").(string),
	}

	if ip, ok := data.GetOk("ip"); ok {
		networkAttachOpts.IP = ip.(string)
	}

	if enablePortSecurity, ok := data.GetOk("enable_port_security"); ok {
		networkAttachOpts.EnablePortSecurity = enablePortSecurity.(bool)
	}

	network, err := c.Network.Attach(region, networkId, networkAttachOpts)
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
