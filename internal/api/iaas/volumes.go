package iaas

import (
	"bytes"
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

func NewVolume(r *api.Requester) *Volume {
	return &Volume{
		requester: r,
	}
}

func (v *Volume) FindVolumeId(region, name string) (*string, error) {
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

	responseBody, err := v.requester.DoRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response *api.SuccessResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var volumeDetails []VolumeDetails
	err = json.Unmarshal(dataBytes, &volumeDetails)
	if err != nil {
		return nil, err
	}

	return volumeDetails, nil
}

func (v *Volume) Create(region string, opts *VolumeOpts) (*VolumeDetails, error) {

	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes", ECCEndPoint, Version, region)

	requestBody, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	responseBody, err := v.requester.DoRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	var response *api.SuccessResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var volumeDetails *VolumeDetails
	err = json.Unmarshal(dataBytes, &volumeDetails)
	if err != nil {
		return nil, err
	}

	return volumeDetails, nil
}

func (v *Volume) Detach(region, opts *VolumeAttachmentOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/detach", ECCEndPoint, Version, region)

	requestBody, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	_, err = v.requester.DoRequest("PATCH", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	return nil
}

func (v *Volume) Attach(region, opts *VolumeAttachmentOpts) error {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/attach", ECCEndPoint, Version, region)

	requestBody, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	_, err = v.requester.DoRequest("PATCH", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	return nil
}

func (v *Volume) Limits(region string) (*VolumeLimitDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/limits", ECCEndPoint, Version, region)

	data, err := v.requester.DoRequest("GET", endpoint, nil)
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

	var volumeLimitDetails *VolumeLimitDetails
	err = json.Unmarshal(dataBytes, &volumeLimitDetails)
	if err != nil {
		return nil, err
	}

	return volumeLimitDetails, nil
}

func (v *Volume) Read(region, id string) (*VolumeDetails, error) {
	endpoint := fmt.Sprintf("/%v/%v/regions/%v/volumes/%v", ECCEndPoint, Version, region, id)

	data, err := v.requester.DoRequest("GET", endpoint, nil)
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

	var volumeDetails *VolumeDetails
	err = json.Unmarshal(dataBytes, &volumeDetails)
	if err != nil {
		return nil, err
	}

	return volumeDetails, nil
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

	requestBody, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	responseBody, err := v.requester.DoRequest("PATCH", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	var response *api.SuccessResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(response.Data)
	if err != nil {
		return nil, err
	}

	var volumeDetails *VolumeDetails
	err = json.Unmarshal(dataBytes, &volumeDetails)
	if err != nil {
		return nil, err
	}

	return volumeDetails, nil
}
