package setting

import (
	"fmt"
	"path/filepath"

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

func NewAppSetting(name string, home string, jar string, entryPoint string, command string, pidFile string) (setting *AppSetting, err error) {
	if !filepath.IsAbs(home) {
		err = fmt.Errorf("app.home(%s) must be an absolute path", home)
		return
	}

	if (jar != "") != (command != "") {
		err = fmt.Errorf("either app.jar(%s) or app.command(%s) but not both must be given", jar, command)
		return
	}

	if (jar == "") && (entryPoint != "") {
		err = fmt.Errorf("app.entry_point(%s) must be used with app.jar(%s)", entryPoint, jar)
		return
	}

	var j string
	if jar != "" { j = util.NormalizePath(jar, home) }

	var c string
	if command != "" { c = util.NormalizePath(command, home) }

	var p string
	if pidFile != "" { p = util.NormalizePath(pidFile, home) }
	setting = &AppSetting{
		Name: name,
		Home: home,
		Jar: j,
		EntryPoint: entryPoint,
		Command: c,
		PidFile: p,
	}
	return
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
