package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/nbcx/cronjob/ext"
	"plugin"
)

func NewQueue(conf *Config) (ext.QueueInterface, error) {

	queue := conf.GetString("queue")

	if queue != "" {
		has, err := ext.FileExists(queue)

		if !has {
			return nil, err
		}

		//打开动态库
		p, err := plugin.Open(queue)
		if err != nil {
			panic(err)
		}

		//接口验证
		f, err := p.Lookup("NewQueue")
		if err != nil {
			panic(err)
		}

		d, err := f.(func(conf ext.ConfigInterface, log ext.LoggerInterface) (ext.QueueInterface, error))(conf, logger)
		return d, err
	}

	return nil, nil
}
