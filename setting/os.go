package setting

type OSSetting struct {
	User string
	Env map[string]string
}

func NewOSSetting(user string, env map[string]string) *OSSetting {
	return &OSSetting{user, env}
}
