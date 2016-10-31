package main

import (
	//"flag"
	//"fmt"
	"os"

	"github.com/ocadaruma/go-javactl/executor"
	"github.com/ocadaruma/go-javactl/option"
	"github.com/ocadaruma/go-javactl/setting"
	"github.com/ocadaruma/go-javactl/setting/mapping"
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

	var s *setting.Setting
	s, err = setting.NewSetting(config)
	if err != nil { return }

	var executor = executor.NewExecutor()
}
