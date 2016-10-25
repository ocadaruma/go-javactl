package setting

import (
	"fmt"
	"path/filepath"
)

type Java struct {
	Home string
	Version float32
	Server bool
	Memory *Memory
	JMX *JMX
	Env map[string]string
	Option []string
}

func (this Java) Normalize() (result *Java, err error) {
	if this.Home == "" || !filepath.IsAbs(this.Home) {
		err = fmt.Errorf("java.home(%s) is required and must be an absolute path", this.Home)
		return
	}

	if this.Version < 1 {
		err = fmt.Errorf("invalid java.version(%f)", this.Version)
		return
	}

	memory, err := this.Memory.validate(this.Version)
	if err != nil { return }

	result = &this
	result.Memory = memory

	return
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
	SurvivorRatio int `yaml:"survivor_ratio"`
	TargetSurvivorRatio int `yaml:"target_survivor_ratio"`
}

func (this Memory) validate(javaVersion float32) (result *Memory, err error) {
	if this.PermMin != "" && javaVersion >= 1.8 {
		err = fmt.Errorf("java.memory.perm_min is not applicable to java(%v) >= 1.8", javaVersion)
		return
	}
	if this.PermMax != "" && javaVersion >= 1.8 {
		err = fmt.Errorf("java.memory.perm_max is not applicable to java(%v) >= 1.8", javaVersion)
		return
	}
	if this.MetaspaceMin != "" && javaVersion < 1.8 {
		err = fmt.Errorf("java.memory.metaspace_min is not applicable to java(%v) < 1.8", javaVersion)
		return
	}
	if this.MetaspaceMax != "" && javaVersion < 1.8 {
		err = fmt.Errorf("java.memory.metaspace_max is not applicable to java(%v) < 1.8", javaVersion)
		return
	}
	result = &this

	return
}

type JMX struct {
	Port *int
	SSL *bool
	Authenticate *string
}

