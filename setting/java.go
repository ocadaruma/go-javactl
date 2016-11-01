package setting

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/ocadaruma/go-javactl/setting/mapping"
	"github.com/ocadaruma/go-javactl/util"
)

type JavaSetting struct {
	Home string
	Version float32
	Server bool
	Memory *MemorySetting
	JMX *JMXSetting
	Env map[string]string
	Option []string
}

func NewJavaSetting(java mapping.Java) (result *JavaSetting, err error) {
	if java.Home == "" || !filepath.IsAbs(java.Home) {
		err = fmt.Errorf("java.home(%s) is required and must be an absolute path", java.Home)
		return
	}

	if java.Version < 1 {
		err = fmt.Errorf("invalid java.version(%f)", java.Version)
		return
	}

	var memory *MemorySetting
	if java.Memory != nil {
		memory, err = newMemorySetting(java.Version, *java.Memory)
		if err != nil { return }
	}

	var jmxSetting *JMXSetting
	if java.JMX != nil {
		jmx := newJMXSetting(*java.JMX)
		jmxSetting = &jmx
	}

	result = &JavaSetting{
		Home: java.Home,
		Version: java.Version,
		Server: java.Server,
		Memory: memory,
		JMX: jmxSetting,
		Env: java.Env,
		Option: java.Option,
	}

	return
}

func (this JavaSetting) GetArgs() []string {
	return append([]string{this.GetExecutable()}, this.getOpts()...)
}

func (this JavaSetting) GetExecutable() string {
	return filepath.Join(this.Home, "bin", "java")
}

func (this JavaSetting) getOpts() (result []string) {
	var memoryOpts []string
	if this.Memory != nil { memoryOpts = this.Memory.getOpts() }

	var jmxOpts []string
	if this.JMX != nil { jmxOpts = this.JMX.getOpts() }

	var server []string
	if this.Server { server = []string{"-server"} }

	var keys []string
	for k := range this.Env { keys = append(keys, k) }
	sort.Strings(keys)

	var env []string
	for _, k := range keys {
		env = append(env, fmt.Sprintf("-D%s=%s", k, this.Env[k]))
	}

	result = append(result, server...)
	result = append(result, memoryOpts...)
	result = append(result, jmxOpts...)
	result = append(result, env...)
	result = append(result, this.Option...)

	return
}

type MemorySetting struct {
	HeapMin string
	HeapMax string
	PermMin string
	PermMax string
	MetaspaceMin string
	MetaspaceMax string
	NewMin string
	NewMax string
	SurvivorRatio *int
	TargetSurvivorRatio *int
}

func newMemorySetting(javaVersion float32, memory mapping.Memory) (result *MemorySetting, err error) {
	if memory.PermMin != "" && javaVersion >= 1.8 {
		err = fmt.Errorf("java.memory.perm_min is not applicable to java(%v) >= 1.8", javaVersion)
		return
	}
	if memory.PermMax != "" && javaVersion >= 1.8 {
		err = fmt.Errorf("java.memory.perm_max is not applicable to java(%v) >= 1.8", javaVersion)
		return
	}
	if memory.MetaspaceMin != "" && javaVersion < 1.8 {
		err = fmt.Errorf("java.memory.metaspace_min is not applicable to java(%v) < 1.8", javaVersion)
		return
	}
	if memory.MetaspaceMax != "" && javaVersion < 1.8 {
		err = fmt.Errorf("java.memory.metaspace_max is not applicable to java(%v) < 1.8", javaVersion)
		return
	}

	result = &MemorySetting{
		HeapMin: memory.HeapMin,
		HeapMax: memory.HeapMax,
		PermMin: memory.PermMin,
		PermMax: memory.PermMax,
		MetaspaceMin: memory.MetaspaceMin,
		MetaspaceMax: memory.MetaspaceMax,
		NewMin: memory.NewMin,
		NewMax: memory.NewMax,
		SurvivorRatio: memory.SurvivorRatio,
		TargetSurvivorRatio: memory.TargetSurvivorRatio,
	}

	return
}

func (this MemorySetting) getOpts() (result []string) {
	opts := []string{
		util.FmtIfNonZero("-Xms%s", this.HeapMin),
		util.FmtIfNonZero("-Xmx%s", this.HeapMax),
		util.FmtIfNonZero("-XX:PermSize=%s", this.PermMin),
		util.FmtIfNonZero("-XX:MaxPermSize=%s", this.PermMax),
		util.FmtIfNonZero("-XX:MetaspaceSize=%s", this.MetaspaceMin),
		util.FmtIfNonZero("-XX:MaxMetaspaceSize=%s", this.MetaspaceMax),
		util.FmtIfNonZero("-Xmn%s", this.NewMin),
		util.FmtIfNonZero("-XX:MaxNewSize=%s", this.NewMax),
		util.FmtIfNonNilInt("-XX:SurvivorRatio=%d", this.SurvivorRatio),
		util.FmtIfNonNilInt("-XX:TargetSurvivorRatio=%d", this.TargetSurvivorRatio),
	}

	for _, o := range opts {
		if o != "" { result = append(result, o) }
	}

	return
}

type JMXSetting struct {
	Port *int
	SSL *bool
	Authenticate *bool
}

func newJMXSetting(jmx mapping.JMX) JMXSetting {
	return JMXSetting{
		Port: jmx.Port,
		SSL: jmx.SSL,
		Authenticate: jmx.Authenticate,
	}
}

func (this JMXSetting) getOpts() (result []string) {
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
