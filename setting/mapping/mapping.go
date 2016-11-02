package mapping

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type YAMLConfig struct {
	App App
	Java Java
	Log *Log
	OS OS
	PreCommands  []string `yaml:"pre"`
	PostCommands []string `yaml:"post"`
}

func LoadConfig(configPath string) (result *YAMLConfig, err error) {
	var buf []byte
	buf, err = ioutil.ReadFile(configPath)

	if err != nil { return }

	result = &YAMLConfig{}
	err = yaml.Unmarshal(buf, result)

	return
}

// represents app entries
type App struct {
	Name string
	Home string
	Jar string
	EntryPoint string `yaml:"entry_point"`
	Command string
	PidFile string `yaml:"pid_file"`
}

// represents java entries
type Java struct {
	Home string
	Version float32
	Server bool
	Memory *Memory
	JMX *JMX
	Prop map[string]string
	Option []string
}

type JMX struct {
	Port *int
	SSL *bool
	Authenticate *bool
}

type Memory struct {
	HeapMin string `yaml:"heap_min"`
	HeapMax string `yaml:"heap_max"`
	PermMin string `yaml:"perm_min"`
	PermMax string `yaml:"perm_max"`
	MetaspaceMin string `yaml:"metaspace_min"`
	MetaspaceMax string `yaml:"metaspace_max"`
	NewMin string `yaml:"new_min"`
	NewMax string `yaml:"new_max"`
	SurvivorRatio *int `yaml:"survivor_ratio"`
	TargetSurvivorRatio *int `yaml:"target_survivor_ratio"`
}

// represents log entries
type Log struct {
	ConsoleLog *ConsoleLog `yaml:"console"`
	GCLog *GCLog `yaml:"gc"`
	Dump *Dump
	ErrorLog *ErrorLog `yaml:"error"`
}

type ConsoleLog struct {
	Prefix string
	MaxSize string `yaml:"max_size"`
	Backup int
	Preserve int
}

type GCLog struct {
	Prefix string
	MaxSize string `yaml:"max_size"`
	Backup int
	Preserve int
}

type Dump struct {
	Prefix string
}

type ErrorLog struct {
	Path string
}

// represents os entries
type OS struct {
	User string
	Env map[string]string
}
