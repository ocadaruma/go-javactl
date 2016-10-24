package executor

import (
	"fmt"
	"os/user"

	"github.com/ocadaruma/javagtl/setting"
)

type Executor struct {
	Setting *setting.Setting
	Failed bool
}

//func (this *Executor) CheckRequirement() (err error) {
//	this.checkUser()
//	this.checkJavaVersion()
//	this.checkDuplicate()
//}

func (this *Executor) checkUser() (err error) {
	actual, _ := user.Current()
	expect := this.Setting.OSSetting.User

	if actual != expect {
		err = fmt.Errorf("This application must be run as '%s', but you are '%s'.", expect, actual)
	}
	return
}

//func (this *Executor) checkJavaVersion() error {
//
//}
//
//func (this *Executor) checkDuplicate() error {
//}
//
//func (this *Executor) createDirectories() error {
//
//}

//func Execute() {
//
//}
