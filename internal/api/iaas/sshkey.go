package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type SSHKeyDetails struct {
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   string `json:"created_at"`
	PublicKey   string `json:"public_key"`
}

type SSHKeyOpts struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type SSHKeyUpdateOpts struct {
}

type SSHKey struct {
	requester *api.Requester
}

// NewSSHKey - init communicator with SSHKey
func NewSSHKey(ctx context.Context) *SSHKey {
	return &SSHKey{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// Find - looking for a sshkey by name
func (s *SSHKey) Find(region, name string) (*SSHKeyDetails, error) {
	sshkeys, err := s.List(region)
	if err != nil {
		return nil, err
	}

	for _, sshkey := range sshkeys {
		if sshkey.Name == name {
			return &sshkey, nil
		}
	}

	return nil, fmt.Errorf("sshkey %v not found", name)
}

// List - return all sshkeys
func (s *SSHKey) List(region string) ([]SSHKeyDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/ssh-keys", ECCEndPoint, Version, region)

	data, err := s.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []SSHKeyDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Create - create a sshkey
func (s *SSHKey) Create(region string, opts *SSHKeyOpts) (*SSHKeyDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/ssh-keys", ECCEndPoint, Version, region)

	data, err := s.requester.Create(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *SSHKeyDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Delete - delete a sshkey
func (s *SSHKey) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/ssh-keys/%v", ECCEndPoint, Version, region, id)
	return s.requester.Delete(endpoint, nil)
}
