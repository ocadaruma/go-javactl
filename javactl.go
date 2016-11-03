package main

import (
	//"flag"
	//"fmt"
	"os"

	"github.com/ocadaruma/go-javactl/executor"
	"github.com/ocadaruma/go-javactl/option"
	"github.com/ocadaruma/go-javactl/setting"
	"github.com/ocadaruma/go-javactl/setting/mapping"
	"github.com/ocadaruma/go-javactl/logger"
	"time"
)

var AppVersion string

func main() {
	var err error

	var opts *option.JavactlOpts
	opts, err = option.ParseArgs(AppVersion, os.Args[1:])
	if err != nil { return }

	var config *mapping.YAMLConfig
	config, err = mapping.LoadConfig(opts.ConfigPath)
	if err != nil { return }

	var sett *setting.Setting
	sett, err = setting.NewSetting(config)
	if err != nil { return }

	executor := executor.NewExecutor(logger.NewSystemLogger(opts.DryRun), sett, opts)
	now := time.Now()

	err = executor.CheckRequirement()
	if err != nil { return }

	err = executor.CreateDirectories()
	if err != nil { return }

	err = executor.CleanOldLogs(now)
	if err != nil { return }

	err = executor.Execute(now)
	if err != nil { return }
}
