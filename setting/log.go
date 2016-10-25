package setting

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
