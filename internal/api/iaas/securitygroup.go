package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

const (
	DefaultSecurityGroup = "arDefault"
)

type SecurityGroupRule struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	EtherType   string `json:"ether_type"`
	GroupID     string `json:"group_id"`
	IP          string `json:"ip"`
	PortStart   int32  `json:"port_start"`
	PortEnd     int32  `json:"port_end"`
	Protocol    string `json:"protocol"`
}

type SecurityGroupAbrak struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	IPs  []string `json:"ips"`
}

type SecurityGroupDetails struct {
	Abraks      []SecurityGroupAbrak `json:"abraks"`
	Description string               `json:"description"`
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	ReadOnly    bool                 `json:"readonly"`
	RealName    string               `json:"real_name"`
	Rules       []SecurityGroupRule  `json:"rules"`
	Tags        []Tag                `json:"tags"`
}

type SecurityGroup struct {
	requester *api.Requester
}

func NewSecurityGroup(ctx context.Context) *SecurityGroup {
	return &SecurityGroup{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

func (s *SecurityGroup) List(region string) ([]SecurityGroupDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities", ECCEndPoint, Version, region)

	data, err := s.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []SecurityGroupDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

func (s *SecurityGroup) Find(region, name string) (*string, error) {
	securityGroups, err := s.List(region)
	if err != nil {
		return nil, err
	}
	for _, sg := range securityGroups {
		if sg.RealName == name {
			return &sg.ID, nil
		}
	}
	return nil, fmt.Errorf("securityGroup %v not found", name)
}
