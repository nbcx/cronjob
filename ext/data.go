package ext

type Task struct {
	Id   string
	Cmd  string
	Args string
	Ct   string
}

type QueueInterface interface {
	Push(cmd string, args string) error
	Pop() *Task
}
