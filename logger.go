package logler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/syslog"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Client struct {
	Trace            *log.Logger
	info             *log.Logger
	warn             *log.Logger
	error            *log.Logger
	emergency        *log.Logger
	component        string
	logglySampleRate int
	syslogClient     *syslog.Writer
}

type Options struct {
	Component        string
	LogglySampleRate int
	MinimalLog       bool
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
		//Connect to local syslog server (the syslog server should be configured to send to logstash)
		syslogclient, err := syslog.New(syslog.LOG_ERR, opts.Component)
		if err != nil {
			log.Println(err.Error)
		}
		result.syslogClient = syslogclient
		result.logglySampleRate = opts.LogglySampleRate
		if len(opts.Component) > 0 {
			result.component = opts.Component
		}
	}
	return result
}

func (c *Client) Info(msg map[string]interface{}) {
	if msg, err := getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		c.info.Println(string(j))

		if c.syslogClient != nil {
			if c.logglySampleRate == 100 {
				c.syslogClient.Info(string(j))
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.syslogClient.Info(string(j))
				}
			}
		}
	}
}

func (c *Client) Warn(msg map[string]interface{}) {
	if msg, err := getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		c.warn.Println(string(j))

		if c.syslogClient != nil {
			if c.logglySampleRate == 100 {
				c.syslogClient.Warning(string(j))
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.syslogClient.Warning(string(j))
				}
			}
		}
	}
}

func (c *Client) Error(msg map[string]interface{}) {
	if msg, err := getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		c.error.Println(string(j))
		if c.syslogClient != nil {
			if c.logglySampleRate == 100 {
				fmt.Println("***about to syslog***")
				c.syslogClient.Err(string(j))
				fmt.Println("***syslog done***")
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.syslogClient.Err(string(j))
				}
			}
		}
	}
}

func (c *Client) Emergency(msg map[string]interface{}) {
	if msg, err := getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		c.emergency.Println(string(j))
		if c.syslogClient != nil {
			c.syslogClient.Emerg(string(j))
		}
	}
}

func getMessage(msg map[string]interface{}) (map[string]interface{}, error) {
	if msg != nil {
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[1])
		file, line := f.FileLine(pc[1])
		tmp := strings.Split(file, string(filepath.Separator))
		msg["filename"] = tmp[len(tmp)-1]
		msg["line"] = line
		msg["func"] = f.Name()
		return msg, nil
	}
	return nil, errors.New("message log nil message")
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
