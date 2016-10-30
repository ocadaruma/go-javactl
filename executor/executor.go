package executor

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"os"
	"os/exec"
	"os/user"
	"strconv"

	"github.com/ocadaruma/go-javactl/setting"
)

type Executor struct {
	Setting setting.Setting
	Failed bool
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

func (this Executor) executeCommands(commands []string) {

}

//func (this *Executor) createDirectories() error {
//
//}

//func Execute() {
//
//}
