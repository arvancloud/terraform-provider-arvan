package iaas

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DatasourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNetworkRead,
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
				Description: "name of network",
			},
		},
	}
}

func datasourceNetworkRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var errors diag.Diagnostics
	c := meta.(*client.Client).Iaas

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	name := data.Get("name").(string)
	id, err := c.Network.FindNetworkId(region, name)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("network %v not found", name),
		})
		return errors
	}

	data.SetId(*id)
	return errors
}
