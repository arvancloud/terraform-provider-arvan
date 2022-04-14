package iaas

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arvancloud/terraform-provider-arvan/internal/api"
)

type VolumeAttachment struct {
	ID           string `json:"id"`
	Device       string `json:"device"`
	ServerId     string `json:"server_id"`
	ServerName   string `json:"server_name"`
	VolumeId     string `json:"volume_id"`
	AttachmentId string `json:"attachment_id"`
	AttachedAt   string `json:"attached_at"`
	HostName     string `json:"host_name"`
}

type VolumeDetails struct {
	ID             string             `json:"id"`
	Size           int                `json:"size"`
	Status         string             `json:"status"`
	CreatedAt      string             `json:"created_at"`
	Description    string             `json:"description"`
	VolumeTypeName string             `json:"volume_type_name"`
	SnapshotId     string             `json:"snapshot_id"`
	SourceVolumeId string             `json:"source_volume_id"`
	Bootable       string             `json:"bootable"`
	Name           string             `json:"name"`
	Attachments    []VolumeAttachment `json:"attachments"`
}

type VolumeOpts struct {
	Name        string `json:"name"`
	Size        int    `json:"size"`
	Description string `json:"description"`
}

type VolumeUpdateOpts struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VolumeAttachmentOpts struct {
	ServerId string `json:"server_id"`
	VolumeId string `json:"volume_id"`
}

type VolumeLimitDetails struct {
	TotalSnapshotsUsed       int `json:"total_snapshots_used"`
	MaxTotalBackups          int `json:"max_total_backups"`
	MaxTotalVolumeGigabytes  int `json:"max_total_volume_gigabytes"`
	MaxTotalSnapshots        int `json:"max_total_snapshots"`
	MaxTotalBackupGigabytes  int `json:"max_total_backup_gigabytes"`
	TotalBackupGigabytesUsed int `json:"total_backup_gigabytes_used"`
	MaxTotalVolumes          int `json:"max_total_volumes"`
	TotalVolumesUsed         int `json:"total_volumes_used"`
	TotalBackupsUsed         int `json:"total_backups_used"`
	TotalGigabytesUsed       int `json:"total_gigabytes_used"`
}

type Volume struct {
	requester *api.Requester
}

func NewVolume(ctx context.Context) *Volume {
	return &Volume{
		requester: ctx.Value(api.RequesterContext).(*api.Requester),
	}
}

func (v *Volume) Find(region, name string) (*string, error) {
	volumes, err := v.List(region)
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		if volume.Name == name {
			return &volume.ID, nil
		}
	}

	return nil, fmt.Errorf("volume %v not found", name)
}

func (v *Volume) List(region string) ([]VolumeDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes", ECCEndPoint, Version, region)

	data, err := v.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details []VolumeDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

func (v *Volume) Create(region string, opts *VolumeOpts) (*VolumeDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes", ECCEndPoint, Version, region)

	data, err := v.requester.Create(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *VolumeDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

func (v *Volume) Detach(region, opts *VolumeAttachmentOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/detach", ECCEndPoint, Version, region)
	_, err := v.requester.Update(endpoint, opts, nil)
	return err
}

func (v *Volume) Attach(region, opts *VolumeAttachmentOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/attach", ECCEndPoint, Version, region)
	_, err := v.requester.Update(endpoint, opts, nil)
	return err
}

func (v *Volume) Limits(region string) (*VolumeLimitDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/limits", ECCEndPoint, Version, region)

	data, err := v.requester.List(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *VolumeLimitDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

func (v *Volume) Read(region, id string) (*VolumeDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/%v", ECCEndPoint, Version, region, id)

	data, err := v.requester.Read(endpoint, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *VolumeDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}

func (v *Volume) Delete(region, id string) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/%v", ECCEndPoint, Version, region, id)

	_, err := v.requester.DoRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	return nil
}

func (v *Volume) Update(region, id string, opts *VolumeUpdateOpts) (*VolumeDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/%v", ECCEndPoint, Version, region, id)

	data, err := v.requester.Update(endpoint, opts, nil)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var details *VolumeDetails
	err = json.Unmarshal(marshal, &details)
	return details, err
}
