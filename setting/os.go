package setting

import (
	"github.com/ocadaruma/go-javactl/setting/mapping"
)

type OSSetting struct {
	User string
	Env map[string]string
}

func NewOSSetting(os *mapping.OS) OSSetting {
	return OSSetting{
		User: os.User,
		Env: os.Env,
	}
}
