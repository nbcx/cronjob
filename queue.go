package main

import (
	"github.com/nbcx/cronjob/ext"
)

func NewQueue(conf *Config) (ext.QueueInterface, error) {

	//获取队列插件
	queue := conf.GetString("queue")

	if queue == "" {
		logger.Info("无队列插件")
		return nil, nil
	}

	f, err := loadOn(queue)
	if err != nil {
		return nil, err
	}

	d, err := f.(func(conf ext.ConfigInterface, log ext.LoggerInterface) (ext.QueueInterface, error))(conf, logger)
	return d, err
}
