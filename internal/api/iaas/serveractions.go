package iaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type SnapshotOpts struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ServerActions struct {
	requester *api.Requester
}

// NewServerActions - init communicator with ServerActions
func NewServerActions(ctx context.Context) *ServerActions {
	return &ServerActions{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// Rename - rename the server
func (s *ServerActions) Rename(region, id, newName string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/rename", ECCEndPoint, Version, region, id)

	var requestBody = struct {
		Name string `json:"name"`
	}{
		Name: newName,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.DoRequest("POST", endpoint, bytes.NewBuffer(body))
	return err
}

// ShutDown - shutdown a server
func (s *ServerActions) ShutDown(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/power-off", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}

// TurnOn - turn on a server
func (s *ServerActions) TurnOn(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/power-on", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

// SoftReboot - apply soft reboot for a server
func (s *ServerActions) SoftReboot(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/reboot", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}

// HardReboot - apply hard reboot for a server
func (s *ServerActions) HardReboot(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/hard-reboot", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}

// Rescue - apply rescue on a server
func (s *ServerActions) Rescue(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/rescue", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}

// UnRescue - apply un-rescue on a server
func (s *ServerActions) UnRescue(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/unrescue", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}

// Rebuild - rebuild a server from an image
func (s *ServerActions) Rebuild(region, id, imageId string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/rebuild", ECCEndPoint, Version, region, id)

	var requestBody any = &struct {
		ImageId string `json:"image_id"`
	}{
		ImageId: imageId,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.DoRequest("POST", endpoint, bytes.NewBuffer(body))
	return err
}

// ChangeFlavor - change flavor of a server
func (s *ServerActions) ChangeFlavor(region, id, flavorId string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/resize", ECCEndPoint, Version, region, id)

	var requestBody any = &struct {
		FlavorId string `json:"flavor_id"`
	}{
		FlavorId: flavorId,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.DoRequest("POST", endpoint, bytes.NewBuffer(body))
	return err
}

// ChangeDiskSize - change disk size of a server
func (s *ServerActions) ChangeDiskSize(region, id string, size int) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/resizeRoot", ECCEndPoint, Version, region, id)

	var requestBody any = &struct {
		NewSize int `json:"new_size"`
	}{
		NewSize: size,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.DoRequest("PUT", endpoint, bytes.NewBuffer(body))
	return err
}

// Snapshot - create a snapshot from a server
func (s *ServerActions) Snapshot(region, id string, opts *SnapshotOpts) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/%v/snapshot", ECCEndPoint, Version, region, id)

	body, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	_, err = s.requester.DoRequest("POST", endpoint, bytes.NewBuffer(body))
	return err
}

// AddSecurityGroup - assign a securityGroup to a server
func (s *ServerActions) AddSecurityGroup(region, id, securityGroupId string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/add-security-group", ECCEndPoint, Version, region, id)

	var requestBody any = &struct {
		SecurityGroupId string `json:"security_group_id"`
	}{
		SecurityGroupId: securityGroupId,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.DoRequest("POST", endpoint, bytes.NewBuffer(body))
	return err
}

// RemoveSecurityGroup - remove a securityGroup from a server
func (s *ServerActions) RemoveSecurityGroup(region, id, securityGroupId string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/remove-security-group", ECCEndPoint, Version, region, id)

	var requestBody any = &struct {
		SecurityGroupId string `json:"security_group_id"`
	}{
		SecurityGroupId: securityGroupId,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.DoRequest("POST", endpoint, bytes.NewBuffer(body))
	return err
}

// ChangePublicIP - change public-ip of a server
func (s *ServerActions) ChangePublicIP(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/change-public-ip", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}

// AddPublicIP - add public-ip of a server
func (s *ServerActions) AddPublicIP(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/add-public-ip", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}

// ResetRootPassword - reset root password of a server
func (s *ServerActions) ResetRootPassword(region, id string) (err error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/reset-root-password", ECCEndPoint, Version, region, id)
	_, err = s.requester.Create(endpoint, nil, nil)
	return err
}
