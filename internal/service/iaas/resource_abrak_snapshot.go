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

func ResourceAbrakSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakSnapshotCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakSnapshotUpdate,
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
			"snapshot_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "snapshot name of abrak",
			},
		},
	}
}

func resourceAbrakSnapshotCreate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
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

	snapshotName := data.Get("snapshot_name").(string)
	err := c.Server.Actions.Snapshot(region, id, snapshotName)
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not take snapshot server %v to %v", id, snapshotName),
		})
		return errors
	}

	data.SetId(id)
	return errors
}

func resourceAbrakSnapshotUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChanges("snapshot_name") {
		return resourceAbrakSnapshotCreate(ctx, data, meta)
	}
	return nil
}
