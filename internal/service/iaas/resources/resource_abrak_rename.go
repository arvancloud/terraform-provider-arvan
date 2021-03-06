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

func ResourceAbrakRename() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakRenameCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakRenameUpdate,
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
			"abrak_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "uuid of abrak",
				ValidateFunc: validation.IsUUID,
			},
			"new_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "new name of abrak",
			},
		},
	}
}

func resourceAbrakRenameCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	uuid := data.Get("abrak_uuid").(string)
	name := data.Get("new_name").(string)

	err := c.Server.Actions.Rename(region, uuid, name)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not rename server %v to %v", uuid, name),
		})
		return errors
	}

	data.SetId(uuid)
	return errors
}

func resourceAbrakRenameUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChanges("new_name") {
		return resourceAbrakRenameCreate(ctx, data, meta)
	}
	return nil
}
