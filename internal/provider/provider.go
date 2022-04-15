package provider

import (
	"context"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/iaas/datasources"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/iaas/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a *schema.Provider
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
			"arvan_iaas_abrak":                       resources.ResourceAbrak(),
			"arvan_iaas_abrak_rename":                resources.ResourceAbrakRename(),
			"arvan_iaas_abrak_action":                resources.ResourceAbrakAction(),
			"arvan_iaas_abrak_rebuild":               resources.ResourceAbrakRebuild(),
			"arvan_iaas_abrak_change_flavor":         resources.ResourceAbrakChangeFlavor(),
			"arvan_iaas_abrak_change_disk_size":      resources.ResourceAbrakChangeDiskSize(),
			"arvan_iaas_abrak_snapshot":              resources.ResourceAbrakSnapshot(),
			"arvan_iaas_abrak_add_security_group":    resources.ResourceAbrakAddSecurityGroup(),
			"arvan_iaas_abrak_remove_security_group": resources.ResourceAbrakRemoveSecurityGroup(),
			"arvan_iaas_cdn_security_group":          resources.ResourceSecurityGroupCdn(),
			"arvan_iaas_security_group_rule":         resources.ResourceSecurityGroupRule(),
			"arvan_iaas_security_group":              resources.ResourceSecurityGroup(),
			"arvan_iaas_sshkey":                      resources.ResourceSSHKey(),
			"arvan_iaas_floatip":                     resources.ResourceFloatIP(),
			"arvan_iaas_ptr":                         resources.ResourcePtr(),
			"arvan_iaas_tag":                         resources.ResourceTag(),
			"arvan_iaas_tag_attach":                  resources.ResourceTagAttach(),
			"arvan_iaas_tag_replace":                 resources.ResourceTagReplaceBatch(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"arvan_iaas_abrak":          datasources.DatasourceAbrak(),
			"arvan_iaas_network":        datasources.DatasourceNetwork(),
			"arvan_iaas_image":          datasources.DatasourceImage(),
			"arvan_iaas_volume":         datasources.DatasourceVolume(),
			"arvan_iaas_security_group": datasources.DatasourceSecurityGroup(),
			"arvan_iaas_sshkey":         datasources.DatasourceSSHKey(),
			"arvan_iaas_quota":          datasources.DatasourceQuota(),
			"arvan_iaas_tag":            datasources.DatasourceQuota(),
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
