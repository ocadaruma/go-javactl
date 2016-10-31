package executor

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ocadaruma/go-javactl/option"
	"github.com/ocadaruma/go-javactl/setting"
)

type Executor struct {
	Setting setting.Setting
	Opts option.JavactlOpts
	Failed bool
}

func NewExecutor() Executor {

}

func (this Executor) CheckRequirement() (err error) {
	err = this.checkUser()
	if err != nil { return }

	err = this.checkJavaVersion()
	if err != nil { return }

	err = this.checkDuplicateProcess()

	return
}

func (this Executor) checkUser() (err error) {
	actual, _ := user.Current()
	expect := this.Setting.OS.User

	if actual != expect {
		err = fmt.Errorf("This application must be run as '%s', but you are '%s'.", expect, actual)
	}
	return
}

func (this Executor) checkJavaVersion() (err error) {
	var out []byte
	out, err = exec.Command(this.Setting.Java.GetExecutable(), "-version").CombinedOutput()
	if err != nil { return }

	actual := regexp.MustCompile(`java version "(\d+[.]\d+)[.]\d+_\d+"`).FindStringSubmatch(out)
	versionString := fmt.Sprintf("%.1f", this.Setting.Java.Version)

	if versionString != actual {
		err = fmt.Errorf("Unexpected Java version: expect='%s', actual='%s'.", versionString, actual)
	}

	return
}

func (this Executor) checkDuplicateProcess() (err error) {
	if !this.Setting.App.IsDuplicateAllowed() {
		// check pid file existence
		_, e := os.Stat(this.Setting.App.PidFile)

		if e == nil {
			var content []byte
			content, err = ioutil.ReadFile(this.Setting.App.PidFile)
			if err != nil { return }

			var pid int
			pid, err = strconv.Atoi(string(content))
			if err != nil { return }

			_, err = os.FindProcess(pid)
			if err == nil {
				return fmt.Errorf("%s is already running.", this.Setting.App.Name)
			}
		}
	}
	return
}

func (this Executor) CreateDirectories() (err error) {
	if this.Setting.Log == nil { return }

	log := *this.Setting.Log

	var consoleLogDir string
	if log.ConsoleLog != nil && log.ConsoleLog.Prefix != "" {
		consoleLogDir = filepath.Dir(log.ConsoleLog.Prefix)
	}

	var gcLogDir string
	if log.GCLog != nil && log.GCLog.Prefix != "" {
		gcLogDir = filepath.Dir(log.GCLog.Prefix)
	}

	var dumpPrefix string
	if log.Dump != nil { dumpPrefix = log.Dump.Prefix }

	var errorLogDir string
	if log.ErrorLog != nil && log.ErrorLog.Path != "" {
		errorLogDir = filepath.Dir(log.ErrorLog.Path)
	}

	for _, dir := range []string{consoleLogDir, gcLogDir, dumpPrefix, errorLogDir} {
		if dir != "" {
			if this.Opts.DryRun {
				fmt.Printf("Would create directory: %s\n", dir)
			} else {
				fmt.Printf("Creating directory: %s\n", dir)
				err = os.MkdirAll(dir, os.ModeDir)
				if err != nil { return }
			}
		}
	}

	return
}

func (this Executor) CleanOldLogs(now time.Time) (err error) {
	if log := this.Setting.Log; log != nil {
		// clear console logs
		if console := log.ConsoleLog; console != nil {
			if console.Prefix != "" && console.Preserve > 0 {
				err = this.deleteFiles(console.Prefix)
				if err != nil { return }
			}
		}
		// clear gc logs
		if gc := log.GCLog; gc != nil {
			if gc.Prefix != "" && gc.Preserve > 0 {
				err = this.deleteFiles(gc.Prefix)
				if err != nil { return }
			}
		}
	}
	return
}

func (this Executor) deleteFiles(prefix string) (err error) {
	var files []string
	files, err = filepath.Glob(prefix + "*")

	if err != nil { return }

	for _, path := range files {
		if this.Opts.DryRun {
			fmt.Printf("Would delete file: %s\n", path)
		} else {
			fmt.Printf("Deleting file: %s\n", path)
			err = os.Remove(path)
			if err != nil { return }
		}
	}

	return
}

func (this Executor) Execute(now time.Time) (err error) {
	err = this.executeCommands(this.Setting.PreCommands, now)
	if err != nil { return }

	err = this.executeApplication(now)
	if err != nil { return }

	err = this.executeCommands(this.Setting.PostCommands, now)

	return
}

func (this Executor) executeApplication(now time.Time) (err error) {

}

func (this Executor) executeCommands(commands []string, now time.Time) (err error) {
	failed := this.Failed
	for cmd := range commands {
		if this.Opts.DryRun {
			fmt.Printf("Would execute: %s\n", cmd)
		} else {

		}
	}
}

//func (this *Executor) createDirectories() error {
//
//}

//func Execute() {
//
//}
