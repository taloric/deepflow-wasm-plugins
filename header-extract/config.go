package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deepflowio/deepflow-wasm-go-sdk/sdk"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PortWhiteList []string `yaml:"port-white-list"`
	ProcName      []string `yaml:"proc-name"`

	ports []ports
}

type ports struct {
	min   uint16
	max   uint16
	point uint16
}

func (c *Config) init(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(file) > 0 {
		yaml.Unmarshal(file, c)
	}

	if len(c.PortWhiteList) > 0 {
		for _, v := range c.PortWhiteList {
			c.ports = append(c.ports, c.parsePort(v))
		}
	}
	return nil
}

func (c *Config) parsePort(port_range string) ports {
	if p := strings.Split(port_range, "-"); len(p) > 1 {
		min, err := strconv.Atoi(p[0])
		if err != nil {
			sdk.Warn(fmt.Sprintf("%s: %s invalid port: %s", plugin_module, port_range, p[0]))
		}
		max, err := strconv.Atoi(p[1])
		if err != nil {
			sdk.Warn(fmt.Sprintf("%s: %s invalid port: %s", plugin_module, port_range, p[1]))
		}
		return ports{min: uint16(min), max: uint16(max)}
	}
	point, err := strconv.Atoi(port_range)
	if err != nil {
		sdk.Warn(fmt.Sprintf("%s: invalid port_range: %s", plugin_module, port_range))
		return ports{}
	}
	return ports{point: uint16(point)}
}

func (c *Config) allowCapturePort(port uint16) bool {
	if len(c.ports) == 0 {
		return true
	}
	for _, v := range c.ports {
		if v.min*v.max > 0 && port >= v.min && port <= v.max {
			return true
		}
		if v.point == port {
			return true
		}
	}
	return false
}
