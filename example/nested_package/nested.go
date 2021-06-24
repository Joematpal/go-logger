package nested

import "errors"

type Client struct {
	Logger interface {
		Errorf(format string, args ...interface{})
	}
}

func (c *Client) Something() error {
	err := errors.New("error")
	c.Logger.Errorf("something=%v", err)
	return err
}
