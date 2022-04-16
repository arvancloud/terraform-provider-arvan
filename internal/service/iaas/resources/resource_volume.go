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

func ResourceVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of volume",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The size of the volume, in gigabytes (GB)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description of volume",
			},
		},
	}
}

func resourceVolumeCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	// volume Options
	volume := &iaas.VolumeOpts{
		Name: data.Get("name").(string),
		Size: data.Get("size").(int),
	}

	if description, ok := data.GetOk("description"); ok {
		volume.Description = description.(string)
	}

	response, err := c.Volume.Create(region, volume)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ID)
	data.Set("name", response.Name)
	data.Set("size", response.Size)
	data.Set("status", response.Status)
	data.Set("created_at", response.CreatedAt)
	data.Set("description", response.Description)
	data.Set("volume_type_name", response.VolumeTypeName)
	data.Set("snapshot_id", response.SnapshotId)
	data.Set("source_volume_id", response.SourceVolumeId)
	data.Set("bootable", response.Bootable)

	return errors
}

func resourceVolumeUpdate(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	if data.HasChanges("name", "description") {
		// volumeUpdateOpts Options
		volumeUpdateOpts := &iaas.VolumeUpdateOpts{
			Name: data.Get("name").(string),
		}

		if description, ok := data.GetOk("description"); ok {
			volumeUpdateOpts.Description = description.(string)
		}

		err := c.Volume.Update(region, data.Id(), volumeUpdateOpts)
		if err != nil {
			return diag.FromErr(err)
		}

		data.Set("name", volumeUpdateOpts.Name)
		data.Set("description", volumeUpdateOpts.Description)
	}

	return errors
}

func resourceVolumeDelete(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	volumes, err := c.Volume.List(region)
	for _, volume := range volumes {
		if volume.ID == data.Id() {
			if len(volume.Attachments) > 0 {
				for _, attachment := range volume.Attachments {
					opts := &iaas.VolumeAttachmentOpts{
						VolumeId: data.Id(),
						ServerId: attachment.ServerId,
					}
					err = c.Volume.Detach(region, opts)
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}

	// TODO: waiting to be un-available

	err = c.Volume.Delete(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete Volume %v", data.Id()),
		})
		return errors
	}

	return nil
}
