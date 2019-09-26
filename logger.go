package main

import (
	"fmt"
	"github.com/nbcx/cronjob/ext"
	"io"
	"log"
	"os"
)

type loger struct {
	hander *log.Logger
}

func NewLogger(conf *Config) ext.LoggerInterface {
	var out io.Writer
	if conf.ip == "sd" {
		file, err := os.OpenFile("gox.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("%s\r\n", err.Error())
			os.Exit(-1)
		}
		out = file
	} else {
		out = os.Stdout
	}
	return &loger{
		hander: log.New(out, "\r\n", log.Ldate|log.Ltime|log.Llongfile),
	}
}

func (log *loger) Write(debug string, args ...interface{}) {
	log.hander.SetPrefix(debug)
	log.hander.Output(3, fmt.Sprint(args...))
}

func (log *loger) Info(args ...interface{}) {
	log.Write("INFO ", args...)
}

func (log *loger) Warning(args ...interface{}) {
	log.Write("WARNING ", args...)
}

func (log *loger) Error(args ...interface{}) {
	log.Write("ERROR ", args...)
}

func (log *loger) Debug(args ...interface{}) {
	log.Write("DEBUG ", args...)
}

func (log *loger) Fatal(args ...interface{}) {
	log.Write("FATAL ", args...)
}
