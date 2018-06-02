package base

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// AppName exposes application name to config module
var AppName string

// Config receives configuration options
type Config struct {
	File     string
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
	Table    string `yaml:"table"`
	Interval int    `yaml:"interval"`
	Timeout  int    `yaml:"timeout"`
	ID       int    `yaml:"id"`
}

func init() {
	AppName = "pgbeat"
}

// NewConfig creates a Config object
func NewConfig() *Config {
	return &Config{}
}

// Read loads options from a configuration file to Config
func (c *Config) Read(file string) error {
	file, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return err
	}

	return nil
}

// Dsn formats a connection string based on Config
func (c *Config) Dsn() string {
	var dsn string
	if c.Host != "" {
		dsn += fmt.Sprintf("host=%s ", c.Host)
	}
	if c.Port != 0 {
		dsn += fmt.Sprintf("port=%d ", c.Port)
	}
	if c.User != "" {
		dsn += fmt.Sprintf("user=%s ", c.User)
	}
	if c.Password != "" {
		dsn += fmt.Sprintf("password=%s ", c.Password)
	}
	if c.Database != "" {
		dsn += fmt.Sprintf("dbname=%s ", c.Database)
	}
	if c.Timeout != 0 {
		dsn += fmt.Sprintf("connect_timeout=%d ", c.Timeout)
	}
	if AppName != "" {
		dsn += fmt.Sprintf("application_name=%s", AppName)
	}
	return dsn
}
