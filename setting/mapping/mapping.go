package mapping

type App struct {
	Name string
	Home string
	Jar string
	EntryPoint string `yaml:"entry_point"`
	Command string
	PidFile string `yaml:"pid_file"`
}

type Java struct {
	Home string
	Version float32
	Server bool
	Memory *Memory
	JMX *JMX
	Env map[string]string
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
