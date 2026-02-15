package elementsw

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
)

// Config is a struct for user input
type configStuct struct {
	User            string
	Password        string
	ElementSwServer string
	APIVersion      string
}

// Client contain the api endpoint
// Removed unused type clientStuct

// APIError is any error the api gives
type APIError struct {
	ID    int `json:"id"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"error"`
}

// Client is the main function to connect to the APi
func (c *configStuct) clientFun() (*Client, error) {
	host := c.ElementSwServer
	if strings.Contains(host, "://") {
		u, err := url.Parse(host)
		if err == nil {
			host = u.Host
		}
	}
	client := &Client{
		Host:     host,
		Username: c.User,
		Password: c.Password,
		HTTPTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true},
		},
	}

	client.SetAPIVersion(c.APIVersion)

	return client, nil
}
