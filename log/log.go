package log

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger implements logrus logger
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.JSONFormatter{}
	return
}

// New logger
// Usage: this is only an example what this function can do
// New(DBTopic).Error("db connection problem..")
func New() *Entry {
	_, _, fName := caller()
	l := map[string]interface{}{
		"context": fName,
	}
	return WithFields(l)
}

// WithFields logger
// Usage: this is only an example what this function can do
// WithFields(IOTopic, "Error:", errors.New("got a 200 http status but got empty response")).Info("got an error...")
func WithFields(messages map[string]interface{}) *Entry {
	_, _, fName := caller()
	l := map[string]interface{}{
		"context": fName,
	}

	for i, v := range messages {
		l[i] = v
	}

	return &Entry{Entry: Logger.WithFields(l)}
}

func caller() (file string, line int, name string) {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)

	// get the info of the actual function that's in the pointer
	f := runtime.FuncForPC(pc[0] - 1)
	file, line = f.FileLine(pc[0])

	fn := strings.Split(f.Name(), "/")
	fName := fn[len(fn)-1]

	return file, line, fName
}
