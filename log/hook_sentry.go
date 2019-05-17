package log

import (
	"errors"
	"fmt"
	"os"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// AddSentryHook logrus
func AddSentryHook() {
	if !cast.ToBool(os.Getenv("sentry.enabled")) {
		return
	}

	dsn := os.Getenv("sentry.dsn")
	hook, err := logrus_sentry.NewSentryHook(dsn, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})
	if err != nil {
		fmt.Printf(errFormat, errors.New("Unable to create sentry log hook"))
		return
	}

	Logger.Hooks.Add(hook)
}
