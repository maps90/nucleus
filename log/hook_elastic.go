package log

import (
	"fmt"
	"os"
	"time"

	"github.com/maps90/nucleus/storage/elastic"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
)

const (
	errFormat    = "(%v)\n"
	loggerFormat = "account-2006.01.02"
)

// AddElasticLogHook required appName
func AddElasticLogHook(appName string) {
	client, err := elastic.New().Connect("logger")
	if err != nil {
		fmt.Printf(errFormat, "Unable to connect with elastic log server")
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf(errFormat, err)
		return
	}

	format := func() string {
		t := time.Now()
		ft := fmt.Sprintf("%s-2006.01.02", appName)
		return fmt.Sprintf(t.Format(ft))
	}

	hook, err := elogrus.NewElasticHookWithFunc(client, hostname, logrus.InfoLevel, format)
	if err != nil {
		fmt.Printf(errFormat, "Unable to create elastic log hook")
		return
	}

	Logger.Hooks.Add(hook)
}
