package main

import (
	//"flag"
	//"fmt"
	"os"

	"github.com/ocadaruma/go-javactl/option"
	"github.com/ocadaruma/go-javactl/setting"
)

var AppVersion string

func main() {
	opts, err := option.ParseArgs(AppVersion, os.Args[1:])

	if err != nil { return }

	var s *setting.Setting
	s, err = setting.LoadSetting(opts.ConfigPath)

	if err != nil { return }

}
