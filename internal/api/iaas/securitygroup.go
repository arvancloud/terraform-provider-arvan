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

const (
	IngressDirection = "ingress"
	EgressDirection  = "egress"
)

var (
	SupportedDirections = []string{
		IngressDirection,
		EgressDirection,
	}
)

type RuleDetails struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	EtherType   string `json:"ether_type"`
	GroupID     string `json:"group_id"`
	IP          string `json:"ip"`
	PortStart   int    `json:"port_start"`
	PortEnd     int    `json:"port_end"`
	Protocol    string `json:"protocol"`
}

type SecurityGroupAbrak struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	IPs  []string `json:"ips"`
}

type SecurityGroupDetails struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Abraks      []SecurityGroupAbrak `json:"abraks"`
	Description string               `json:"description"`
	ReadOnly    bool                 `json:"readonly"`
	RealName    string               `json:"real_name"`
	Rules       []RuleDetails        `json:"rules"`
	Tags        []TagDetails         `json:"tags"`
}

type SecurityGroupOpts struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SecurityGroupRuleOpts struct {
	Direction   string   `json:"direction"`
	PortFrom    string   `json:"port_from"`
	PortTo      string   `json:"port_to"`
	Protocol    string   `json:"protocol"`
	IPs         []string `json:"ips"`
	Description string   `json:"description"`
}

type SecurityGroup struct {
	requester *api.Requester
}

// NewSecurityGroup - init communicator with SecurityGroup
func NewSecurityGroup(ctx context.Context) *SecurityGroup {
	return &SecurityGroup{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// List - return all securityGroups
func (s *SecurityGroup) List(region string) (details []SecurityGroupDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities", ECCEndPoint, Version, region)

	data, err := s.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Find - looking for a securityGroup by name
func (s *SecurityGroup) Find(region, name string) (*SecurityGroupDetails, error) {
	securityGroups, err := s.List(region)
	if err != nil {
		return nil, err
	}
	for _, sg := range securityGroups {
		if sg.RealName == name {
			return &sg, nil
		}
	}
	return nil, fmt.Errorf("securityGroup %v not found", name)
}

// Create - create a securityGroup
func (s *SecurityGroup) Create(region string, opts *SecurityGroupOpts) (details *SecurityGroupDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities", ECCEndPoint, Version, region)

	data, err := s.requester.Create(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, &details)
	return details, err
}

// CreateCdn - create a securityGroup
func (s *SecurityGroup) CreateCdn(region string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities/cdn", ECCEndPoint, Version, region)
	_, err = s.requester.DoRequest("POST", endpoint, nil)
	return err
}

// Delete - delete a securityGroup
func (s *SecurityGroup) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities/%v", ECCEndPoint, Version, region, id)
	return s.requester.Delete(endpoint, nil)
}

// Read - get details of a securityGroup
func (s *SecurityGroup) Read(region, id string) (details *SecurityGroupDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities/security-rules/%v", ECCEndPoint, Version, region, id)

	data, err := s.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(marshal, &details)
	return details, err
}

// CreateRule - create a rule for a securityGroup
func (s *SecurityGroup) CreateRule(region, id string, opts *SecurityGroupRuleOpts) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities/security-rules/%v", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, opts, nil)
	return err
}

// DeleteRule - delete rule of a securityGroup
func (s *SecurityGroup) DeleteRule(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/securities/security-rules/%v", ECCEndPoint, Version, region, id)
	return s.requester.Delete(endpoint, nil)
}
