package iaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type ServerActions struct {
	requester *api.Requester
}

func NewServerActions(ctx context.Context) *ServerActions {
	return &ServerActions{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

func (s *ServerActions) Rename(region, id, newName string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/rename", ECCEndPoint, Version, region, id)

	var requestBody any = &struct {
		Name string `json:"name"`
	}{
		Name: newName,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.Create(endpoint, bytes.NewBuffer(body), nil)
	return err
}

func (s *ServerActions) ShutDown(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/power-off", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

func (s *ServerActions) TurnOn(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/power-on", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

func (s *ServerActions) SoftReboot(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/reboot", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

func (s *ServerActions) HardReboot(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/hard-reboot", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

func (s *ServerActions) Rescue(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/rescue", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

func (s *ServerActions) UnRescue(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/unrescue", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

func (s *ServerActions) Rebuild(region, id, imageId string) error {
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

	_, err = s.requester.Create(endpoint, bytes.NewBuffer(body), nil)
	return err
}

func (s *ServerActions) ChangeFlavor(region, id, flavorId string) error {
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

	_, err = s.requester.Create(endpoint, bytes.NewBuffer(body), nil)
	return err
}

func (s *ServerActions) ChangeDiskSize(region, id string, size int) error {
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

	_, err = s.requester.Create(endpoint, bytes.NewBuffer(body), nil)
	return err
}

func (s *ServerActions) Snapshot(region, id, name string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/snapshot", ECCEndPoint, Version, region, id)

	var requestBody any = &struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	_, err = s.requester.Create(endpoint, bytes.NewBuffer(body), nil)
	return err
}

func (s *ServerActions) AddSecurityGroup(region, id, securityGroupId string) error {
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

	_, err = s.requester.Create(endpoint, bytes.NewBuffer(body), nil)
	return err
}

func (s *ServerActions) RemoveSecurityGroup(region, id, securityGroupId string) error {
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

	_, err = s.requester.Create(endpoint, bytes.NewBuffer(body), nil)
	return err
}

func (s *ServerActions) ChangePublicIP(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/change-public-ip", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}

func (s *ServerActions) ResetRootPassword(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/servers/%v/reset-root-password", ECCEndPoint, Version, region, id)
	_, err := s.requester.Create(endpoint, nil, nil)
	return err
}
