package setting

type LogSetting struct {
	Console Console
}

type Console struct {
	Home string
	Prefix string
	MaxSize string
	Backup string
	Preserve string
}
