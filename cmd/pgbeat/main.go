package main

import (
	"flag"
	"fmt"
	"github.com/jouir/pgbeat/base"
	"github.com/jouir/pgbeat/manager"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// AppVersion stores application version at compilation time
var AppVersion string

func main() {
	var err error
	config := base.NewConfig()

	version := flag.Bool("version", false, "Print version")
	configFile := flag.String("config", "", "Configuration file")
	flag.StringVar(&config.Host, "host", "", "Instance host address")
	flag.IntVar(&config.Port, "port", 0, "Instance port")
	flag.StringVar(&config.User, "user", "", "Instance username")
	flag.StringVar(&config.Password, "password", "", "Instance password")
	prompt := flag.Bool("prompt-password", false, "Prompt for password")
	flag.StringVar(&config.Database, "database", "", "Database name")
	flag.StringVar(&config.Schema, "schema", "public", "Schema name")
	flag.StringVar(&config.Table, "table", "pgbeat", "Table name")
	flag.Float64Var(&config.Interval, "interval", 1, "Time to sleep between updates in seconds")
	flag.IntVar(&config.Timeout, "timeout", 3, "Connection timeout in seconds")
	flag.IntVar(&config.ID, "id", 1, "Differenciate daemons by using an indentifier")
	flag.Parse()

	if *version {
		if AppVersion == "" {
			AppVersion = "unknown"
		}
		fmt.Println(AppVersion)
		return
	}

	if *prompt {
		fmt.Print("Password: ")
		bytes, err := terminal.ReadPassword(syscall.Stdin)
		base.Panic(err)
		config.Password = string(bytes)
		fmt.Print("\n")
	}

	if *configFile != "" {
		err = config.Read(*configFile)
		base.Panic(err)
	}

	beatmaker := manager.NewBeatmaker(config)

	// Signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			log.Printf("Received %v signal\n", sig)
			beatmaker.Terminate()
			os.Exit(0)
		}
	}()

	beatmaker.Run()
}
