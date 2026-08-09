package main

import (
	"flag"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func envLookup(env map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, found := env[key]
		return value, found
	}
}

func TestLoadConfig(t *testing.T) {
	tests := map[string]struct {
		args     []string
		env      map[string]string
		expected Config
	}{
		"defaults": {
			expected: NewConfig(),
		},
		"flags": {
			args: []string{
				"-remote", "flag.example:10011",
				"-user", "flaguser",
				"-password", "flagpassword",
				"-listen", "127.0.0.1:1234",
				"-enablechannelmetrics",
				"-ignorefloodlimits",
			},
			expected: Config{
				Remote:               "flag.example:10011",
				ListenAddr:           "127.0.0.1:1234",
				User:                 "flaguser",
				Password:             "flagpassword",
				EnableChannelMetrics: true,
				IgnoreFloodLimits:    true,
			},
		},
		"environment": {
			env: map[string]string{
				"REMOTE":                 "env.example:10011",
				"SERVERQUERY_USER":       "envuser",
				"SERVERQUERY_PASSWORD":   "envpassword",
				"LISTEN_ADDRESS":         "127.0.0.1:4321",
				"ENABLE_CHANNEL_METRICS": "1",
				"IGNORE_FLOOD_LIMITS":    "true",
			},
			expected: Config{
				Remote:               "env.example:10011",
				ListenAddr:           "127.0.0.1:4321",
				User:                 "envuser",
				Password:             "envpassword",
				EnableChannelMetrics: true,
				IgnoreFloodLimits:    true,
			},
		},
		"flags take precedence over the environment": {
			args: []string{"-remote", "flag.example:10011", "-enablechannelmetrics=false"},
			env: map[string]string{
				"REMOTE":                 "env.example:10011",
				"SERVERQUERY_USER":       "envuser",
				"ENABLE_CHANNEL_METRICS": "true",
			},
			expected: Config{
				Remote:               "flag.example:10011",
				ListenAddr:           "0.0.0.0:9189",
				User:                 "envuser",
				EnableChannelMetrics: false,
			},
		},
		"invalid boolean environment value falls back to the default": {
			env: map[string]string{"ENABLE_CHANNEL_METRICS": "definitely"},
			expected: Config{
				Remote:               "localhost:10011",
				ListenAddr:           "0.0.0.0:9189",
				User:                 "serveradmin",
				EnableChannelMetrics: false,
			},
		},
		"invalid boolean environment value is still overridden by the flag": {
			args: []string{"-enablechannelmetrics"},
			env:  map[string]string{"ENABLE_CHANNEL_METRICS": "definitely"},
			expected: Config{
				Remote:               "localhost:10011",
				ListenAddr:           "0.0.0.0:9189",
				User:                 "serveradmin",
				EnableChannelMetrics: true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fs := flag.NewFlagSet(name, flag.ContinueOnError)
			fs.SetOutput(io.Discard)

			config, err := loadConfig(fs, tc.args, envLookup(tc.env))
			require.Nil(t, err)
			assert.Equal(t, tc.expected, config)
		})
	}
}

func TestLoadConfigParseError(t *testing.T) {
	fs := flag.NewFlagSet("parse error", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	_, err := loadConfig(fs, []string{"-unknown"}, envLookup(nil))
	assert.NotNil(t, err)
}
