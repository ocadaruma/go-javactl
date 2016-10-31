package setting

import (
	"time"

	"github.com/ocadaruma/go-javactl/setting/mapping"
)

type Setting struct {
	App AppSetting
	Java JavaSetting
	Log *LogSetting
	OS OSSetting
	PreCommands  []string
	PostCommands []string
}

func NewSetting(config mapping.YAMLConfig) (result *Setting, err error) {
	var app *AppSetting
	app, err = NewAppSetting(config.App)
	if err != nil { return }

	var java *JavaSetting
	java, err = NewJavaSetting(config.Java)
	if err != nil { return }

	var logSetting *LogSetting
	if config.Log != nil {
		log := NewLogSetting(app.Home, *config.Log)
		logSetting = &log
	}

	os := NewOSSetting(config.OS)

	result = &Setting{
		App: *app,
		Java: *java,
		Log: logSetting,
		OS: os,
		PreCommands: config.PreCommands,
		PostCommands: config.PostCommands,
	}

	return
}

func (this Setting) GetArgs(extraArgs []string, now time.Time) []string {
	return append(
		this.App.GetArgs(append(this.Java.GetArgs(), this.Log.GetOpts(now)...)),
		extraArgs...
	)
}
