package datasources

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

func DatasourceQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceQuotaRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"max_image_meta": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum image meta",
			},
			"max_personality": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum personality",
			},
			"max_security_group_rules": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum security group rules",
			},
			"max_security_groups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum security groups",
			},
			"max_server_group_members": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum server group members",
			},
			"max_server_groups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum server groups",
			},
			"max_server_meta": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum server meta",
			},
			"max_total_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum CPU cores",
			},
			"max_total_floating_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum floating ips",
			},
			"max_total_instances": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum instances",
			},
			"max_total_keypairs": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum key-pairs",
			},
			"max_total_ram_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "maximum RAM size",
			},
			"total_cores_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total CPU cores used",
			},
			"total_floating_ip_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total floating ips used",
			},
			"total_instances_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total instances used",
			},
			"total_ram_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total RAM used",
			},
			"total_security_groups_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total security groups used",
			},
			"total_server_groups_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "total server groups used",
			},
		},
	}
}

func datasourceQuotaRead(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	quota, err := c.Quota.Read(region)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("quota not found"),
		})
		return errors
	}

	data.SetId(helper.GenUUID())
	data.Set("max_image_meta", quota.MaxImageMeta)
	data.Set("max_personality", quota.MaxPersonality)
	data.Set("max_security_group_rules", quota.MaxSecurityGroupRules)
	data.Set("max_security_groups", quota.MaxSecurityGroups)
	data.Set("max_server_group_members", quota.MaxServerGroupMembers)
	data.Set("max_server_groups", quota.MaxServerGroups)
	data.Set("max_server_meta", quota.MaxServerMeta)
	data.Set("max_total_cores", quota.MaxTotalCores)
	data.Set("max_total_floating_ips", quota.MaxTotalFloatingIPs)
	data.Set("max_total_instances", quota.MaxTotalInstances)
	data.Set("max_total_keypairs", quota.MaxTotalKeyPairs)
	data.Set("max_total_ram_size", quota.MaxTotalRamSize)
	data.Set("total_cores_used", quota.TotalCoresUsed)
	data.Set("total_floating_ip_used", quota.TotalFloatingIPUsed)
	data.Set("total_instances_used", quota.TotalInstancesUsed)
	data.Set("total_ram_used", quota.TotalRamUsed)
	data.Set("total_security_groups_used", quota.TotalSecurityGroupsUsed)
	data.Set("total_server_groups_used", quota.TotalServerGroupsUsed)

	return errors
}
