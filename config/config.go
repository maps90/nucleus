package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

var v *viper.Viper

func init() {
	v = viper.New()
}

// ErrMsgTmpl is an error template
const ErrMsgTmpl = "something went wrong. Error: %s"
const defaultBaseConfig = "/resources"

func NewEnv(prefix string) {
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()
}

// New configuration initialization, it requires configName and baseConfigPath
// by default configName is "config", and the base config is the root application folder.
func New(args ...string) (err error) {
	var configName = "application"
	var baseConfig string

	if len(args) > 0 {
		configName = args[0]
	}

	if len(args) == 2 && args[1] != "" {
		baseConfig = args[1]
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	if baseConfig == "" {
		baseConfig = path + defaultBaseConfig
	}

	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(baseConfig)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf(ErrMsgTmpl, err)
	}

	return
}

// Set sets the value for the key in the override regiser.
// Set is case-insensitive for a key.
// Will be used instead of values obtained via
// flags, config file, ENV, default, or key/value store.
func Set(key string, value interface{}) { v.Set(key, value) }

// Get can retrieve any value given the key to use.
// Get is case-insensitive for a key.
// Get has the behavior of returning the value associated with the first
// place from where it is set. Viper will check in the following order:
// override, flag, env, config file, key/value store, default
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func Get(key string) interface{} { return v.Get(key) }

// Sub returns new Viper instance representing a sub tree of this instance.
// Sub is case-insensitive for a key.
func Sub(key string) *viper.Viper { return v.Sub(key) }

// GetString returns the value associated with the key as a string.
func GetString(key string) string { return v.GetString(key) }

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool { return v.GetBool(key) }

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int { return v.GetInt(key) }

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(key string) int64 { return v.GetInt64(key) }

// GetFloat64 returns the value associated with the key as a float.
func GetFloat64(key string) float64 { return v.GetFloat64(key) }

// GetTime returns the value associated with the key as time.
func GetTime(key string) time.Time { return v.GetTime(key) }

// GetStringSlice returns the value associated with the key as a string slice.
func GetStringSlice(key string) []string { return v.GetStringSlice(key) }

// GetDuration returns the value associated with the key as a time.duration.
func GetDuration(key string) time.Duration { return v.GetDuration(key) }

// GetStringMap returns the value associated with the key as a string map.
func GetStringMap(key string) map[string]interface{} { return v.GetStringMap(key) }

// GetStringMapStringSlice returns the value associated with the key as a string map string slice.
func GetStringMapStringSlice(key string) map[string][]string { return v.GetStringMapStringSlice(key) }
