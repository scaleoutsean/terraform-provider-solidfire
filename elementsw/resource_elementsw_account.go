package elementsw

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scaleoutsean/solidfire-go/sdk"
)

func resourceElementSwAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceElementSwAccountCreate,
		Read:   resourceElementSwAccountRead,
		Update: resourceElementSwAccountUpdate,
		Delete: resourceElementSwAccountDelete,
		Exists: resourceElementSwAccountExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"initiator_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"target_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceElementSwAccountCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating account: %#v", d)
	client := meta.(*Client)

	req := sdk.AddAccountRequest{}

	if v, ok := d.GetOk("username"); ok {
		req.Username = v.(string)
	} else {
		return fmt.Errorf("username argument is required")
	}

	if v, ok := d.GetOk("initiator_secret"); ok {
		req.InitiatorSecret = v.(string)
	}

	if v, ok := d.GetOk("target_secret"); ok {
		req.TargetSecret = v.(string)
	}

	client.initOnce.Do(client.init)
	resp, sdkErr := client.sdkClient.AddAccount(context.TODO(), &req)
	if sdkErr != nil {
		log.Printf("Error creating account: %s", sdkErr.Detail)
		return sdkErr
	}

	d.SetId(fmt.Sprintf("%v", resp.Account.AccountID))

	log.Printf("Created account: %v %v", req.Username, resp.Account.AccountID)

	return resourceElementSwAccountRead(d, meta)
}

func resourceElementSwAccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading account: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}

	res, err := client.GetAccountByID(convID)
	if err != nil {
		log.Print("GetAccountByID failed")
		return err
	}

	d.Set("username", res.Username)

	// Since we drop secrets in GetAccountByID, we don't update them here
	// to avoid clearing them in the state if they were set during Create/Update.
	// However, if the user wants them to be totally gone from UI/logs,
	// maybe we should set them to empty?
	// The instruction says "drop from the response body ... to avoid exposing them in UI or logs".
	// If they are in the state, they ARE in the UI of terraform (masked as sensitive).
	// If I don't Set them, they stay as they were (probably what we want).
	/*
		if res.InitiatorSecret != "" {
			d.Set("initiator_secret", res.InitiatorSecret)
		}
		if res.TargetSecret != "" {
			d.Set("target_secret", res.TargetSecret)
		}
	*/

	return nil
}

func resourceElementSwAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating account %#v", d)
	client := meta.(*Client)

	req := sdk.ModifyAccountRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	req.AccountID = convID

	if d.HasChange("username") {
		req.Username = d.Get("username").(string)
	}

	if d.HasChange("initiator_secret") {
		req.InitiatorSecret = d.Get("initiator_secret").(string)
	}

	if d.HasChange("target_secret") {
		req.TargetSecret = d.Get("target_secret").(string)
	}

	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.ModifyAccount(context.TODO(), &req)
	if sdkErr != nil {
		return sdkErr
	}

	return nil
}

func resourceElementSwAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting account: %#v", d)
	client := meta.(*Client)

	req := sdk.RemoveAccountRequest{}

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return fmt.Errorf("id argument is required")
	}
	req.AccountID = convID

	client.initOnce.Do(client.init)
	_, sdkErr := client.sdkClient.RemoveAccount(context.TODO(), &req)
	if sdkErr != nil {
		return sdkErr
	}

	return nil
}

func resourceElementSwAccountExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of account: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	convID, convErr := strconv.ParseInt(id, 10, 64)

	if convErr != nil {
		return false, fmt.Errorf("id argument is required")
	}

	_, err := client.GetAccountByID(convID)
	if err != nil {
		// In account.go, we return sdkErr directly.
		if sdkErr, ok := err.(*sdk.SdkError); ok {
			if fmt.Sprintf("%s", sdkErr.Detail) == "500:xUnknownAccount" {
				d.SetId("")
				return false, nil
			}
		}
		log.Print("AccountExists failed")
		return false, err
	}

	return true, nil
}
