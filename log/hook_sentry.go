package log

import (
	"errors"
	"fmt"

	"github.com/evalphobia/logrus_sentry"
	"github.com/maps90/nucleus/config"
	"github.com/sirupsen/logrus"
)

// AddSentryHook logrus
func AddSentryHook() {
	if !config.GetBool("sentry.enabled") {
		return
	}

	dsn := config.GetString("sentry.dsn")
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
