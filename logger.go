package logler

import (
	"encoding/json"
	"github.com/streamrail/go-loggly"
	"log"
	"os"
)

type Client struct {
	Trace        *log.Logger
	info         *log.Logger
	warn         *log.Logger
	error        *log.Logger
	emergency    *log.Logger
	component    string
	logglyClient *loggly.Client
}

type Options struct {
	LogglyToken string
	Component   string
}

func New(opts *Options) *Client {
	result := &Client{
		info: log.New(os.Stdout,
			"INFO: ",
			log.Ldate|log.Ltime),
		warn: log.New(os.Stdout,
			"WARNING: ",
			log.Ldate|log.Ltime),
		error: log.New(os.Stderr,
			"ERROR: ",
			log.Ldate|log.Ltime),
		emergency: log.New(os.Stderr,
			"Emergency: ",
			log.Ldate|log.Ltime),
		Trace: log.New(os.Stdout,
			"TRACE: ",
			log.Ldate|log.Ltime),
	}
	if opts != nil {
		result.logglyClient = loggly.New(opts.LogglyToken)
	}
	return result
}

func (c *Client) Info(msg map[string]interface{}) {
	j, _ := json.Marshal(msg)
	c.info.Println(string(j))
	c.logglyClient.Info(c.component, msg)
}
func (c *Client) Warn(msg map[string]interface{}) {
	j, _ := json.Marshal(msg)
	c.warn.Println(string(j))
	c.logglyClient.Warn(c.component, msg)
}
func (c *Client) Error(msg map[string]interface{}) {
	j, _ := json.Marshal(msg)
	c.error.Println(string(j))
	c.logglyClient.Error(c.component, msg)
}
func (c *Client) Emergency(msg map[string]interface{}) {
	j, _ := json.Marshal(msg)
	c.emergency.Println(string(j))
	c.logglyClient.Emergency(c.component, msg)
}
