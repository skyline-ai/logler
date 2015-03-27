package logler

import (
	"encoding/json"
	"github.com/streamrail/go-loggly"
	"log"
	"os"
)

type Logler struct {
	info   *log.Logger
	warn   *log.Logger
	error  *log.Logger
	loggly *loggly.Client
}

type Options struct {
	LogglyToken string
}

type Message loggly.Message

func New(opts *Options) *Logler {
	result := &Logler{
		info: log.New(os.Stdout,
			"INFO: ",
			log.Ldate|log.Ltime),

		warn: log.New(os.Stdout,
			"WARNING: ",
			log.Ldate|log.Ltime),

		error: log.New(os.Stderr,
			"ERROR: ",
			log.Ldate|log.Ltime),
	}
	if opts != nil {
		result.loggly = loggly.New(opts.LogglyToken)
	}
	return result
}

func (l *Logler) Info(msg *Message) {
	j, _ := json.Marshal(msg)
	l.info.Println(string(j))
}
func (l *Logler) Warn(msg *Message) {
	j, _ := json.Marshal(msg)
	l.warn.Println(string(j))
}
func (l *Logler) Error(msg *Message) {
	j, _ := json.Marshal(msg)
	l.error.Println(string(j))
}