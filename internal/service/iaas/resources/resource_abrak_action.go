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

const (
	SoftRebootAction     = "soft-reboot"
	HardRebootAction     = "hard-reboot"
	ShutDownAction       = "shutdown"
	TurnOnAction         = "turnon"
	RescueAction         = "rescue"
	UnRescueAction       = "unrescue"
	ResetPasswordAction  = "reset-password"
	ChangePublicIPAction = "change-ip"
	AddPublicIPAction    = "add-public-ip"
)

var (
	SupportedActions = []string{
		SoftRebootAction,
		HardRebootAction,
		ShutDownAction,
		TurnOnAction,
		RescueAction,
		UnRescueAction,
		ResetPasswordAction,
		ChangePublicIPAction,
		AddPublicIPAction,
	}
)

func ResourceAbrakAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakActionCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: resourceAbrakActionUpdate,
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
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "action",
				ValidateFunc: validation.StringInSlice(SupportedActions, false),
			},
		},
	}
}

func resourceAbrakActionCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	var err error
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	uuid := data.Get("uuid").(string)

	action := data.Get("action").(string)
	switch action {
	case SoftRebootAction:
		err = c.Server.Actions.SoftReboot(region, uuid)
	case HardRebootAction:
		err = c.Server.Actions.HardReboot(region, uuid)
	case ShutDownAction:
		err = c.Server.Actions.ShutDown(region, uuid)
	case TurnOnAction:
		err = c.Server.Actions.TurnOn(region, uuid)
	case RescueAction:
		err = c.Server.Actions.Rescue(region, uuid)
	case UnRescueAction:
		err = c.Server.Actions.UnRescue(region, uuid)
	case ResetPasswordAction:
		err = c.Server.Actions.ResetRootPassword(region, uuid)
	case ChangePublicIPAction:
		err = c.Server.Actions.ChangePublicIP(region, uuid)
	case AddPublicIPAction:
		err = c.Server.Actions.AddPublicIP(region, uuid)
	default:
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("could not apply action %v on server %v", action, uuid),
		})
		return errors
	}

	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(uuid)
	return errors
}

func resourceAbrakActionUpdate(ctx context.Context, data *schema.ResourceData, meta any) diag.Diagnostics {
	if data.HasChange("action") {
		return resourceAbrakActionCreate(ctx, data, meta)
	}
	return nil
}
