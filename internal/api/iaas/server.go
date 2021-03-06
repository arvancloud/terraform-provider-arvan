package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

const (
	ServerActiveStatus           = "active"
	ServerBuildStatus            = "build"
	ServerDeletedStatus          = "deleted"
	ServerErrorStatus            = "error"
	ServerHardRebootStatus       = "hard_reboot"
	ServerMigratingStatus        = "migrating"
	ServerPasswordStatus         = "password"
	ServerPausedStatus           = "paused"
	ServerRebootStatus           = "reboot"
	ServerReBuildStatus          = "rebuild"
	ServerRescueStatus           = "rescue"
	ServerResizeStatus           = "resize"
	ServerRevertResizeStatus     = "revert_resize"
	ServerShelvedStatus          = "shelved"
	ServerShelvedOffloadedStatus = "shelved_offloaded"
	ServerShutoffStatus          = "shutoff"
	ServerSoftDeletedStatus      = "soft_deleted"
	ServerSuspendedStatus        = "suspended"
	ServerUnknownStatus          = "unknown"
	ServerVerifyResizeStatus     = "verify_resize"
)

type ServerSecurityGroupOpts struct {
	Name string `json:"name"`
}

type ServerOpts struct {
	Name           string                    `json:"name"`
	NetworkIds     []string                  `json:"network_ids"`
	FlavorId       string                    `json:"flavor_id"`
	ImageId        string                    `json:"image_id"`
	SecurityGroups []ServerSecurityGroupOpts `json:"security_groups"`
	SshKey         bool                      `json:"ssh_key,omitempty"`
	KeyName        any                       `json:"key_name" default:"0"`
	Count          int                       `json:"count" default:"1"`
	CreateType     string                    `json:"create_type,omitempty"`
	DiskSize       int                       `json:"disk_size"`
	InitScript     string                    `json:"init_script" default:""`
	HAEnabled      bool                      `json:"ha_enabled" default:"false"`
}

type ServerFlavor struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	RAM   int32  `json:"ram"`
	Swap  string `json:"swap"`
	VCPUs int32  `json:"vcpus"`
	Disk  int32  `json:"disk"`
}

type ServerImage struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	MinDisk   int32             `json:"min_disk"`
	MinRam    int32             `json:"min_ram"`
	OS        string            `json:"os"`
	OSVersion string            `json:"os_version"`
	Progress  int32             `json:"progress"`
	Size      int64             `json:"size"`
	Status    string            `json:"status"`
	Created   string            `json:"created"`
	UserName  string            `json:"username"`
	MetaData  map[string]string `json:"metadata"`
}

type ServerSecurityGroup struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Name        string        `json:"name"`
	ReadOnly    string        `json:"readonly"`
	RealName    string        `json:"real_name"`
	Rules       []RuleDetails `json:"rules"`
}

type ServerDetails struct {
	ID             string                     `json:"id"`
	Name           string                     `json:"name"`
	Flavor         *ServerFlavor              `json:"flavor"`
	Status         string                     `json:"status"`
	Image          *ServerImage               `json:"image"`
	Created        string                     `json:"created"`
	Password       string                     `json:"password"`
	TaskState      *string                    `json:"task_state"`
	KeyName        string                     `json:"key_name"`
	ArNext         string                     `json:"ar_next"`
	SecurityGroups []ServerSecurityGroup      `json:"security_groups"`
	Addresses      map[string][]ServerAddress `json:"addresses"`
	Tags           []TagDetails               `json:"tags"`
	HAEnabled      bool                       `json:"ha_enabled"`
}

type ServerAddress struct {
	MAC     string `json:"mac_addr"`
	Version string `json:"version"`
	Addr    string `json:"addr"`
	Type    string `json:"type"`
}

type Server struct {
	requester *api.Requester
	Actions   *ServerActions
}

// NewServer - init communicator with Server
func NewServer(ctx context.Context) *Server {
	return &Server{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
		Actions:   NewServerActions(ctx),
	}
}

// Find - looking for a server by name
func (s *Server) Find(region, name string) (*ServerDetails, error) {
	servers, err := s.List(region)
	if err != nil {
		return nil, err
	}

	for _, server := range servers {
		if server.Name == name {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("server %v not found", name)
}

// List - return all servers
func (s *Server) List(region string) (details []ServerDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers", ECCEndPoint, Version, region)

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

// Create - create a server
func (s *Server) Create(region string, opts *ServerOpts) (details *ServerDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers", ECCEndPoint, Version, region)

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

// Read - get details of a server
func (s *Server) Read(region, id string) (details *ServerDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v", ECCEndPoint, Version, region, id)

	data, err := s.requester.Read(endpoint, nil)
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

// Delete - delete a server
func (s *Server) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v", ECCEndPoint, Version, region, id)
	return s.requester.Delete(endpoint, map[string]string{
		"forceDelete": "true",
	})
}

type ServerOptionsDetails struct {
	RegionId              int    `json:"region_id"`
	RequiresPaymentMethod bool   `json:"requires_payment_method"`
	DropletCount          int    `json:"droplet_count"`
	DropletLimit          int    `json:"droplet_limit"`
	Currency              string `json:"currency"`
	ImageName             string `json:"image_name"`
	ImageVersion          string `json:"image_version"`
	NetworkId             string `json:"network_id"`
}

// Options - return a region options (default network and image)
func (s *Server) Options(region string) (details *ServerOptionsDetails, err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/options", ECCEndPoint, Version, region)

	data, err := s.requester.Read(endpoint, nil)
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
