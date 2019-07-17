package jenkins

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

var (
	// ErrNot2xxStatusCode 代表 response status code 並非 2xx
	ErrNot2xxStatusCode = fmt.Errorf(`expected response status code 2xx.
Use the '--verbose' flag to see the full stacktrace`)
)

// Client 代表一個 jenkins client
type Client struct {
	c   *resty.Client
	log *logrus.Logger
}

// NewClient 產生一個 jenkins client
func NewClient(url string) *Client {
	return &Client{
		log: logrus.StandardLogger(),
		c:   resty.New().SetHostURL(url).SetDisableWarn(true),
	}
}

// SetLogger sets the logger
func (c *Client) SetLogger(log *logrus.Logger) *Client {
	c.log = log
	return c
}

// SetVerbose enables verbose mode
func (c *Client) SetVerbose(v bool) *Client {
	c.c.SetDebug(v)
	return c
}

// SetBasicAuth set the basic auth for jenkins
func (c *Client) SetBasicAuth(username, password string) *Client {
	c.c.SetBasicAuth(username, password)
	return c
}
