package setting

import (
	"io/ioutil"
	"time"

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

func LoadConfig(configPath string) (result *Setting, err error) {
	result, err = load(configPath)

	if err != nil { return }

	err = result.App.Normalize()
	if err != nil { return }

	err = result.Java.Normalize()

	return
}

func (this Setting) GetArgs(now time.Time) []string {
	result := []string{}

	logArgs := []string{}
	if this.Log != nil { logArgs = this.Log.GetOpts() }

	this.App.GetArgs()

	return result
}

func load(configPath string) (result *Setting, err error) {
	var buf []byte
	buf, err = ioutil.ReadFile(configPath)

	if err != nil { return }

	result = &Setting{}
	err = yaml.Unmarshal(buf, result)

	return
}
