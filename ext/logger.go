package ext

type LoggerInterface interface {
	Write(debug string, args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}
