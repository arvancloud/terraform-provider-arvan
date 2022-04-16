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

func ResourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupRuleCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: resourceSecurityGroupRuleDelete,
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
			"security_group_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "uuid of security group",
				ValidateFunc: validation.IsUUID,
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "direction of rule",
				ValidateFunc: validation.StringInSlice(iaas.SupportedDirections, false),
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "protocol of rule",
			},
			"ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of ips, ['any'] for all ips",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Description:  "ip address",
					ValidateFunc: validateIPs,
				},
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "security group rule description",
			},
			"port_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "port from (If the rule is for one port, just fill this field with that port)",
			},
			"port_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "port to (If the rule is for all ports, leave port_from and port_to empty)",
			},
		},
	}
}

func resourceSecurityGroupRuleCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	securityGroupId := data.Get("security_group_uuid").(string)

	// get ips
	var ips []string
	if items, ok := data.GetOk("ips"); ok {
		ips = make([]string, len(items.([]any)))
		for i, ip := range items.([]any) {
			ips[i] = ip.(string)
		}
	} else {
		// default is accept all
		ips = append(ips, "any")
	}

	// securityGroupRule Options
	securityGroupRule := &iaas.SecurityGroupRuleOpts{
		Direction:   data.Get("direction").(string),
		PortFrom:    data.Get("port_from").(string),
		PortTo:      data.Get("port_to").(string),
		Protocol:    data.Get("protocol").(string),
		IPs:         ips,
		Description: data.Get("description").(string),
	}

	err := c.SecurityGroup.CreateRule(region, securityGroupId, securityGroupRule)
	if err != nil {
		return diag.FromErr(err)
	}

	rules, err := c.SecurityGroup.Read(region, securityGroupId)
	if err != nil {
		return diag.FromErr(err)
	}

	rule := helper.FindRule(rules.Rules, *securityGroupRule)
	if rule == nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not find new rule",
		})
		return errors
	}

	data.SetId(rule.ID)
	return errors
}

func resourceSecurityGroupRuleDelete(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	err := c.SecurityGroup.DeleteRule(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete SecurityGroup %v", data.Id()),
		})
		return errors
	}

	return nil
}

func validateIPs(i any, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("its not string"))
		return warnings, errors
	}
	if v == "any" {
		return warnings, errors
	}
	return validation.IsCIDR(i, k)
}
