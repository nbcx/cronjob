package main

import (
	"database/sql"
	"fmt"
	"github.com/nbcx/cronjob/ext"
	"sync"
)

func init() {
	fmt.Println("hello plugin")
}

type redisQueue struct {
	db   sql.DB
	lock sync.Mutex
	log  ext.LoggerInterface
}

func NewQueue(conf ext.ConfigInterface, log ext.LoggerInterface) ext.QueueInterface {
	log.Error(conf.GetSectionString("save", "paths", "没有值啊"))
	return &redisQueue{
		log: log,
	}
}

func (this *redisQueue) Push(cmd string, args string) error {
	fmt.Println("Push?")
	return nil
}

func (this *redisQueue) Pop() *ext.Task {
	fmt.Println("Pop?")
	task := &ext.Task{}
	return task
}
