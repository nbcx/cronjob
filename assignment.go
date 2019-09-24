package main

import (
	"github.com/nbcx/cronjob/ext"
	"os"
	"os/exec"
	"strings"
)

//https://blog.csdn.net/dodod2012/article/details/82774257
type AssignmentInterface interface {
	run(cmd string, args string)
}

type assignment struct {
	queue  ext.QueueInterface
	conf   *Config
	flag   bool
	count  int64
	Ccount int64
}

func NewAssignment(conf *Config) *assignment {
	queue, _ := NewQueue(conf)

	assignment := &assignment{
		queue: queue,
		conf:  conf,
		flag:  false,
		count: 0,
	}
	return assignment
}

func (this *assignment) script(key string, value string, args string) {

}

func (this *assignment) run(key string, value string, args string) {
	err := this.queue.Push(key, value, args)

	if err != nil {
		logger.Error("this.queue.push error: ", err)
		return
	}

	//有协程在执行出队操作，则返回
	if this.flag {
		return
	}

	//TODO 是否需要用锁，待考虑
	this.flag = true

	for {
		data := this.queue.Pop()
		if data == nil {
			logger.Info("no tasks to be performed")
			break
		}

		//通过ID执行
		if data.Key == "id" {
			cmds, ok := this.conf.jobs[data.Value]
			if !ok {
				logger.Info("job ", data.Value, " is not support")
				continue
			}
			command := strings.Join(cmds[3:], " ") + " " + data.Args

			go this.shell(command)

			//原生调用
		} else if data.Key == "pro" {
			go this.command(data.Value + " " + data.Args)

			//通过bash调用
		} else {
			go this.shell(data.Value + " " + data.Args)
			break
		}
	}

	this.count++
	this.flag = false
	logger.Info("go goroutines number ", this.count)
}

//以shell运行指令
func (this *assignment) shell(command string) {

	logger.Info("shell: ", command)

	// TODO 跨平台待支持
	cmd := exec.Command("/bin/bash", "-c", command)
	_, err := cmd.Output()

	if err != nil {
		logger.Error("execute shell:%s failed with error:%s", command, err.Error())
	}
}

//原生指令
func (this *assignment) command(command string) error {
	// Split the input separate the command and the arguments.
	args := strings.Split(command, " ")

	// Check for built-in commands.
	switch args[0] {
	case "cd":
		// 'cd' to home with empty path not yet supported.
		if len(args) < 2 {
			return nil //todo ErrNoPath
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	}

	// Prepare the command to execute.
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command return the error.
	return cmd.Run()
}
