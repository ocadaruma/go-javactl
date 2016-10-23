package setting

import (
	"path/filepath"

	"github.com/ocadaruma/javagtl/assert"
	"github.com/ocadaruma/javagtl/util"
)

type AppSetting struct {
	Name string
	Home string
	Jar string
	EntryPoint string
	Command string
	PidFile string
}

func NewAppSetting(name string, home string, jar string, entryPoint string, command string, pidFile string) AppSetting {
	assert.Assert(filepath.IsAbs(home), "app.home must be an absolute path")
	assert.Assert((jar != "") != (command != ""), "either app.jar or app.command but not both must be given")
	assert.Assert((jar != "") || (entryPoint == ""), "app.entryPoint must be used with app.jar")

	var j string
	if jar != "" { j = util.NormalizePath(jar, home) }

	var c string
	if command != "" { c = util.NormalizePath(command, home) }

	var p string
	if pidFile != "" { p = util.NormalizePath(pidFile, home) }
	return AppSetting{
		Name: name,
		Home: home,
		Jar: j,
		EntryPoint: entryPoint,
		Command: c,
		PidFile: p,
	}
}

func (this AppSetting) IsDuplicateAllowed() bool {
	return this.PidFile != ""
}

func (this AppSetting) GetArgs(javaArgs []string) []string {
	if this.Jar != "" {
		if this.EntryPoint != "" {
			return append(javaArgs, "-cp", this.Jar, this.EntryPoint)
		} else {
			return append(javaArgs, "-jar", this.Jar)
		}
	} else {
		return []string{this.Command}
	}
}
