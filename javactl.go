package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/ocadaruma/go-javactl/executor"
	"github.com/ocadaruma/go-javactl/logger"
	"github.com/ocadaruma/go-javactl/option"
	"github.com/ocadaruma/go-javactl/setting"
	"github.com/ocadaruma/go-javactl/setting/mapping"
)

var AppVersion string

func main() {
	var err error

	var opts *option.JavactlOpts
	opts, err = option.ParseArgs(AppVersion, os.Args[1:])
	if err != nil {
		return
	}

	defer catchUnexpected(opts.Debug)

	var config *mapping.YAMLConfig
	config, err = mapping.LoadConfig(opts.ConfigPath)
	if err != nil {
		printErrorThenExit(err)
	}

	var sett *setting.Setting
	sett, err = setting.NewSetting(config)
	if err != nil {
		printErrorThenExit(err)
	}

	executor := executor.NewExecutor(logger.NewSystemLogger(opts.DryRun), sett, opts)
	now := time.Now()

	err = executor.CheckRequirement()
	if err != nil {
		printErrorThenExit(err)
	}

	err = executor.CheckDuplicateProcess()
	if err != nil {
		os.Exit(0)
	} else {
		err = executor.CreateDirectories()
		if err != nil {
			printErrorThenExit(err)
		}

		err = executor.CleanOldLogs(now)
		if err != nil {
			printErrorThenExit(err)
		}

		err = executor.Execute(now)
		if err != nil {
			os.Exit(1)
		}
	}
}

func catchUnexpected(debugEnabled bool) {
	err := recover()
	if err != nil {
		fmt.Printf("unexpected error occurred: %v\n", err)
		if debugEnabled {
			debug.PrintStack()
		}
		os.Exit(2)
	}
}

func printErrorThenExit(err error) {
	fmt.Printf("an error occurred: %v\n", err)
	os.Exit(1)
}
