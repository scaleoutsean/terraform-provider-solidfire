package elementsw

import (
	"encoding/json"
	"log"

	"github.com/fatih/structs"
)

type getAccountByIDRequest struct {
	AccountID int `structs:"accountID"`
}

type getAccountByNameRequest struct {
	AccountName string `structs:"username"`
}

type getAccountResult struct {
	Account account `json:"account"`
}

type account struct {
	AccountID       int         `json:"accountID"`
	Attributes      interface{} `json:"attributes"`
	InitiatorSecret string      `json:"initiatorSecret"`
	Status          string      `json:"status"`
	TargetSecret    string      `json:"targetSecret"`
	Username        string      `json:"username"`
}

func (c *Client) GetAccountByID(id int) (account, error) {
	   params := structs.Map(getAccountByIDRequest{AccountID: id})
	   return c.GetAccountDetails(params, "GetAccountByID")
}


func (c *Client) GetAccountDetails(params map[string]interface{}, method string) (account, error) {
	   response, err := c.CallAPIMethod(method, params)
	   if err != nil {
			   log.Print(method + " request failed")
			   return account{}, err
	   }

	   var result getAccountResult
	   if err := json.Unmarshal([]byte(*response), &result); err != nil {
			   log.Print("Failed to unmarshal response from GetAccountByID")
			   return account{}, err
	   }

	   return result.Account, nil
}

