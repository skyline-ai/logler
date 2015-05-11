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

	"google.golang.org/api/bigquery/v2"
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
	MinLog           bool
}

type Options struct {
	Component        string
	LogglySampleRate int
	MinimalLog       bool
}

var (
	bqSchemaMap = map[string]bool{
		"bq":         true,
		"component":  true,
		"appversion": true,
		"timestamp":  true,
		"category":   true,
		"label":      true,
		"label1":     true,
		"label2":     true,
		"label3":     true,
		"label4":     true,
		"label5":     true,
		"label6":     true,
		"label7":     true,
		"label8":     true,
		"label9":     true,
		"label10":    true,
		"action":     true,
		"clientip":   true,
		"ua":         true,
		"geoip":      true,
	}
)

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
		result.MinLog = opts.MinimalLog
		result.logglySampleRate = opts.LogglySampleRate
		if len(opts.Component) > 0 {
			result.component = opts.Component
		}
	}
	initBQ()

	return result
}

//Sends data to google bigquery (if fits bqSchemaMap)
func (c *Client) BQ(bqs map[string]interface{}) error {
	bqmap := make(map[string]bigquery.JsonValue)
	for key, value := range bqs {
		if _, ok := bqSchemaMap[key]; ok {
			bqmap[key] = fmt.Sprintf("%v", value)
		}
	}
	bqmap["timestamp"] = time.Now().UTC()
	sendBQ(bqmap)
	if message, err := json.Marshal(bqs); err != nil {
		errmsg := fmt.Sprintf("Could not marshal %v\n", bqs)
		return errors.New(errmsg)
	} else {
		c.info.Println(string(message))
	}
	return nil
}

func (c *Client) Info(msg map[string]interface{}) {
	if msg, err := c.getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		message := string(j)
		c.info.Println(message)

		if c.syslogClient != nil {
			if c.logglySampleRate == 100 {
				c.syslogClient.Info(message)
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.syslogClient.Info(message)
				}
			}
		}
	}
}

func (c *Client) Warn(msg map[string]interface{}) {
	if msg, err := c.getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		message := string(j)
		c.warn.Println(message)

		if c.syslogClient != nil {
			if c.logglySampleRate == 100 {
				c.syslogClient.Warning(message)
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.syslogClient.Warning(message)
				}
			}
		}
	}
}

func (c *Client) Error(msg map[string]interface{}) {
	if msg, err := c.getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		message := string(j)
		c.error.Println(message)
		if c.syslogClient != nil {
			if c.logglySampleRate == 100 {
				c.syslogClient.Err(message)
			} else {
				if random(1, 100) <= c.logglySampleRate {
					c.syslogClient.Err(message)
				}
			}
		}
	}
}

func (c *Client) Emergency(msg map[string]interface{}) {
	if msg, err := c.getMessage(msg); err != nil {
		log.Println(err.Error())
	} else {
		j, _ := json.Marshal(msg)
		c.emergency.Println(string(j))
		if c.syslogClient != nil {
			c.syslogClient.Emerg(string(j))
		}
	}
}

func (c *Client) getMessage(msg map[string]interface{}) (map[string]interface{}, error) {
	if msg != nil && !c.MinLog {
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[1])
		file, line := f.FileLine(pc[1])
		tmp := strings.Split(file, string(filepath.Separator))
		msg["filename"] = tmp[len(tmp)-1]
		msg["line"] = line
		msg["func"] = f.Name()
		return msg, nil
	} else if msg != nil {
		return msg, nil
	}
	return nil, errors.New("message log nil message")
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
