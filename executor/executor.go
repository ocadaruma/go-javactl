package executor

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ocadaruma/go-javactl/logger"
	"github.com/ocadaruma/go-javactl/option"
	"github.com/ocadaruma/go-javactl/setting"
)

type Executor struct {
	Logger *logger.Logger
	Setting *setting.Setting
	Opts *option.JavactlOpts
}

func NewExecutor(log *logger.Logger, sett *setting.Setting, opts *option.JavactlOpts) (result *Executor) {
	result = &Executor{
		Logger: log,
		Setting: sett,
		Opts: opts,
	}

	return
}

func (this *Executor) CheckRequirement() (err error) {
	err = this.checkUser()
	if err != nil { return }

	err = this.checkJavaVersion()
	if err != nil { return }

	err = this.checkDuplicateProcess()

	return
}

func (this *Executor) checkUser() (err error) {
	actual, _ := user.Current()
	expect := this.Setting.OS.User

	if actual.Username != expect {
		err = fmt.Errorf("This application must be run as '%s', but you are '%s'.", expect, actual.Username)
	}
	return
}

func (this *Executor) checkJavaVersion() (err error) {
	var out []byte
	out, err = exec.Command(this.Setting.Java.GetExecutable(), "-version").CombinedOutput()
	if err != nil { return }

	actual := regexp.MustCompile(`java version "(\d+[.]\d+)[.]\d+_\d+"`).FindStringSubmatch(string(out))[0]
	versionString := fmt.Sprintf("%.1f", this.Setting.Java.Version)

	if versionString != actual {
		err = fmt.Errorf("Unexpected Java version: expect='%s', actual='%s'.", versionString, actual)
	}

	return
}

func (this *Executor) checkDuplicateProcess() (err error) {
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
				this.Logger.Info(
					fmt.Sprintf("%s is already running. Skipped: config=%s",
						this.Setting.App.Name, this.Opts.ConfigPath))
				return fmt.Errorf("%s is already running.", this.Setting.App.Name)
			}
		}
	}
	return
}

func (this *Executor) CreateDirectories() (err error) {
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

func (this *Executor) CleanOldLogs(now time.Time) (err error) {
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

func (this *Executor) deleteFiles(prefix string) (err error) {
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

func (this *Executor) Execute(now time.Time) error {
	allErrors := []error{}

	var err error
	err = this.executeCommands(this.Setting.PreCommands, now)
	if err != nil { allErrors = append(allErrors, err) }

	if err == nil {
		err = this.executeApplication(now)
		if err != nil { allErrors = append(allErrors, err) }
	}

	err = this.executeCommands(this.Setting.PostCommands, now)
	if err != nil { allErrors = append(allErrors, err) }

	errStrings := []string{}
	for _, e := range allErrors {
		errStrings = append(errStrings, fmt.Sprintf("  - %v", e))
	}

	if len(errStrings) > 0 {
		return fmt.Errorf("Failed to execute. errors: [\n%s\n]\n", strings.Join(errStrings, "\n"))
	} else {
		return nil
	}
}

func (this *Executor) executeApplication(now time.Time) (err error) {
	this.Logger.Info(
		fmt.Sprintf("%s started: config=%s, args=%s",
			this.Setting.App.Name, this.Opts.ConfigPath, this.Opts.ExtraArgs))

	startTime := time.Now()
	if this.Opts.DryRun {
		var args, env []string
		indentation := "  "
		for _, arg := range this.Setting.GetArgs(this.Opts.ExtraArgs, now) {
			args = append(args, indentation + arg)
		}
		for _, e := range this.Setting.GetEnviron(now) {
			env = append(env, indentation + e)
		}
		var output string
		if this.Setting.Log != nil && this.Setting.Log.ConsoleLog != nil {
			output = this.Setting.Log.ConsoleLog.GetPath(now)
		} else {
			output = "stdout"
		}
		fmt.Printf("Would execute: cwd=%s, cmd=[\n%s\n], env={\n%s\n}, output=%s\n",
			this.Setting.App.Home,
			strings.Join(args, "\n"),
			strings.Join(env, "\n"),
			output)
	} else {
		var stdout, stderr io.Writer
		if this.Setting.Log != nil && this.Setting.Log.ConsoleLog != nil {
			outpath := this.Setting.Log.ConsoleLog.GetPath(now)
			var maxBytes int64
			if this.Setting.Log.ConsoleLog.MaxSize != nil {
				maxBytes = this.Setting.Log.ConsoleLog.MaxSize.Bytes()
			}
			out := logger.NewConsoleLogger(outpath, maxBytes, this.Setting.Log.ConsoleLog.Backup)
			stdout = out
			stderr = out
		} else {
			stdout = os.Stdout
			stderr = os.Stderr
		}
		env := make(map[string]string)
		for k, v := range osEnviron() {
			env[k] = v
		}
		for k, v := range this.Setting.GetEnviron(now) {
			env[k] = v
		}
		args := subprocessArgs{
			shell: false,
			args: this.Setting.GetArgs(this.Opts.ExtraArgs, now),
			cwd: this.Setting.App.Home,
			pidFile: this.Setting.App.PidFile,
			env: env,
			stdin: os.Stdin,
			stdout: stdout,
			stderr: stderr,
		}
		err = callSubprocess(&args)
	}
	elapsed := time.Now().Sub(startTime)

	if err != nil {
		this.Logger.Error(
			fmt.Sprintf("%s ended with error: error=%v, elapsed=%ds",
				this.Setting.App.Name, err, elapsed.Seconds()))
	} else {
		this.Logger.Info(
			fmt.Sprintf("%s ended successfully: elapsed=%ds",
				this.Setting.App.Name, elapsed.Seconds()))
	}

	return
}

// execute commands sequentially.
// attempt to execute all commands even if some command fails.
func (this *Executor) executeCommands(commands []string, now time.Time) error {
	allErrors := []error{}

	for _, cmd := range commands {
		if this.Opts.DryRun {
			fmt.Printf("Would execute: %s\n", cmd)
		} else {
			args := subprocessArgs{
				shell: true,
				args: []string{cmd},
				cwd: this.Setting.App.Home,
				env: this.Setting.GetEnviron(now),
			}
			err := callSubprocess(&args)

			if err != nil {
				this.Logger.Error(fmt.Sprintf("Failed to execute: app=%s, cmd=%s, error: %v",
					this.Setting.App.Name, cmd, err))
				allErrors = append(allErrors, err)
			}
		}
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("some commands failed. %v", allErrors)
	} else {
		return nil
	}
}

type subprocessArgs struct {
	shell bool
	args []string
	cwd string
	pidFile string
	env map[string]string
	stdin io.Reader
	stdout io.Writer
	stderr io.Writer
}

func callSubprocess(args *subprocessArgs) (err error) {
	var wd string
	wd, err = os.Getwd()
	if err != nil { return }

	defer os.Chdir(wd)

	// if cwd is specified, exec command in cwd.
	if args.cwd != "" {
		err = os.Chdir(args.cwd)
		if err != nil { return }
	}

	var cmd *exec.Cmd
	if args.shell {
		osname := runtime.GOOS
		if osname == "windows" {
			shell := os.Getenv("COMSPEC")
			cmd = exec.Command(shell, append([]string{"/c"}, args.args...)...)
		} else {
			shell := os.Getenv("SHELL")
			cmd = exec.Command(shell, append([]string{"-c"}, args.args...)...)
		}
	} else {
		if len(args.args) < 1 {
			err = fmt.Errorf("args must not be empty. args: %v", args.args)
			return
		}
		cmd = exec.Command(args.args[0], args.args[1:]...)
	}

	var envPairs []string
	for key, value := range args.env {
		envPairs = append(envPairs, fmt.Sprintf("%s=%s", key, value))
	}
	if len(envPairs) > 0 { cmd.Env = envPairs }

	if args.stdout == nil { cmd.Stdout = os.Stdout } else { cmd.Stdout = args.stdout }
	if args.stderr == nil { cmd.Stderr = os.Stderr } else { cmd.Stderr = args.stderr }
	if args.stdin != nil { cmd.Stdin = args.stdin }

	if args.pidFile == "" {
		err = cmd.Run()
	} else {
		err = cmd.Start()
		if err != nil { return }

		pid := cmd.Process.Pid
		defer os.Remove(args.pidFile)

		err = ioutil.WriteFile(args.pidFile, []byte(string(pid)), 0644)
		if err != nil { return }

		err = cmd.Wait()
	}

	return
}

func osEnviron() map[string]string {
	items := make(map[string]string)
	entries := os.Environ()

	for _, entry := range entries {
		split := strings.SplitN(entry, "=", 2)
		items[split[0]] = split[1]
	}

	return items
}
