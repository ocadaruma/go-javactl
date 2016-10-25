package setting

import (
	"fmt"
	"path/filepath"

	"github.com/ocadaruma/javagtl/util"
)

type App struct {
	Name string
	Home string
	Jar string
	EntryPoint string `yaml:"entry_point"`
	Command string
	PidFile string `yaml:"pid_file"`
}

func (this App) Normalize() (setting *App, err error) {
	if !filepath.IsAbs(this.Home) {
		err = fmt.Errorf("app.home(%s) must be an absolute path", this.Home)
		return
	}

	if (this.Jar != "") == (this.Command != "") {
		err = fmt.Errorf("either app.jar(%s) or app.command(%s) but not both must be given", this.Jar, this.Command)
		return
	}

	if (this.Jar == "") && (this.EntryPoint != "") {
		err = fmt.Errorf("app.entry_point(%s) must be used with app.jar(%s)", this.EntryPoint, this.Jar)
		return
	}

	var j string
	if this.Jar != "" { j = util.NormalizePath(this.Jar, this.Home) }

	var c string
	if this.Command != "" { c = util.NormalizePath(this.Command, this.Home) }

	var p string
	if this.PidFile != "" { p = util.NormalizePath(this.PidFile, this.Home) }
	setting = &App{
		Name: this.Name,
		Home: this.Home,
		Jar: j,
		EntryPoint: this.EntryPoint,
		Command: c,
		PidFile: p,
	}
	return
}

func (this App) IsDuplicateAllowed() bool {
	return this.PidFile != ""
}

func (this App) GetArgs(javaArgs []string) []string {
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
