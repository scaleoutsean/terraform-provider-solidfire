package elementsw

import (
	"encoding/json"
)

// ModifySnapshotRequest for SolidFire API
type ModifySnapshotRequest struct {
   SnapshotID              int     `json:"snapshotID"`
   Name                    string  `json:"name,omitempty"`
   SnapMirrorLabel         string  `json:"snapMirrorLabel,omitempty"`
   EnableRemoteReplication *bool   `json:"enableRemoteReplication,omitempty"`
   ExpirationTime          string  `json:"expirationTime,omitempty"`
}

func (c *Client) ModifySnapshot(req ModifySnapshotRequest) error {
   params := map[string]interface{}{
      "snapshotID": req.SnapshotID,
   }
   if req.Name != "" { params["name"] = req.Name }
   if req.SnapMirrorLabel != "" { params["snapMirrorLabel"] = req.SnapMirrorLabel }
   if req.EnableRemoteReplication != nil { params["enableRemoteReplication"] = *req.EnableRemoteReplication }
   if req.ExpirationTime != "" { params["expirationTime"] = req.ExpirationTime }
   _, err := c.CallAPIMethod("ModifySnapshot", params)
   return err
}

// ModifyGroupSnapshotRequest for SolidFire API
type ModifyGroupSnapshotRequest struct {
   GroupSnapshotID         int     `json:"groupSnapshotID"`
   Name                    string  `json:"name,omitempty"`
   EnableRemoteReplication *bool   `json:"enableRemoteReplication,omitempty"`
   ExpirationTime          string  `json:"expirationTime,omitempty"`
}

func (c *Client) ModifyGroupSnapshot(req ModifyGroupSnapshotRequest) error {
   params := map[string]interface{}{
      "groupSnapshotID": req.GroupSnapshotID,
   }
   if req.Name != "" { params["name"] = req.Name }
   if req.EnableRemoteReplication != nil { params["enableRemoteReplication"] = *req.EnableRemoteReplication }
   if req.ExpirationTime != "" { params["expirationTime"] = req.ExpirationTime }
   _, err := c.CallAPIMethod("ModifyGroupSnapshot", params)
   return err
}
// CreateGroupSnapshotRequest for SolidFire API
type CreateGroupSnapshotRequest struct {
   Volumes             []int                  `json:"volumes"`
   Name                string                 `json:"name,omitempty"`
   EnableRemoteReplication *bool              `json:"enableRemoteReplication,omitempty"`
   EnsureSerialCreation   *bool               `json:"ensureSerialCreation,omitempty"`
   Retention           string                 `json:"retention,omitempty"`
   ExpirationTime      string                 `json:"expirationTime,omitempty"`
   Attributes          map[string]interface{} `json:"attributes,omitempty"`
}

type CreateGroupSnapshotResponse struct {
   GroupSnapshotID   int    `json:"groupSnapshotID"`
   GroupSnapshotUUID string `json:"groupSnapshotUUID"`
}

func (c *Client) CreateGroupSnapshot(req CreateGroupSnapshotRequest) (*CreateGroupSnapshotResponse, error) {
   params := map[string]interface{}{
     "volumes": req.Volumes,
   }
   if req.Name != "" { params["name"] = req.Name }
   if req.EnableRemoteReplication != nil { params["enableRemoteReplication"] = *req.EnableRemoteReplication }
   if req.EnsureSerialCreation != nil { params["ensureSerialCreation"] = *req.EnsureSerialCreation }
   if req.Retention != "" { params["retention"] = req.Retention }
   if req.ExpirationTime != "" { params["expirationTime"] = req.ExpirationTime }
   if req.Attributes != nil { params["attributes"] = req.Attributes }
   result, err := c.CallAPIMethod("CreateGroupSnapshot", params)
   if err != nil {
     return nil, err
   }
   var resp CreateGroupSnapshotResponse
   if err := json.Unmarshal(*result, &resp); err != nil {
     return nil, err
   }
   return &resp, nil
}

// CreateSnapshotRequest for SolidFire API
type CreateSnapshotRequest struct {
   VolumeID                int                    `json:"volumeID"`
   SnapshotID              *int                   `json:"snapshotID,omitempty"`
   SnapMirrorLabel         string                 `json:"snapMirrorLabel,omitempty"`
   Name                    string                 `json:"name,omitempty"`
   EnableRemoteReplication *bool                  `json:"enableRemoteReplication,omitempty"`
   EnsureSerialCreation    *bool                  `json:"ensureSerialCreation,omitempty"`
   Retention               string                 `json:"retention,omitempty"`
   ExpirationTime          string                 `json:"expirationTime,omitempty"`
   Attributes              map[string]interface{} `json:"attributes,omitempty"`
}

type CreateSnapshotResponse struct {
   Snapshot struct {
     SnapshotID int    `json:"snapshotID"`
     CreateTime string `json:"createTime"`
   } `json:"snapshot"`
}

func (c *Client) CreateSnapshot(req CreateSnapshotRequest) (*CreateSnapshotResponse, error) {
   params := map[string]interface{}{
     "volumeID": req.VolumeID,
   }
   if req.SnapshotID != nil { params["snapshotID"] = *req.SnapshotID }
   if req.SnapMirrorLabel != "" { params["snapMirrorLabel"] = req.SnapMirrorLabel }
   if req.Name != "" { params["name"] = req.Name }
   if req.EnableRemoteReplication != nil { params["enableRemoteReplication"] = *req.EnableRemoteReplication }
   if req.EnsureSerialCreation != nil { params["ensureSerialCreation"] = *req.EnsureSerialCreation }
   if req.Retention != "" { params["retention"] = req.Retention }
   if req.ExpirationTime != "" { params["expirationTime"] = req.ExpirationTime }
   if req.Attributes != nil { params["attributes"] = req.Attributes }
   result, err := c.CallAPIMethod("CreateSnapshot", params)
   if err != nil {
     return nil, err
   }
   var resp CreateSnapshotResponse
   if err := json.Unmarshal(*result, &resp); err != nil {
     return nil, err
   }
   return &resp, nil
}

// DeleteGroupSnapshotRequest for SolidFire API
type DeleteGroupSnapshotRequest struct {
   GroupSnapshotID int  `json:"groupSnapshotID"`
   SaveMembers     bool `json:"saveMembers"`
}

func (c *Client) DeleteGroupSnapshot(req DeleteGroupSnapshotRequest) error {
   params := map[string]interface{}{
     "groupSnapshotID": req.GroupSnapshotID,
     "saveMembers": req.SaveMembers,
   }
   _, err := c.CallAPIMethod("DeleteGroupSnapshot", params)
   return err
}

// DeleteSnapshotRequest for SolidFire API
type DeleteSnapshotRequest struct {
   SnapshotID int `json:"snapshotID"`
}

func (c *Client) DeleteSnapshot(req DeleteSnapshotRequest) error {
   params := map[string]interface{}{
     "snapshotID": req.SnapshotID,
   }
   _, err := c.CallAPIMethod("DeleteSnapshot", params)
   return err
}

// ListGroupSnapshotsRequest for SolidFire API
// https://docs.netapp.com/sfe-122/topic/com.netapp.doc.sfe-api/GUID-6B2B2B2B-2B2B-2B2B-2B2B-2B2B2B2B2B2B.html
// Volumes is a slice of volume IDs

type ListGroupSnapshotsRequest struct {
   Volumes []int `json:"volumes,omitempty"`
}

type ListSnapshotsRequest struct{}

// GroupSnapshotMember represents a member snapshot in a group
// (fields may need adjustment based on API response)
type GroupSnapshotMember struct {
   SnapshotID int    `json:"snapshotID"`
   VolumeID   int    `json:"volumeID"`
   CreateTime string `json:"createTime"`
}

type GroupSnapshot struct {
   GroupSnapshotID   int                   `json:"groupSnapshotID"`
   GroupSnapshotUUID string                `json:"groupSnapshotUUID"`
   Members           []GroupSnapshotMember `json:"members"`
}

type ListGroupSnapshotsResponse struct {
   GroupSnapshots []GroupSnapshot `json:"groupSnapshots"`
}

type Snapshot struct {
   SnapshotID  int    `json:"snapshotID"`
   VolumeID    int    `json:"volumeID"`
   CreateTime  string `json:"createTime"`
}

type ListSnapshotsResponse struct {
   Snapshots []Snapshot `json:"snapshots"`
}

// ListGroupSnapshots calls the SolidFire API to list group snapshots
func (c *Client) ListGroupSnapshots(req ListGroupSnapshotsRequest) (*ListGroupSnapshotsResponse, error) {
   params := map[string]interface{}{}
   if len(req.Volumes) > 0 {
      params["volumes"] = req.Volumes
   }
   result, err := c.CallAPIMethod("ListGroupSnapshots", params)
   if err != nil {
      return nil, err
   }
   var resp ListGroupSnapshotsResponse
   if err := json.Unmarshal(*result, &resp); err != nil {
      return nil, err
   }
   return &resp, nil
}

// ListSnapshots calls the SolidFire API to list individual snapshots
func (c *Client) ListSnapshots(req ListSnapshotsRequest) (*ListSnapshotsResponse, error) {
   result, err := c.CallAPIMethod("ListSnapshots", map[string]interface{}{})
   if err != nil {
      return nil, err
   }
   var resp ListSnapshotsResponse
   if err := json.Unmarshal(*result, &resp); err != nil {
      return nil, err
   }
   return &resp, nil
}
