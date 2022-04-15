package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type SubnetOpts struct {
	Name          string   `json:"name"`
	SubnetIP      string   `json:"subnet_ip"`
	EnableGateway bool     `json:"enable_gateway"`
	SubnetGateway string   `json:"subnet_gateway"`
	Dhcp          string   `json:"dhcp"`
	DnsServers    []string `json:"dns_servers"`
}

type NetworkAttachOpts struct {
	ServerId           string `json:"server_id"`
	IP                 string `json:"ip"`
	EnablePortSecurity bool   `json:"enablePortSecurity"`
}

type NetworkDetachOpts struct {
	ServerId string `json:"server_id"`
}

type PoolDetails struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type NetworkServerAddress struct {
	MacAddress string `json:"mac_address"`
	Version    string `json:"version"`
	Address    string `json:"addr"`
	Type       string `json:"type"`
	IsPublic   bool   `json:"is_public"`
}

type NetworkServerIP struct {
	IsFloatIP           bool   `json:"float_ip"`
	IP                  string `json:"ip"`
	MacAddress          string `json:"mac_address"`
	PortId              string `json:"port_id"`
	PortSecurityEnabled bool   `json:"port_security_enabled"`
	Ptr                 string `json:"ptr"`
	Public              bool   `json:"public"`
	SubnetId            string `json:"subnet_id"`
	SubnetName          string `json:"subnet_name"`
	Version             string `json:"version"`
}

type NetworkServerPublicIP struct {
	SubnetId  string `json:"subnet_id"`
	IPAddress string `json:"ip_address"`
}

type NetworkServer struct {
	ID             string                            `json:"id"`
	Name           string                            `json:"name"`
	Addresses      map[string][]NetworkServerAddress `json:"addresses"`
	IPs            []NetworkServerIP                 `json:"ips"`
	PublicIPs      []NetworkServerPublicIP           `json:"public_ip"`
	SecurityGroups []string                          `json:"security_groups"`
}

type SubnetDetails struct {
	ID              string          `json:"id"`
	NetworkId       string          `json:"network_id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	IPVersion       string          `json:"ip_version"`
	Cidr            string          `json:"cidr"`
	GatewayIP       string          `json:"gateway_ip"`
	DnsNameservers  []string        `json:"dns_nameservers"`
	AllocationPool  []PoolDetails   `json:"allocation_pool"`
	HostRoutes      string          `json:"host_routes"`
	EnableDhcp      bool            `json:"enable_dhcp"`
	TenantId        string          `json:"tenant_id"`
	ProjectId       string          `json:"project_id"`
	IPv6AddressMode string          `json:"ipv6_address_mode"`
	IPv6RaMode      string          `json:"ipv6_ra_mode"`
	SubnetPoolId    string          `json:"subnetpool_id"`
	ServiceType     string          `json:"service_type"`
	RevisionNumber  int             `json:"revision_number"`
	Tags            []TagDetails    `json:"tags"`
	Servers         []NetworkServer `json:"servers"`
}

type NetworkDetails struct {
	ID                    string          `json:"id"`
	Name                  string          `json:"name"`
	Description           string          `json:"description"`
	AdminStateUp          bool            `json:"admin_state_up"`
	Shared                bool            `json:"shared"`
	Status                string          `json:"status"`
	Subnets               []SubnetDetails `json:"subnets"`
	TenantId              string          `json:"tenant_id"`
	DhcpIP                string          `json:"dhcp_ip"`
	UpdatedAt             string          `json:"updated_at"`
	CreatedAt             string          `json:"created_at"`
	IPv4AddressScope      string          `json:"ipv4_address_scope"`
	IPv6AddressScope      string          `json:"ipv6_address_scope"`
	QosPolicyId           string          `json:"qos_policy_id"`
	RevisionNumber        int             `json:"revision_number"`
	RouteExternal         string          `json:"route:external"`
	Mtu                   int             `json:"mtu"`
	PortSecurityEnabled   bool            `json:"port_security_enabled"`
	AvailabilityZoneHints string          `json:"availability_zone_hints"`
	AvailabilityZones     string          `json:"availability_zones"`
	Tags                  []TagDetails    `json:"tags"`
}

type Network struct {
	requester *api.Requester
}

// NewNetwork - init communicator with network
func NewNetwork(ctx context.Context) *Network {
	return &Network{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// List - return all network
func (n *Network) List(region string) ([]NetworkDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/networks", ECCEndPoint, Version, region)

	data, err := n.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []NetworkDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Find - looking for a network by name
func (n *Network) Find(region, name string) (*NetworkDetails, error) {
	networks, err := n.List(region)
	if err != nil {
		return nil, err
	}

	for _, network := range networks {
		if network.Name == name {
			return &network, nil
		}
	}

	return nil, fmt.Errorf("network %v not found", name)
}

// Detach - detach a network from a server
func (n *Network) Detach(region, id string, opts *NetworkDetachOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/networks/%v/detach", ECCEndPoint, Version, region, id)
	_, err := n.requester.Patch(endpoint, opts, nil)
	return err
}

// Attach - attach a network to a server
func (n *Network) Attach(region, id string, opts *NetworkAttachOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/networks/%v/attach", ECCEndPoint, Version, region, id)
	_, err := n.requester.Patch(endpoint, opts, nil)
	return err
}

// ReadSubnet - get subnet details
func (n *Network) ReadSubnet(region, id string) ([]SubnetDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/subnets/%v", ECCEndPoint, Version, region, id)

	data, err := n.requester.Read(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []SubnetDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// CreateSubnet - create a subnet
func (n *Network) CreateSubnet(region string, opts *SubnetOpts) (*SubnetDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/subnets", ECCEndPoint, Version, region)

	data, err := n.requester.Create(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *SubnetDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// UpdateSubnet - edit a subnet
func (n *Network) UpdateSubnet(region, id string, opts *SubnetOpts) (*SubnetDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/subnets/%v", ECCEndPoint, Version, region, id)

	data, err := n.requester.Patch(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *SubnetDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// DeleteSubnet - delete a subnet
func (n *Network) DeleteSubnet(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/subnets/%v", ECCEndPoint, Version, region, id)
	return n.requester.Delete(endpoint, nil)
}
