package datasources

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DatasourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSecurityGroupRead,
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
				Description: "name of security group",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "description of security group",
			},
			"real_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "real name of security group",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "real name of security group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description of rule",
						},
						"direction": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "direction of rule",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ip of rule",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ip of rule",
						},
					},
				},
			},
		},
	}
}

func datasourceSecurityGroupRead(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	name := data.Get("name").(string)
	securityGroup, err := c.SecurityGroup.Find(region, name)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("security group %v not found", name),
		})
		return errors
	}

	data.SetId(securityGroup.ID)

	data.Set("description", data.Get("description").(string))
	data.Set("real_name", data.Get("real_name").(string))

	// TODO: we have to complete the variables

	return errors
}
