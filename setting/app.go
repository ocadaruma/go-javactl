package setting

import (
	"fmt"
	"path/filepath"

	"github.com/ocadaruma/go-javactl/setting/mapping"
	"github.com/ocadaruma/go-javactl/util"
)

type AppSetting struct {
	Name string
	Home string
	Jar string
	EntryPoint string
	Command string
	PidFile string
}

func NewAppSetting(app mapping.App) (result *AppSetting, err error) {
	if !filepath.IsAbs(app.Home) {
		err = fmt.Errorf("app.home(%s) must be an absolute path", app.Home)
		return
	}

	if (app.Jar != "") == (app.Command != "") {
		err = fmt.Errorf("either app.jar(%s) or app.command(%s) but not both must be given", app.Jar, app.Command)
		return
	}

	if (app.Jar == "") && (app.EntryPoint != "") {
		err = fmt.Errorf("app.entry_point(%s) must be used with app.jar(%s)", app.EntryPoint, app.Jar)
		return
	}

	result = &AppSetting{
		Name: app.Name,
		Home: app.Home,
		Jar: app.Jar,
		EntryPoint: app.EntryPoint,
		Command: app.Command,
		PidFile: app.PidFile,
	}

	// normalize paths
	if result.Jar != "" { result.Jar = util.NormalizePath(result.Jar, result.Home) }

	if result.Command != "" { result.Command = util.NormalizePath(result.Command, result.Home) }

	if result.PidFile != "" { result.PidFile = util.NormalizePath(result.PidFile, result.Home) }

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
