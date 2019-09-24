package ext

type Task struct {
	Key   string
	Value string
	Args  string
}

type QueueInterface interface {
	Push(key string, value string, args string) error
	Pop() *Task
}
