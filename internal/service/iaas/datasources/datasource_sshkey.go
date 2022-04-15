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

func DatasourceSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSSHKeyRead,
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
				Description: "name of ssh-key",
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "public key",
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "fingerprint of ssh-key",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "creation time of ssh-key",
			},
		},
	}
}

func datasourceSSHKeyRead(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
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
	sshKey, err := c.SSHKey.Find(region, name)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("ssh-key %v not found", name),
		})
		return errors
	}

	data.SetId(sshKey.Name)
	data.Set("public_key", sshKey.PublicKey)
	data.Set("fingerprint", sshKey.Fingerprint)
	data.Set("created_at", sshKey.CreatedAt)
	return errors
}
