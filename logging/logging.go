package logging

import (
	"os"

	gologging "github.com/op/go-logging"
)

var (
	Log            Logger
	leveledBackend gologging.LeveledBackend
)

type Logger struct {
	log *gologging.Logger
}

func init() {
	Log = Logger{
		log: gologging.MustGetLogger("default"),
	}

	format := gologging.MustStringFormatter(`%{shortfile} %{color} %{time:2006-01-02T15:04:05.000000Z} %{level:.4s} %{color:reset} %{message}`)
	backend := gologging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := gologging.NewBackendFormatter(backend, format)
	leveledBackend = gologging.AddModuleLevel(backendFormatter)
	gologging.SetBackend(leveledBackend)
	SetLogLevel(SILENT)
}

func SetLogLevel(l LogLevel) {
	leveledBackend.SetLevel(convertLogLevel(l), "")
}

func (l Logger) Error(msg string, err error) {
	// NOTICE: go vet doesn't like the following line:
	//   l.log.Error("%s : %v", msg, err)
	// Therefore:
	l.log.Error(""+"%s : %v", msg, err)
}

func (l Logger) Errorf(msg string, values ...interface{}) {
	l.log.Errorf(msg, values...)
}

func (l Logger) Info(msg string) {
	l.log.Info(msg)
}

func (l Logger) Infof(msg string, values ...interface{}) {
	l.log.Infof(msg, values...)
}

func (l Logger) Debug(msg string) {
	l.log.Debug(msg)
}

func (l Logger) Debugf(msg string, values ...interface{}) {
	l.log.Debugf(msg, values...)
}

func (l Logger) Panic(msg string, err error) {
	l.log.Panicf("%s : %v", msg, err)
}

func (l Logger) Panicf(msg string, values ...interface{}) {
	l.log.Panicf(msg, values...)
}

func convertLogLevel(l LogLevel) gologging.Level {
	switch l {
	case SILENT:
		return -1
	case FATAL:
		return gologging.CRITICAL
	case ERROR:
		return gologging.ERROR
	case INFO:
		return gologging.INFO
	case DEBUG:
		return gologging.DEBUG
	default:
		Log.Panicf("Unknown LogLevel: %v", l)
		return gologging.CRITICAL
	}
}
