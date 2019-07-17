package jenkins

import (
	"encoding/json"
	"gopkg.in/resty.v1"
)

const (
	pathCrumb = "/crumbIssuer/api/json"
)

func (c *Client) csrf() (*resty.Request, error) {
	field, crumb, err := c.crumb()
	if err != nil {
		return nil, err
	}
	r := c.c.R()
	r.SetHeader(field, crumb)
	return r, nil
}

func (c *Client) crumb() (field string, value string, err error) {
	resp, err := c.c.R().
		Get(pathCrumb)
	if err != nil {
		return "", "", err
	}
	if !resp.IsSuccess() {
		return "", "", ErrNot2xxStatusCode
	}
	var data map[string]string
	if err = json.Unmarshal(resp.Body(), &data); err != nil {
		return "", "", err
	}
	field = data["crumbRequestField"]
	value = data["crumb"]
	return
}
