package logler

import (
	"encoding/json"
	"errors"
	"github.com/streamrail/go-loggly"
	"log"
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
	logglyClient     *loggly.Client
}

type Options struct {
	LogglyToken      string
	Component        string
	LogglySampleRate int
	LogglyBufferSize int
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
		if len(opts.LogglyToken) > 0 && opts.LogglySampleRate > 0 {
			bufferSize := 100
			if opts.LogglyBufferSize > 0 {
				bufferSize = opts.LogglyBufferSize
			}
			minLog := false
			if opts.MinimalLog {
				minLog = true
			}
			result.logglyClient = loggly.New(opts.LogglyToken, bufferSize, minLog)
			result.logglySampleRate = opts.LogglySampleRate
		}
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

		if c.logglyClient != nil {
			if c.logglySampleRate == 100 {
				c.logglyClient.Info(c.component, msg)
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.logglyClient.Info(c.component, msg)
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

		if c.logglyClient != nil {
			if c.logglySampleRate == 100 {
				c.logglyClient.Warn(c.component, msg)
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.logglyClient.Warn(c.component, msg)
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

		if c.logglyClient != nil {
			if c.logglySampleRate == 100 {
				c.logglyClient.Error(c.component, msg)
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.logglyClient.Error(c.component, msg)
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
		if c.logglyClient != nil {
			c.logglyClient.Emergency(c.component, msg)
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
