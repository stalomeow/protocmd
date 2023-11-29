package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type cmdIdConfig struct {
	args     map[string]string
	cmdIdMap map[string]uint16
}

func loadCmdIdConfig(params string) (*cmdIdConfig, error) {
	buf, err := os.ReadFile("cmdid.yaml")
	if err != nil {
		return nil, err
	}

	cmdId := make(map[string]uint16)
	err = yaml.Unmarshal(buf, &cmdId)
	if err != nil {
		return nil, err
	}

	args := make(map[string]string)
	for _, param := range strings.Split(params, ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}

		if params != "" {
			args[param] = value
		}
	}

	return &cmdIdConfig{
		args:     args,
		cmdIdMap: cmdId,
	}, nil
}

func (config *cmdIdConfig) writeFlags(flags *flag.FlagSet) error {
	for k, v := range config.args {
		err := flags.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
