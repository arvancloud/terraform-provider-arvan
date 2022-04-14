package iaas

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

func ResourceAbrakAddSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakAddSecurityGroupCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakAddSecurityGroupUpdate,
		DeleteContext: helper.DummyResourceAction,
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
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "uuid of abrak",
			},
			"security_group_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "uuid of security group",
			},
		},
	}
}

func resourceAbrakAddSecurityGroupCreate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
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

	id := data.Get("uuid").(string)

	securityGroupId := data.Get("security_group_uuid").(string)
	err := c.Server.Actions.AddSecurityGroup(region, id, securityGroupId)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not add security group %v server %v", securityGroupId, id),
		})
		return errors
	}

	data.SetId(id)
	return errors
}

func resourceAbrakAddSecurityGroupUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChange("security_group_uuid") {
		return resourceAbrakAddSecurityGroupCreate(ctx, data, meta)
	}
	return nil
}
