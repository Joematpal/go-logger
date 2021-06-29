package nested

import (
	"errors"

	logger "github.com/digital-dream-labs/go-logger/example/logger"
)

type Client struct {
	Logger logger.Logger
}

func (c *Client) Something() error {
	err := errors.New("error")
	c.Logger.Errorf("something=%v", err)
	return err
}
