package option

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

const (
	usage = "[options...] <config_path> [args...]"
)

type JavactlOpts struct {
	ConfigPath string
	ExtraArgs []string
	DryRun bool `long:"check" description:"dry-run mode"`
	Debug bool `long:"debug" description:"debug mode"`
	Version func() `short:"v" long:"version" description:"show program's version number and exit"`
}

func ParseArgs(appVersion string, args []string) (result *JavactlOpts, err error) {
	opts := JavactlOpts{}
	opts.Version = func() {
		fmt.Printf("javactl %s\n", appVersion)
		os.Exit(0)
	}

	result = &opts
	parser := flags.NewParser(result, flags.Default)
	parser.Usage = usage

	var rest []string
	rest, err = parser.ParseArgs(args)

	if err != nil { return }

	if len(rest) < 1 {
		fmt.Println("must specify config path.")
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	result.ConfigPath = rest[0]
	result.ExtraArgs = rest[1:]

	return
}
