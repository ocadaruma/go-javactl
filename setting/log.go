package setting

import (
	"fmt"
	"time"

	"github.com/ocadaruma/go-javactl/setting/mapping"
	"github.com/ocadaruma/go-javactl/util"
)

type LogSetting struct {
	ConsoleLog *ConsoleLogSetting
	GCLog *GCLogSetting
	Dump *DumpSetting
	ErrorLog *ErrorLogSetting
}

func NewLogSetting(home string, log mapping.Log) (result LogSetting) {
	if log.ConsoleLog != nil {
		console := newConsoleLogSetting(home, *log.ConsoleLog)
		result.ConsoleLog = &console
	}
	if log.GCLog != nil {
		gc := newGCLogSetting(home, *log.GCLog)
		result.GCLog = &gc
	}
	if log.Dump != nil {
		dump := newDumpSetting(home, *log.Dump)
		result.Dump = &dump
	}
	if log.ErrorLog != nil {
		errorLog := newErrorLogSetting(home, *log.ErrorLog)
		result.ErrorLog = &errorLog
	}

	return
}

func (this LogSetting) GetOpts(now time.Time) (result []string) {
	if this.GCLog != nil {
		result = append(result, this.GCLog.getOpts(now)...)
	}
	if this.Dump != nil {
		result = append(result, this.Dump.getOpts()...)
	}
	if this.ErrorLog != nil {
		result = append(result, this.ErrorLog.getOpts()...)
	}

	return
}

type ConsoleLogSetting struct {
	Prefix string
	MaxSize *util.MemSize
	Backup int
	Preserve int
}

func newConsoleLogSetting(home string, consoleLog mapping.ConsoleLog) (result ConsoleLogSetting) {
	result = ConsoleLogSetting{
		Prefix: consoleLog.Prefix,
		MaxSize: nil,
		Backup: consoleLog.Backup,
		Preserve: consoleLog.Preserve,
	}

	if consoleLog.MaxSize != "" {
		result.MaxSize = util.NewMemSize(consoleLog.MaxSize)
	}

	if result.Prefix != "" {
		result.Prefix = util.NormalizePath(result.Prefix, home)
	}

	return
}

func (this ConsoleLogSetting) GetPath(now time.Time) (result string) {
	if this.Prefix != "" {
		result = getPath(this.Prefix, now)
	}
	return
}

type GCLogSetting struct {
	Prefix string
	MaxSize string
	Backup int
	Preserve int
}

func (this GCLogSetting) GetPath(now time.Time) (result string) {
	if this.Prefix != "" {
		result = getPath(this.Prefix, now)
	}
	return
}

func (this GCLogSetting) getOpts(now time.Time) (result []string) {
	if this.Prefix != "" {
		result = append(result, []string{
			"-verbose:gc",
			"-XX:+PrintGCDateStamps",
			"-XX:+PrintGCDetails",
			fmt.Sprintf("-Xloggc:%s", this.GetPath(now)),
		}...)

		if this.MaxSize != "" {
			result = append(result, []string{
				"-XX:+UseGCLogFileRotation",
				fmt.Sprintf("-XX:GCLogFileSize=%s", this.MaxSize),
			}...)

			if this.Backup > 0 {
				result = append(result, fmt.Sprintf("-XX:NumberOfGCLogFiles=%d", this.Backup))
			}
		}
	}
	return
}

func newGCLogSetting(home string, gcLog mapping.GCLog) (result GCLogSetting) {
	result = GCLogSetting{
		Prefix: gcLog.Prefix,
		MaxSize: gcLog.MaxSize,
		Backup: gcLog.Backup,
		Preserve: gcLog.Preserve,
	}
	if result.Prefix != "" {
		result.Prefix = util.NormalizePath(result.Prefix, home)
	}
	return
}

type DumpSetting struct {
	Prefix string
}

func newDumpSetting(home string, dump mapping.Dump) (result DumpSetting) {
	result = DumpSetting{
		Prefix: dump.Prefix,
	}
	if result.Prefix != "" { result.Prefix = util.NormalizePath(result.Prefix, home) }
	return
}

func (this DumpSetting) getOpts() (result []string) {
	if this.Prefix != "" {
		result = append(result,
			"-XX:+HeapDumpOnOutOfMemoryError",
			fmt.Sprintf("-XX:HeapDumpPath=%s", this.Prefix),
		)
	}
	return
}

type ErrorLogSetting struct {
	Path string
}

func newErrorLogSetting(home string, errorLog mapping.ErrorLog) (result ErrorLogSetting) {
	result = ErrorLogSetting{
		Path: errorLog.Path,
	}
	if result.Path != "" { result.Path = util.NormalizePath(result.Path, home) }
	return
}

func (this ErrorLogSetting) getOpts() []string {
	if this.Path != "" {
		return []string{fmt.Sprintf("-XX:ErrorFile=%s", this.Path)}
	} else {
		return []string{}
	}
}

func getPath(prefix string, now time.Time) (result string) {
	if prefix != "" {
		suffix := now.Format("_20060102_150405.log")
		result = fmt.Sprintf("%s%s", prefix, suffix)
	}
	return
}
