package logler

import (
	"encoding/json"
	"errors"
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
	MinLog           bool
}

type Options struct {
	Component        string
	LogglySampleRate int
	MinimalLog       bool
}

type BQSchema struct {
	Bq         string `json:"bq"`
	Component  string `json:"component"`
	Sid        string `json:"sid"`
	AppVersion string `json:"appversion"`
	Category   string `json:"category"`
	Label      string `json:"label"`
	Label1     string `json:"label1"`
	Label2     string `json:"label2"`
	Label3     string `json:"label3"`
	Label4     string `json:"label4"`
	Label5     string `json:"label5"`
	Label6     string `json:"label6"`
	Label7     string `json:"label7"`
	Label8     string `json:"label8"`
	Label9     string `json:"label9"`
	Label10    string `json:"label10"`
	Action     string `json:"action"`
	ClientIP   string `json:"clientip"`
	Ua         string `json:"ua"`
	GeoIP      string `json:"geoip"`
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
		result.MinLog = opts.MinimalLog
		result.logglySampleRate = opts.LogglySampleRate
		if len(opts.Component) > 0 {
			result.component = opts.Component
		}
	}
	return result
}

//Sends data to google bigquery (only if json has bq=true)
// The data to bigquery should have the BQScheme struct
func (c *Client) BQ(bqs map[string]string) error {
	var msg []byte
	var err error
	if msg, err = json.Marshal(bqs); err != nil {
		log.Println(err.Error())
		return err
	}
	message := string(msg)
	c.info.Println(message)
	c.syslogClient.Info(message)
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
