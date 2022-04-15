package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type QuotaDetails struct {
	MaxImageMeta            int `json:"max_image_meta"`
	MaxPersonality          int `json:"max_personality"`
	MaxSecurityGroupRules   int `json:"max_security_group_rules"`
	MaxSecurityGroups       int `json:"max_security_groups"`
	MaxServerGroupMembers   int `json:"max_server_group_members"`
	MaxServerGroups         int `json:"max_server_groups"`
	MaxServerMeta           int `json:"max_server_meta"`
	MaxTotalCores           int `json:"max_total_cores"`
	MaxTotalFloatingIPs     int `json:"max_total_floating_ips"`
	MaxTotalInstances       int `json:"max_total_instances"`
	MaxTotalKeyPairs        int `json:"max_total_keypairs"`
	MaxTotalRamSize         int `json:"max_total_ram_size"`
	TotalCoresUsed          int `json:"total_cores_used"`
	TotalFloatingIPUsed     int `json:"total_floating_ip_used"`
	TotalInstancesUsed      int `json:"total_instances_used"`
	TotalRamUsed            int `json:"total_ram_used"`
	TotalSecurityGroupsUsed int `json:"total_security_groups_used"`
	TotalServerGroupsUsed   int `json:"total_server_groups_used"`
}

type Quota struct {
	requester *api.Requester
}

// NewQuota - init communicator with Quota
func NewQuota(ctx context.Context) *Quota {
	return &Quota{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// List - return all quotas
func (q *Quota) List(region string) ([]QuotaDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/quotas", ECCEndPoint, Version, region)

	data, err := q.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []QuotaDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}
