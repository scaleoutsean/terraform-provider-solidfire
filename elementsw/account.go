package elementsw

import (
	"context"

	"github.com/scaleoutsean/solidfire-go/sdk"
)

type account struct {
	AccountID       int64       `json:"accountID"`
	Attributes      interface{} `json:"attributes"`
	InitiatorSecret string      `json:"initiatorSecret"`
	Status          string      `json:"status"`
	TargetSecret    string      `json:"targetSecret"`
	Username        string      `json:"username"`
}

func (c *Client) GetAccountByID(id int64) (account, error) {
	req := sdk.GetAccountByIDRequest{
		AccountID: id,
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.GetAccountByID(context.TODO(), &req)
	if sdkErr != nil {
		return account{}, sdkErr
	}

	return c.processAccount(res.Account), nil
}

func (c *Client) GetAccountByName(name string) (account, error) {
	req := sdk.GetAccountByNameRequest{
		Username: name,
	}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.GetAccountByName(context.TODO(), &req)
	if sdkErr != nil {
		return account{}, sdkErr
	}

	return c.processAccount(res.Account), nil
}

func (c *Client) ListAccounts() ([]account, error) {
	req := sdk.ListAccountsRequest{}
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.ListAccounts(context.TODO(), &req)
	if sdkErr != nil {
		return nil, sdkErr
	}

	var accounts []account
	for _, a := range res.Accounts {
		accounts = append(accounts, c.processAccount(a))
	}
	return accounts, nil
}

func (c *Client) processAccount(sdkAccount sdk.Account) account {
	// Security requirement: Always drop initiatorSecret and targetSecret from response
	// to avoid exposing them in UI or logs.
	return account{
		AccountID:       sdkAccount.AccountID,
		Attributes:      sdkAccount.Attributes,
		InitiatorSecret: "", // Dropped for security
		TargetSecret:    "", // Dropped for security
		Status:          sdkAccount.Status,
		Username:        sdkAccount.Username,
	}
}
