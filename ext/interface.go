package ext

//执行器接口
type AssignmentInterface interface {
	Script(key string, value string, args string)
	Queue(key string, value string, args string)
}

//队列接口
type QueueInterface interface {
	Push(key string, value string, args string) error
	Pop() *Task
}

//设置接口
type ConfigInterface interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	GetString(key string, defaultValue ...string) string
	GetSectionString(section string, key string, defaultValue ...string) string
}

//日志接口
type LoggerInterface interface {
	Write(debug string, args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}

//服务接口
type ServerInterface interface {
	Run()
}
