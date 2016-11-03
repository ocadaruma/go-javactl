package setting

import (
	"time"

	"github.com/ocadaruma/go-javactl/setting/mapping"
	"github.com/ocadaruma/go-javactl/util"
)

type Setting struct {
	App AppSetting
	Java JavaSetting
	Log *LogSetting
	OS OSSetting
	PreCommands  []string
	PostCommands []string
}

func NewSetting(config *mapping.YAMLConfig) (result *Setting, err error) {
	var app *AppSetting
	app, err = NewAppSetting(&config.App)
	if err != nil { return }

	var java *JavaSetting
	java, err = NewJavaSetting(&config.Java)
	if err != nil { return }

	var logSetting *LogSetting
	if config.Log != nil {
		log := NewLogSetting(app.Home, config.Log)
		logSetting = &log
	}

	os := NewOSSetting(&config.OS)

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

func (this *Setting) GetArgs(extraArgs []string, now time.Time) []string {
	var logOpts []string
	if this.Log != nil { logOpts = this.Log.GetOpts(now) }
	return append(
		this.App.GetArgs(append(this.Java.GetArgs(), logOpts...)),
		extraArgs...)
}

func (this *Setting) GetEnviron(now time.Time) map[string]string {
	result := make(map[string]string)

	result["JAVA_HOME"] = this.Java.Home

	var logOpts []string
	if this.Log != nil { logOpts = this.Log.GetOpts(now) }
	result["JAVA_OPTS"] = util.List2Cmdline(append(this.Java.getOpts(), logOpts...))

	for key, value := range this.OS.Env {
		result[key] = value
	}

	return result
}
