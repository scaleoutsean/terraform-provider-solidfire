package elementsw

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceElementSwAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceElementSwAccountRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initiator_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"target_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceElementSwAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	var accountID int64
	var username string

	if v, ok := d.GetOk("account_id"); ok {
		accountID = int64(v.(int))
	}
	if v, ok := d.GetOk("username"); ok {
		username = v.(string)
	}

	if accountID == 0 && username == "" {
		return fmt.Errorf("one of account_id or username must be set")
	}

	var acc account
	var err error

	if accountID != 0 {
		acc, err = client.GetAccountByID(accountID)
	} else {
		acc, err = client.GetAccountByName(username)
	}

	if err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	d.SetId(strconv.FormatInt(acc.AccountID, 10))
	d.Set("account_id", int(acc.AccountID))
	d.Set("username", acc.Username)
	d.Set("status", acc.Status)
	// Security: CHAP secrets are usually sensitive and we might want to drop them as per instructions
	// but data source often needs them for automation.
	// The instructions said "drop initiatorSecret and targetSecret from response body... unless explicitly stated"
	// However, for a data source, if they are requested, they should be there but marked Sensitive.
	// Actually, the persona says: "Unless explicitly stated otherwise, always drop initiatorSecret and targetSecret"
	// I'll drop them to be safe and conservative as requested.

	// d.Set("initiator_secret", acc.InitiatorSecret) // Dropped
	// d.Set("target_secret", acc.TargetSecret)       // Dropped

	return nil
}
