package resources

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
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
				Optional:    true,
				Default:     false,
				Description: "enable gateway",
			},
			"enable_dhcp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "enable dhcp",
			},
			"gateway": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "gateway",
				ValidateFunc: validation.IsIPv4Address,
			},
			"dhcp": {
				Type:        schema.TypeSet,
				Optional:    true,
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
				Optional:    true,
				Description: "dns servers",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Description:  "ip address of dns",
					ValidateFunc: validation.IsIPv4Address,
				},
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of network",
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
		Name:     data.Get("name").(string),
		SubnetIP: data.Get("subnet_ip").(string),
	}

	if enableGateway, ok := data.GetOk("enable_gateway"); ok {
		subnetOpts.EnableGateway = enableGateway.(bool)
	}

	if enableDhcp, ok := data.GetOk("enable_dhcp"); ok {
		subnetOpts.EnableDhcp = enableDhcp.(bool)
	}

	if subnetGateway, ok := data.GetOk("gateway"); ok {
		subnetOpts.SubnetGateway = subnetGateway.(string)
	}

	// parse dhcp
	if dhcpSet, ok := data.GetOk("dhcp"); ok {
		var dhcp string
		for _, dhcpRange := range dhcpSet.(*schema.Set).List() {
			ranges, ok := dhcpRange.(map[string]any)
			if ok {
				dhcp += fmt.Sprintf("%v,%v\n", ranges["from"].(string), ranges["to"].(string))
			}
		}
		subnetOpts.Dhcp = strings.TrimSuffix(dhcp, "\n")
	}

	// parse dns servers
	if items, ok := data.GetOk("dns_servers"); ok {
		var dnsServers string
		for _, item := range items.([]any) {
			dnsServers += item.(string) + "\n"
		}
		subnetOpts.DnsServers = strings.TrimSuffix(dnsServers, "\n")
	}

	response, err := c.Network.CreateSubnet(region, subnetOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ID)
	data.Set("network_uuid", response.NetworkId)
	data.Set("gateway", response.GatewayIP)
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
	data.Set("network_uuid", subnet.NetworkId)
	data.Set("gateway", subnet.GatewayIP)
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

	if data.HasChanges("name", "subnet_ip", "enable_gateway", "gateway", "enable_dhcp", "dhcp", "dns_servers") {

		// subnetOpts Options
		subnetOpts := &iaas.SubnetOpts{
			Name:      data.Get("name").(string),
			SubnetIP:  data.Get("subnet_ip").(string),
			SubnetId:  data.Id(),
			NetworkId: data.Get("network_uuid").(string),
		}

		if enableGateway, ok := data.GetOk("enable_gateway"); ok {
			subnetOpts.EnableGateway = enableGateway.(bool)
		}

		if enableDhcp, ok := data.GetOk("enable_dhcp"); ok {
			subnetOpts.EnableDhcp = enableDhcp.(bool)
		}

		if subnetGateway, ok := data.GetOk("gateway"); ok {
			subnetOpts.SubnetGateway = subnetGateway.(string)
		}

		// parse dhcp
		if dhcpSet, ok := data.GetOk("dhcp"); ok {
			var dhcp string
			for _, dhcpRange := range dhcpSet.(*schema.Set).List() {
				ranges, ok := dhcpRange.(map[string]any)
				if ok {
					dhcp += fmt.Sprintf("%v,%v\n", ranges["from"].(string), ranges["to"].(string))
				}
			}
			subnetOpts.Dhcp = strings.TrimSuffix(dhcp, "\n")
		}

		// parse dns servers
		if items, ok := data.GetOk("dns_servers"); ok {
			var dnsServers string
			for _, item := range items.([]any) {
				dnsServers += item.(string) + "\n"
			}
			subnetOpts.DnsServers = strings.TrimSuffix(dnsServers, "\n")
		}

		err := c.Network.UpdateSubnet(region, subnetOpts)
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("can not update subnet %v", data.Id()),
			})
			return errors
		}

		data.SetId(data.Id())
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

	// TODO: do we need to detach the attached servers ?

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
