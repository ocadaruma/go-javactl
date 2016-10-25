package setting

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Setting struct {
	App App
	Java Java
	Log *Log
	OS OS
	PreCommands  []string `yaml:"pre"`
	PostCommands []string `yaml:"post"`
}

func LoadSetting(configPath string) (result *Setting, err error) {
	result, err = load(configPath)

	if err != nil { return }

	app, err := result.App.Normalize()
	if err != nil { return }

	java, err := result.Java.Normalize()
	if err != nil { return }

	result.App = *app
	result.Java = *java

	return
}

func load(configPath string) (result *Setting, err error) {
	var buf []byte
	buf, err = ioutil.ReadFile(configPath)

	if err != nil { return }

	result = &Setting{}
	err = yaml.Unmarshal(buf, result)

	return
}
