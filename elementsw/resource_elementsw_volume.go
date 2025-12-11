package elementsw

import (
	"fmt"
	"log"
	"strconv"

	"encoding/json"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/terraform-provider-solidfire/elementsw/element/jsonrpc"
)

// validateVolumeAttributes enforces max 10 KV pairs and max 1000 bytes JSON
func validateVolumeAttributes(val interface{}, key string) (warns []string, errs []error) {
	   m, ok := val.(map[string]interface{})
	   if !ok {
			   errs = append(errs, fmt.Errorf("attributes must be a map of string to string"))
			   return
	   }
	   if len(m) > 10 {
			   errs = append(errs, fmt.Errorf("attributes cannot have more than 10 key-value pairs"))
	   }
	   // Convert to JSON and check size
	   b, err := json.Marshal(m)
	   if err != nil {
			   errs = append(errs, fmt.Errorf("attributes could not be JSON encoded: %v", err))
	   } else if len(b) > 1000 {
			   errs = append(errs, fmt.Errorf("attributes JSON encoding exceeds 1000 bytes (actual: %d)", len(b)))
	   }
   return
}
// CreateVolumeRequest the users input for creating a Volume
type CreateVolumeRequest struct {
	   Name                    string           `structs:"name"`
	   AccountID               int              `structs:"accountID"`
	   TotalSize               int              `structs:"totalSize"`
	   Enable512E              bool             `structs:"enable512e"`
	   Attributes              interface{}      `structs:"attributes"`
	   QOS                     QualityOfService `structs:"qos,omitempty"`
	   QoSPolicyID             int              `structs:"qosPolicyID,omitempty"`
	   AssociateWithQoSPolicy  bool             `structs:"associateWithQoSPolicy,omitempty"`
	   EnableSnapMirrorReplication bool          `structs:"enableSnapMirrorReplication,omitempty"`
	   FifoSize                string           `structs:"fifoSize,omitempty"`
	   MinFifoSize             string           `structs:"minFifoSize,omitempty"`
}

// CreateVolumeResult the api results for creating a volume
type CreateVolumeResult struct {
	   VolumeID   int    `json:"volumeID"`
	   Volume     volume `json:"volume"`
	   QoSPolicyID *int  `json:"qosPolicyID"`
}

// DeleteVolumeRequest the user input for deleteing a volume
type DeleteVolumeRequest struct {
	VolumeID int `structs:"volumeID"`
}

// ModifyVolumeRequest the user input for modify a volume
type ModifyVolumeRequest struct {
	VolumeID   int              `structs:"volumeID"`
	AccountID  int              `structs:"accountID"`
	Attributes interface{}      `structs:"attributes"`
	QOS        QualityOfService `structs:"qos"`
	TotalSize  int              `structs:"totalSize"`
}

// QualityOfService quailty of service information
type QualityOfService struct {
	MinIOPS   int `structs:"minIOPS"`
	MaxIOPS   int `structs:"maxIOPS"`
	BurstIOPS int `structs:"burstIOPS"`
}

func resourceElementSwVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwVolumeCreate,
		Read:   resourceElementSwVolumeRead,
		Update: resourceElementSwVolumeUpdate,
		Delete: resourceElementSwVolumeDelete,
		Exists: resourceElementSwVolumeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"total_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable512e": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"min_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"burst_iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		   "attributes": {
			   Type:     schema.TypeMap,
			   Optional: true,
			   Elem:     &schema.Schema{Type: schema.TypeString},
			   ValidateFunc: validateVolumeAttributes,
		   },
			"iqn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceElementSwVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %#v", d)
	client := meta.(*Client)

	volume := CreateVolumeRequest{}

	if v, ok := d.GetOk("name"); ok {
		volume.Name = v.(string)
	} else {
		return fmt.Errorf("name argument is required")
	}

	   if a, ok := d.GetOk("account_id"); ok {
			   accountID, ok := a.(int)
			   if !ok {
					   return fmt.Errorf("account_id must be an integer")
			   }
			   volume.AccountID = accountID
	   } else {
			   return fmt.Errorf("account_id argument is required")
	   }

	if v, ok := d.GetOk("total_size"); ok {
		volume.TotalSize = v.(int)
	} else {
		return fmt.Errorf("total_size argument is required")
	}

	if v, ok := d.GetOkExists("enable512e"); ok {
		volume.Enable512E = v.(bool)
	} else {
		return fmt.Errorf("enable512e argument is required")
	}

	   if v, ok := d.GetOk("min_iops"); ok {
			   volume.QOS.MinIOPS = v.(int)
	   }

	   if v, ok := d.GetOk("max_iops"); ok {
			   volume.QOS.MaxIOPS = v.(int)
	   }

	   if v, ok := d.GetOk("burst_iops"); ok {
			   volume.QOS.BurstIOPS = v.(int)
	   }

	   if v, ok := d.GetOk("attributes"); ok {
			   // Validate again at runtime (defensive)
			   if warns, errs := validateVolumeAttributes(v, "attributes"); len(errs) > 0 {
					   return fmt.Errorf("invalid attributes: %v", errs)
			   } else if len(warns) > 0 {
					   log.Printf("attributes validation warning: %v", warns)
			   }
			   volume.Attributes = v
	   }

	resp, err := createVolume(client, volume)
	if err != nil {
		log.Print("Error creating volume")
		return err
	}

	d.SetId(fmt.Sprintf("%v", resp.VolumeID))
	d.Set("iqn", resp.Volume.Iqn)
	log.Printf("Created volume: %v %v", volume.Name, resp.VolumeID)

	return resourceElementSwVolumeRead(d, meta)
}

func createVolume(client *Client, request CreateVolumeRequest) (CreateVolumeResult, error) {
	   params := structs.Map(request)

	   // Remove qos if qosPolicyID is set and associateWithQoSPolicy is true
	   if request.QoSPolicyID > 0 && request.AssociateWithQoSPolicy {
			   delete(params, "qos")
	   }

	   log.Printf("Parameters: %v", params)

	   response, err := client.CallAPIMethod("CreateVolume", params)
	   if err != nil {
			   log.Print("CreateVolume request failed")
			   return CreateVolumeResult{}, err
	   }

	   var result CreateVolumeResult
	   if err := json.Unmarshal([]byte(*response), &result); err != nil {
			   log.Print("Failed to unmarshall response from CreateVolume")
			   return CreateVolumeResult{}, err
	   }
	   return result, nil
}

func resourceElementSwVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading volume: %#v", d)
	client := meta.(*Client)

	volumes := listVolumesRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	s[0] = convID
	volumes.Volumes = s

	res, err := client.listVolumes(volumes)
	if err != nil {
		return err
	}

	if len(res.Volumes) != 1 {
	   return fmt.Errorf("expected one volume to be found. response contained %v results", len(res.Volumes))
	}

	d.Set("name", res.Volumes[0].Name)

	return nil
}

func resourceElementSwVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume access group %#v", d)
	client := meta.(*Client)

	volume := ModifyVolumeRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	volume.VolumeID = convID

	   if a, ok := d.GetOk("account_id"); ok {
			   accountID, ok := a.(int)
			   if !ok {
					   return fmt.Errorf("account_id must be an integer")
			   }
			   volume.AccountID = accountID
	   } else {
			   return fmt.Errorf("account_id argument is required")
	   }

	if v, ok := d.GetOk("total_size"); ok {
		volume.TotalSize = v.(int)
	}

	   if v, ok := d.GetOk("min_iops"); ok {
			   volume.QOS.MinIOPS = v.(int)
	   }

	   if v, ok := d.GetOk("max_iops"); ok {
			   volume.QOS.MaxIOPS = v.(int)
	   }

	   if v, ok := d.GetOk("burst_iops"); ok {
			   volume.QOS.BurstIOPS = v.(int)
	   }

	   if v, ok := d.GetOk("attributes"); ok {
			   if warns, errs := validateVolumeAttributes(v, "attributes"); len(errs) > 0 {
					   return fmt.Errorf("invalid attributes: %v", errs)
			   } else if len(warns) > 0 {
					   log.Printf("attributes validation warning: %v", warns)
			   }
			   volume.Attributes = v
	   }

	err := updateVolume(client, volume)
	if err != nil {
		return err
	}

	return nil
}

func updateVolume(client *Client, request ModifyVolumeRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("ModifyVolume", params)
	if err != nil {
		log.Print("ModifyVolume request failed")
		return err
	}

	return nil
}

func resourceElementSwVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume access group: %#v", d)
	client := meta.(*Client)

	volume := DeleteVolumeRequest{}

	id := d.Id()
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	volume.VolumeID = convID

	deleteErr := deleteVolume(client, volume)
	if deleteErr != nil {
		return deleteErr
	}

	purgeErr := purgeDeletedVolume(client, volume)
	if purgeErr != nil {
		return purgeErr
	}

	return nil
}

func deleteVolume(client *Client, request DeleteVolumeRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("DeleteVolume", params)
	if err != nil {
		log.Print("DeleteVolume request failed")
		return err
	}

	return nil
}

func purgeDeletedVolume(client *Client, request DeleteVolumeRequest) error {
	params := structs.Map(request)

	_, err := client.CallAPIMethod("PurgeDeletedVolume", params)
	if err != nil {
		log.Print("PurgeDeletedVolume request failed")
		return err
	}

	return nil
}

func resourceElementSwVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume: %#v", d)
	client := meta.(*Client)

	volumes := listVolumesRequest{}

	id := d.Id()
	s := make([]int, 1)
	convID, convErr := strconv.Atoi(id)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	s[0] = convID
	volumes.Volumes = s

	res, err := client.listVolumes(volumes)
	if err != nil {
		if err, ok := err.(*jsonrpc.ResponseError); ok {
			if err.Name == "xUnknown" {
				d.SetId("")
				return false, nil
			}
		}
		return false, err
	}

	if len(res.Volumes) != 1 {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
