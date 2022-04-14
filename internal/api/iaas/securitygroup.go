package iaas

import (
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

func NewSecurityGroup(r *api.Requester) *SecurityGroup {
	return &SecurityGroup{
		requester: r,
	}
}

func (s *SecurityGroup) List(region string) ([]SecurityGroupDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities", ECCEndPoint, Version, region)

	data, err := s.requester.DoRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response *api.SuccessResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var securityGroups []SecurityGroupDetails
	err = json.Unmarshal(dataBytes, &securityGroups)
	if err != nil {
		return nil, err
	}

	return securityGroups, nil
}

func (s *SecurityGroup) FindSecurityGroupId(region, name string) (*string, error) {
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
