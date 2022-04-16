package resources

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubnetCreate,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetUpdate,
		DeleteContext: resourceSubnetDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of subnet",
			},
			"subnet_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "IP range of subnet",
				ValidateFunc: validation.IsCIDR,
			},
			"enable_gateway": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "enable gateway",
			},
			"gateway": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "gateway",
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "DHCP range(s)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "DHCP range beginning",
							ValidateFunc: validation.IsIPv4Address,
						},
						"to": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "DHCP range ending",
							ValidateFunc: validation.IsIPv4Address,
						},
					},
				},
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "dns servers",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Description:  "ip address of dns",
					ValidateFunc: validation.IsIPv4Address,
				},
			},
		},
	}
}

func resourceSubnetCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	// subnetOpts Options
	subnetOpts := &iaas.SubnetOpts{
		Name:          data.Get("name").(string),
		SubnetIP:      data.Get("subnet_ip").(string),
		EnableGateway: data.Get("enable_gateway").(bool),
		SubnetGateway: data.Get("gateway").(string),
	}

	// parse dhcp
	var dhcp string
	dhcpSet := data.Get("dhcp").(*schema.Set)
	for _, dhcpRange := range dhcpSet.List() {
		ranges, ok := dhcpRange.(map[string]any)
		if ok {
			dhcp += fmt.Sprintf("%v,%v\n", ranges["from"].(string), ranges["to"].(string))
		}
	}
	subnetOpts.Dhcp = dhcp

	// parse dns servers
	var dnsServers string
	items := data.Get("dns_servers").([]any)
	for _, item := range items {
		dnsServers += item.(string) + "\n"
	}
	subnetOpts.DnsServers = dnsServers

	response, err := c.Network.CreateSubnet(region, subnetOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprint(response.ID))
	return errors
}

func resourceSubnetRead(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	subnet, err := c.Network.ReadSubnet(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not retrieve Subnet %v", data.Id()),
		})
		return errors
	}

	data.SetId(subnet.ID)

	// TODO: we have to set other details

	return nil
}

func resourceSubnetUpdate(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	if data.HasChanges("name", "subnet_ip", "enable_gateway", "gateway", "dhcp", "dns_servers") {

		// subnetOpts Options
		subnetOpts := &iaas.SubnetOpts{
			Name:          data.Get("name").(string),
			SubnetIP:      data.Get("subnet_ip").(string),
			EnableGateway: data.Get("enable_gateway").(bool),
			SubnetGateway: data.Get("gateway").(string),
		}

		// parse dhcp
		var dhcp string
		dhcpSet := data.Get("dhcp").(*schema.Set)
		for _, dhcpRange := range dhcpSet.List() {
			ranges, ok := dhcpRange.(map[string]any)
			if ok {
				dhcp += fmt.Sprintf("%v,%v\n", ranges["from"].(string), ranges["to"].(string))
			}
		}
		subnetOpts.Dhcp = dhcp

		// parse dns servers
		var dnsServers string
		items := data.Get("dns_servers").([]any)
		for _, item := range items {
			dnsServers += item.(string) + "\n"
		}
		subnetOpts.DnsServers = dnsServers

		subnet, err := c.Network.UpdateSubnet(region, data.Id(), subnetOpts)
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("can not update subnet %v", data.Id()),
			})
			return errors
		}

		data.SetId(subnet.ID)
	}

	return errors
}

func resourceSubnetDelete(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	err := c.Network.DeleteSubnet(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete Subnet %v", data.Id()),
		})
		return errors
	}

	return nil
}
