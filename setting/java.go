package setting

import (
	"fmt"
	"path/filepath"

	"github.com/ocadaruma/go-javactl/util"
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

func (this Java) Normalize() (err error) {
	if this.Home == "" || !filepath.IsAbs(this.Home) {
		err = fmt.Errorf("java.home(%s) is required and must be an absolute path", this.Home)
		return
	}

	if this.Version < 1 {
		err = fmt.Errorf("invalid java.version(%f)", this.Version)
		return
	}

	err = this.Memory.validate(this.Version)

	return
}

func (this Java) getOpts() []string {

}

func (this Java) GetExecutable() string {
	return filepath.Join(this.Home, "bin", "java")
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

func (this Memory) validate(javaVersion float32) (err error) {
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

	return
}

func (this Memory) getOpts() (result []string) {
	opts := []string{
		util.EmptyIfZero("-Xms%s", this.HeapMin),
		util.EmptyIfZero("-Xmx%s", this.HeapMax),
		util.EmptyIfZero("-XX:PermSize=%s", this.PermMin),
		util.EmptyIfZero("-XX:MaxPermSize=%s", this.PermMax),
		util.EmptyIfZero("-XX:MetaspaceSize=%s", this.MetaspaceMin),
		util.EmptyIfZero("-XX:MaxMetaspaceSize=%s", this.MetaspaceMax),
		util.EmptyIfZero("-Xmn%s", this.NewMin),
		util.EmptyIfZero("-XX:MaxNewSize=%s", this.NewMax),
		util.EmptyIfZero("-XX:SurvivorRatio=%d", this.SurvivorRatio),
		util.EmptyIfZero("-XX:TargetSurvivorRatio=%d", this.TargetSurvivorRatio),
	}

	for _, o := range opts {
		if o != "" { result = append(result, o) }
	}

	return
}

type JMX struct {
	Port *int
	SSL *bool
	Authenticate *bool
}

func (this JMX) getOpts() (result []string) {
	if this.Port != nil {
		result = append(result,
			"-Dcom.sun.management.jmxremote",
			fmt.Sprintf("-Dcom.sun.management.jmxremote.port=%d", *this.Port),
		)

		if this.SSL != nil {
			result = append(result,
				fmt.Sprintf("-Dcom.sun.management.jmxremote.ssl=%t", *this.SSL))
		}
		if this.Authenticate != nil {
			result = append(result,
				fmt.Sprintf("-Dcom.sun.management.jmxremote.authenticate=%t", *this.Authenticate))
		}
	}
	return
}
