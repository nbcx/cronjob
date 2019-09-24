package main

import (
	"fmt"
	"os"
	"github.com/nbcx/cronjob/ext"
)

func main() {
	fmt.Printf("process %d\n", os.Getpid())
	ejob := NewEJobApplication()
	ejob.run()
}

type EJob struct {
	conf       *Config
	commandMap map[string]func(string) /*创建指令集合 */
	server     ext.ServerInterface
}

//获取EJob
func NewEJobApplication() *EJob {
	conf := NewConfigObject()
	ejob := &EJob{
		conf:       conf,
		commandMap: make(map[string]func(string)),
		server:     NewServer(conf),
	}
	//ejob.handleCommand("f", f)
	//ejob.handleCommand("g", g)
	return ejob
}

//注册指令
func (ejob *EJob) handleCommand(pattern string, handler func(string)) {
	ejob.commandMap[pattern] = handler
}

//启动EJob
func (ejob *EJob) run() {

	//ejob.commandMap["g"]("hello")

	//刷新
	//-s
	ejob.server.Run()
}
