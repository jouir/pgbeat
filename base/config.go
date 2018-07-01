package base

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// AppName exposes application name to config module
var AppName string

// Config receives configuration options
type Config struct {
	Host             string  `yaml:"host"`
	Port             int     `yaml:"port"`
	User             string  `yaml:"user"`
	Password         string  `yaml:"password"`
	Database         string  `yaml:"database"`
	Schema           string  `yaml:"schema"`
	Table            string  `yaml:"table"`
	Interval         float64 `yaml:"interval"`
	Timeout          int     `yaml:"timeout"`
	ID               int     `yaml:"id"`
	RecoveryInterval float64 `yaml:"recovery-interval"`
	CreateDatabase   bool    `yaml:"create-database"`
	ConnectDatabase  string  `yaml:"connect-database"`
	CreateTable      bool    `yaml:"create-table"`
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
	return c.DsnWithDatabase(c.Database)
}

// DsnWithDatabase formats a connection string based on Config and overrides
// dbname
func (c *Config) DsnWithDatabase(database string) string {
	dsn := strings.Split(c.DsnWithoutDatabase(), " ")
	if c.Database != "" {
		dsn = append(dsn, fmt.Sprintf("dbname=%s", database))
	}
	return strings.Join(dsn, " ")
}

// DsnWithoutDatabase formats a connection string based on Config without
// dbname
func (c *Config) DsnWithoutDatabase() string {
	var dsn []string
	if c.Host != "" {
		dsn = append(dsn, fmt.Sprintf("host=%s", c.Host))
	}
	if c.Port != 0 {
		dsn = append(dsn, fmt.Sprintf("port=%d", c.Port))
	}
	if c.User != "" {
		dsn = append(dsn, fmt.Sprintf("user=%s", c.User))
	}
	if c.Password != "" {
		dsn = append(dsn, fmt.Sprintf("password=%s", c.Password))
	}
	if c.Timeout != 0 {
		dsn = append(dsn, fmt.Sprintf("connect_timeout=%d", c.Timeout))
	}
	if AppName != "" {
		dsn = append(dsn, fmt.Sprintf("application_name=%s", AppName))
	}
	return strings.Join(dsn, " ")
}
