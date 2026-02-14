package elementsw

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/scaleoutsean/solidfire-go/sdk"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var ourlog = logrus.WithFields(logrus.Fields{
	"prefix": "main",
})

func init() {
	logrus.SetFormatter(new(prefixed.TextFormatter))
}

// A Client to interact with the Element API
type Client struct {
	Host                  string
	Username              string
	Password              string
	MaxConcurrentRequests int
	HTTPTransport         http.RoundTripper

	apiVersion string

	initOnce     sync.Once
	sdkClient    *sdk.SFClient
	requestSlots chan int
}

func (c *Client) GetClusterInfo() (*sdk.GetClusterInfoResult, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.GetClusterInfo(context.TODO())
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res, nil
}

func (c *Client) GetClusterVersionInfo() (*sdk.GetClusterVersionInfoResult, error) {
	c.initOnce.Do(c.init)
	res, sdkErr := c.sdkClient.GetClusterVersionInfo(context.TODO())
	if sdkErr != nil {
		return nil, sdkErr
	}
	return res, nil
}

// CallAPIMethod can be used to make a request to any Element API method, receiving results as raw JSON
func (c *Client) CallAPIMethod(method string, params map[string]interface{}) (*json.RawMessage, error) {
	c.initOnce.Do(c.init)

	c.waitForAvailableSlot()
	defer c.releaseSlot()

	ourlog.WithFields(logrus.Fields{
		"method": method,
		"params": params,
	}).Debug("Calling API")

	// This is a bridge method for migration.
	// We'll eventually replace individual calls with SDK methods.
	// For now, we can use a generic call if the SDK supports it, or start migrating methods.
	// Looking at the SDK, it has MakeSFCall which is exported if it was SfClient.MakeSFCall but wait...
	// base_methods.go: func (sfClient *SFClient) MakeSFCall(...)
	// It IS exported (starts with Uppercase).

	var res interface{}
	_, sdkErr := c.sdkClient.MakeSFCall(context.TODO(), method, 1, params, &res)
	if sdkErr != nil {
		return nil, fmt.Errorf("%s: %s", sdkErr.Code, sdkErr.Detail)
	}

	resultBits, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	rawRes := json.RawMessage(resultBits)

	ourlog.WithFields(logrus.Fields{
		"method": method,
	}).Debug("Received successful API response")
	return &rawRes, nil
}

func (c *Client) init() {
	if c.MaxConcurrentRequests == 0 {
		c.MaxConcurrentRequests = 6
	}
	c.requestSlots = make(chan int, c.MaxConcurrentRequests)

	c.sdkClient = &sdk.SFClient{}
	// Note: solidfire-go's Connect method uses SSL and InsecureSkipVerify by default.
	// It also builds the URL from host and version.
	c.sdkClient.Connect(context.TODO(), c.Host, c.GetAPIVersion(), c.Username, c.Password)
}

// SetAPIVersion for the client to use for requests to the Element API
func (c *Client) SetAPIVersion(apiVersion string) {
	c.apiVersion = apiVersion
}

// GetAPIVersion returns the API version that will be used for Element API requests
func (c *Client) GetAPIVersion() string {
	if c.apiVersion == "" {
		return "1.0"
	}
	return c.apiVersion
}

func (c *Client) waitForAvailableSlot() {
	c.requestSlots <- 1
}

func (c *Client) releaseSlot() {
	<-c.requestSlots
}
