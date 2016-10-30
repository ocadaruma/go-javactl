package setting

type Log struct {
	ConsoleLog *ConsoleLog `yaml:"console"`
	GCLog *GCLog `yaml:"gc"`
	Dump *Dump
	ErrorLog *ErrorLog `yaml:"error"`
}

//func (this Log) GetOpts() []string {
//	result := []string{}
//
//	if this.ConsoleLog != nil {
//		//append(result, this.ConsoleLog.)
//	}
//}

type ConsoleLog struct {
	Prefix string
	MaxSize string `yaml:"max_size"`
	Backup int
	Preserve int
}

//func (this ConsoleLog) getPath() string {
//
//}

type GCLog struct {
	Prefix string
	MaxSize string `yaml:"max_size"`
	Backup int
	Preserve int
}

//func (this GCLog) getPath() string {
//
//}
//
//func (this GCLog) getOpts() []string {
//
//}
//
type Dump struct {
	Prefix string
}
//
//func (this Dump) getOpts() []string {
//
//}
//
type ErrorLog struct {
	Path string
}
//
//func (this ErrorLog) getOpts() []string {
//
//}
