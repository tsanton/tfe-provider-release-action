package utilities

type ILogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}
