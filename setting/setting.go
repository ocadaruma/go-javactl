package setting

type Setting struct {
	ConfigPath string
	ExtraArgs []string
	DryRun bool
	Debug bool
	AppSetting AppSetting
	OSSetting OSSetting
	PreCommands []string
	PostCommands []string
}

//func NewSetting(configPath string, extraArgs []string, dryRun bool) *Setting {
//
//}
//
//func ParseArgs(argv ) {
//
//}



//
//type JavaSetting struct {
//
//}
//
//type LogSetting struct {
//
//}
//
//type OSSetting struct {
//
//}
//
//type Setting struct {
//	ConfigPath string
//	ExtraArgs *[]string
//	DryRun bool
//	Debug bool
//	AppSetting AppSetting
//	JavaSetting JavaSetting
//	LogSetting LogSetting
//	OSSetting OSSetting
//	PreCommands *[]string
//	PostCommands *[]string
//}
//
//func (setting *AppSetting) GetArgs(javaArgs *[]string) string {
//
//}
