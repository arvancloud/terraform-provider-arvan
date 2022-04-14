package provider

import (
	"context"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/iaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {

	// The actual provider
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				Description: "The API Key for API operations. You can retrieve this\n" +
					"from the ArvanCloud dashboard.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"arvan_iaas_abrak":                       iaas.ResourceAbrak(),
			"arvan_iaas_abrak_rename":                iaas.ResourceAbrakRename(),
			"arvan_iaas_abrak_shutdown":              iaas.ResourceAbrakShutDown(),
			"arvan_iaas_abrak_turn_on":               iaas.ResourceAbrakTurnOn(),
			"arvan_iaas_abrak_reboot":                iaas.ResourceAbrakReboot(),
			"arvan_iaas_abrak_rescue":                iaas.ResourceAbrakRescue(),
			"arvan_iaas_abrak_rebuild":               iaas.ResourceAbrakRebuild(),
			"arvan_iaas_abrak_change_flavor":         iaas.ResourceAbrakChangeFlavor(),
			"arvan_iaas_abrak_change_disk_size":      iaas.ResourceAbrakChangeDiskSize(),
			"arvan_iaas_abrak_snapshot":              iaas.ResourceAbrakSnapshot(),
			"arvan_iaas_abrak_add_security_group":    iaas.ResourceAbrakAddSecurityGroup(),
			"arvan_iaas_abrak_remove_security_group": iaas.ResourceAbrakRemoveSecurityGroup(),
			"arvan_iaas_abrak_reset_root_password":   iaas.ResourceAbrakResetRootPassword(),
			"arvan_iaas_abrak_change_public_ip":      iaas.ResourceAbrakChangePublicIP(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"arvan_iaas_abrak":          iaas.DatasourceAbrak(),
			"arvan_iaas_network":        iaas.DatasourceNetwork(),
			"arvan_iaas_image":          iaas.DatasourceImage(),
			"arvan_iaas_volume":         iaas.DatasourceVolume(),
			"arvan_iaas_security_group": iaas.DatasourceSecurityGroup(),
		},
		ConfigureContextFunc: providerConfigure,
	}
	return provider
}

// providerConfigure returns (client.Client, diag.Diagnostics)
func providerConfigure(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	c, err := client.NewClient(&client.Config{
		ApiKey: d.Get("api_key").(string),
	})

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, nil
}
