package main

import (
	"encoding/json"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

const pluginName = "github.com/chronotc/monorepo-diff"

// Plugin buildkite monorepo diff plugin structure
type Plugin struct {
	Diff          string
	Wait          bool
	LogLevel      string
	Interpolation bool
	Hooks         []struct{ Command string }
	Watch         []struct {
		Path   string
		Config struct {
			Trigger string
		}
		Label string
		Build struct {
			Message string
			Branch  string
			Commit  string
			Env     map[string]string
		}
		Command string
		Async   bool
		Agents  struct {
			Queue string
		}
		Env map[string]string
	}
}

// UnmarshalJSON set defaults properties
func (s *Plugin) UnmarshalJSON(data []byte) error {
	type plain Plugin
	test := &plain{
		Diff:          "git diff --name-only HEAD~1",
		Wait:          false,
		LogLevel:      "info",
		Interpolation: false,
	}

	_ = json.Unmarshal(data, test)

	*s = Plugin(*test)
	return nil
}

func initializePlugin(data string) (Plugin, error) {
	var plugins []map[string]Plugin

	err := json.Unmarshal([]byte(data), &plugins)

	if err != nil {
		log.Debug(err)
		return Plugin{}, errors.New("Failed to parse plugin configuration")
	}

	for _, p := range plugins {
		for key, plugin := range p {
			if strings.HasPrefix(key, pluginName) {
				return plugin, nil
			}
		}
	}

	return Plugin{}, errors.New("Could not initialize plugin")
}
