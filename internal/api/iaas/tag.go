package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

var (
	ValidTagTypes = []string{
		"server",
		"network",
		"image",
		"float_ip",
		"volume",
		"security_group",
	}
)

type ReplaceTagOpts struct {
	InstanceList []string `json:"instance_list"`
	TagList      []string `json:"tag_list"`
	InstanceType string   `json:"instance_type"`
}

type TagAttachmentOpts struct {
	InstanceId   string `json:"instance_id"`
	InstanceType string `json:"instance_type"`
}

type TagDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TagOpts struct {
	TagName string `json:"tag_name"`
}

type TagUpdateOpts struct {
	TagName string `json:"tag_name"`
}

type Tag struct {
	requester *api.Requester
}

// NewTag - init communicator with tag
func NewTag(ctx context.Context) *Tag {
	return &Tag{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

// Find - looking for a tag by name
func (t *Tag) Find(region, name string) (*TagDetails, error) {
	tags, err := t.List(region)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.Name == name {
			return &tag, nil
		}
	}

	return nil, fmt.Errorf("tag %v not found", name)
}

// List - return all tags
func (t *Tag) List(region string) ([]TagDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/tags", ECCEndPoint, Version, region)

	data, err := t.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []TagDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Create - create a tag
func (t *Tag) Create(region string, opts *TagOpts) (*TagDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/tags", ECCEndPoint, Version, region)

	data, err := t.requester.Create(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *TagDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Delete - delete a tag
func (t *Tag) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/tags/%v", ECCEndPoint, Version, region, id)
	return t.requester.Delete(endpoint, nil)
}

// Update - edit a tag
func (t *Tag) Update(region, id string, opts *TagUpdateOpts) (*TagDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/tags/%v", ECCEndPoint, Version, region, id)

	data, err := t.requester.Put(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *TagDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

// Attach - attach tag to an instance
func (t *Tag) Attach(region string, opts *TagAttachmentOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/tags/attach", ECCEndPoint, Version, region)
	_, err := t.requester.Create(endpoint, opts, nil)
	return err
}

// Detach - detach tag from an instance
func (t *Tag) Detach(region string, opts *TagAttachmentOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/tags/detach", ECCEndPoint, Version, region)
	_, err := t.requester.Create(endpoint, opts, nil)
	return err
}

// Replace - replace a list of tags with instance list (for a list of instances)
func (t *Tag) Replace(region string, opts *ReplaceTagOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/tags/batch", ECCEndPoint, Version, region)
	_, err := t.requester.Create(endpoint, opts, nil)
	return err
}
