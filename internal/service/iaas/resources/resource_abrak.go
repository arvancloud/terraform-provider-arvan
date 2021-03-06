package resources

import (
	"context"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
	"time"
)

const (
	AbrakDelay         = 5 * time.Second
	AbrakMinTimeout    = 3 * time.Second
	AbrakCreateTimeout = 10 * time.Minute
	AbrakDeleteTimeout = 10 * time.Minute
)

func ResourceAbrak() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbrakCreate,
		ReadContext:   resourceAbrakRead,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: resourceAbrakDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(AbrakCreateTimeout),
			Delete: schema.DefaultTimeout(AbrakDeleteTimeout),
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
				Description: "name of abrak",
			},
			"networks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network(s) of abrak",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plan ID of abrak, you can get list of plan IDs of each region from sizes api",
			},
			"image": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Required:    true,
				Description: "image of abrak",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "type of image",
							ValidateFunc: validation.StringInSlice([]string{
								iaas.ImageTypeServer,
								iaas.ImageTypeSnapshot,
								iaas.ImageTypeDistributions,
							}, false),
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name of image",
						},
					},
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "array of security group names",
			},
			"addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "array of abrak addresses",
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSH Key name (for password: 0)",
			},
			"ssh_key": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Can use ssh key",
			},
			"number": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "Count of abraks we want to create",
			},
			"create_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "type of abrak creation, " +
					"a flag which shows client sends arguments or we should read from image",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Disk size of abraks we want to create",
			},
			"init_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Init script",
			},
			"ha_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Instance HA enabled or not",
			},
		},
	}
}

func resourceAbrakCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS
	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	// Abrak Options
	abrak := &iaas.ServerOpts{
		Name:       data.Get("name").(string),
		FlavorId:   data.Get("flavor").(string),
		SshKey:     data.Get("ssh_key").(bool),
		Count:      data.Get("number").(int),
		CreateType: data.Get("create_type").(string),
		DiskSize:   data.Get("disk_size").(int),
		InitScript: data.Get("init_script").(string),
		HAEnabled:  data.Get("ha_enabled").(bool),
	}

	// KeyName
	abrak.KeyName = 0
	if v := data.Get("key_name"); v.(string) != "0" {
		abrak.KeyName = v.(string)
	}

	// ImageID
	if v, ok := data.GetOk("image"); ok && v.(*schema.Set).Len() > 0 {
		var iType, iName string
		for _, imageMap := range v.(*schema.Set).List() {
			imageInfo, ok := imageMap.(map[string]any)
			if !ok {
				errors = append(errors, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "invalid image details",
				})
				return errors
			}

			if v, ok := imageInfo["type"]; ok {
				iType = v.(string)
			}
			if v, ok := imageInfo["name"]; ok {
				iName = v.(string)
			}

			if iType != "" && iName != "" {
				image, err := c.Image.Find(
					region,
					iName,
					iType,
				)
				if err != nil {
					errors = append(errors, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "image not found",
					})
					return errors
				}
				switch image.(type) {
				case iaas.ImageDetails:
					abrak.ImageId = image.(iaas.ImageDetails).ID
				case iaas.ImageServerDetails:
					abrak.ImageId = image.(iaas.ImageServerDetails).ID
				}
				break
			}
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "image block in incomplete (name and type are required)",
			})
			return errors
		}
	} else {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "image block not found or invalid.",
		})
		return errors
	}
	// NetworkIDs
	if networks, ok := data.GetOk("networks"); !ok {
		opts, err := c.Server.Options(region)
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "there is no default network",
			})
			return errors
		}
		abrak.NetworkIds = []string{opts.NetworkId}
	} else {
		var networkIds []string
		for _, networkName := range networks.([]any) {
			network, err := c.Network.Find(region, networkName.(string))
			if err != nil {
				errors = append(errors, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("there is no %v network", networkName.(string)),
				})
				return errors
			}
			networkIds = append(networkIds, network.ID)
		}
		abrak.NetworkIds = networkIds
	}

	// SecurityGroup
	if securityGroups, ok := data.GetOk("security_groups"); !ok {
		securityGroup, err := c.SecurityGroup.Find(region, iaas.DefaultSecurityGroup)
		if err != nil {
			errors = append(errors, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "there is no default security group",
			})
			return errors
		}
		abrak.SecurityGroups = []iaas.ServerSecurityGroupOpts{{
			Name: securityGroup.ID,
		}}
	} else {
		var sgs []iaas.ServerSecurityGroupOpts
		for _, sg := range securityGroups.([]any) {
			securityGroup, err := c.SecurityGroup.Find(region, sg.(string))
			if err != nil {
				errors = append(errors, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("there is no %v security group", sg.(string)),
				})
				return errors
			}
			sgs = append(sgs, iaas.ServerSecurityGroupOpts{
				Name: securityGroup.ID,
			})
		}
		abrak.SecurityGroups = sgs
	}

	response, err := c.Server.Create(region, abrak)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(response.ID)

	_, err = AbrakWaitFroAvailable(ctx, data.Timeout(schema.TimeoutCreate), region, response.ID, c)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAbrakRead(ctx, data, meta)
}

func resourceAbrakRead(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	abrak, err := c.Server.Read(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not retrive abrak %v", data.Id()),
		})
		return errors
	}

	data.Set("flavor", abrak.Flavor.ID)
	data.Set("status", abrak.Status)
	data.Set("created_at", abrak.Created)
	data.Set("addresses", flattenAbrakAddresses(abrak.Addresses))
	data.Set("tags", flattenAbrakTags(abrak.Tags))
	data.Set("ha_enabled", abrak.HAEnabled)
	data.Set("key_name", abrak.KeyName)
	return nil
}

func flattenAbrakTags(tags []iaas.TagDetails) (abrakTags []string) {
	for _, tag := range tags {
		abrakTags = append(abrakTags, tag.Name)
	}
	return abrakTags
}

func flattenAbrakAddresses(addresses map[string][]iaas.ServerAddress) (abrakAddresses []string) {
	for _, network := range addresses {
		for _, address := range network {
			abrakAddresses = append(abrakAddresses, address.Addr)
		}
	}
	return abrakAddresses
}

func resourceAbrakDelete(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS
	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}
	err := c.Server.Delete(region, data.Id())
	if err != nil {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("can not delete server %v", data.Id()),
		})
		return errors
	}

	err = abrakWaitFroDestroy(ctx, data.Timeout(schema.TimeoutDelete), region, data.Id(), c)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func abrakStateRefreshFunc(i *iaas.IaaS, region, id string) resource.StateRefreshFunc {
	return func() (server any, status string, err error) {
		emptyResponse := &iaas.ServerDetails{}

		server, err = i.Server.Read(region, id)
		if err != nil {
			return emptyResponse, iaas.ServerDeletedStatus, nil
		}

		return server, strings.ToLower(server.(*iaas.ServerDetails).Status), err
	}
}

func AbrakWaitFroAvailable(ctx context.Context, timeout time.Duration, region, id string, i *iaas.IaaS) (abrak *iaas.ServerDetails, err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{iaas.ServerBuildStatus, iaas.ServerReBuildStatus, iaas.ServerMigratingStatus},
		Target:     []string{iaas.ServerActiveStatus},
		Refresh:    abrakStateRefreshFunc(i, region, id),
		Timeout:    timeout,
		Delay:      AbrakDelay,
		MinTimeout: AbrakMinTimeout,
	}

	info, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for Abrak (%s) to be ready: %v", id, err)
	}

	return info.(*iaas.ServerDetails), nil
}

func abrakWaitFroDestroy(ctx context.Context, timeout time.Duration, region, id string, i *iaas.IaaS) (err error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{iaas.ServerActiveStatus},
		Target:     []string{iaas.ServerDeletedStatus},
		Refresh:    abrakStateRefreshFunc(i, region, id),
		Timeout:    timeout,
		Delay:      AbrakDelay,
		MinTimeout: AbrakMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for delete Abrak (%s) to be ready: %v", id, err)
	}

	return nil
}
