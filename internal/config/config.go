package config

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

type LoggerConf struct {
	Level string
}

type OS string

const (
	OSWindows OS = "windows"
	OSLinux   OS = "linux"
)

type StatsParamsConf struct {
	OS          OS
	M           int  `yaml:"m"`
	N           int  `yaml:"n"`
	CPU         bool `yaml:"cpu"`
	DisksUsage  bool `yaml:"disksUsage"`
	DisksIoStat bool `yaml:"disksIoStat"`
	NetStat     bool `yaml:"netStat"`
}

type Config struct {
	Logger      LoggerConf
	StatsParams StatsParamsConf `yaml:"statsParams"`
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) Read(fpath string) (err error) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	switch os := runtime.GOOS; os {
	case "linux":
		c.StatsParams.OS = OSLinux
	case "windows":
		c.StatsParams.OS = OSWindows
	default:
		panic(fmt.Sprintf("program not working on %s", os))
	}
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("%+v", *c)
}
