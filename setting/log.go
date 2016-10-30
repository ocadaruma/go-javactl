package setting

import (
	"fmt"
	"time"

	"github.com/ocadaruma/go-javactl/util"
)

type Log struct {
	ConsoleLog *ConsoleLog `yaml:"console"`
	GCLog *GCLog `yaml:"gc"`
	Dump *Dump
	ErrorLog *ErrorLog `yaml:"error"`
}

func (this Log) GetOpts(home string, now time.Time) []string {
	result := []string{}

	if this.GCLog != nil {
		result = append(result, this.GCLog.getOpts(home, now)...)
	}
	if this.Dump != nil {
		result = append(result, this.Dump.getOpts()...)
	}
	if this.ErrorLog != nil {
		result = append(result, this.ErrorLog.getOpts()...)
	}

	return result
}

type ConsoleLog struct {
	Prefix string
	MaxSize string `yaml:"max_size"`
	Backup int
	Preserve int
}

func (this ConsoleLog) GetNormalizedPath(home string, now time.Time) (result string) {
	if this.Prefix != "" {
		result = getPath(util.NormalizePath(this.Prefix, home), now)
	}
	return
}

type GCLog struct {
	Prefix string
	MaxSize string `yaml:"max_size"`
	Backup int
	Preserve int
}

func (this GCLog) GetNormalizedPath(home string, now time.Time) (result string) {
	if this.Prefix != "" {
		result = getPath(util.NormalizePath(this.Prefix, home), now)
	}
	return
}

func (this GCLog) getOpts(home string, now time.Time) []string {
	result := []string{}

	if this.Prefix != "" {
		result = append(result, []string{
			"-verbose:gc",
			"-XX:+PrintGCDateStamps",
			"-XX:+PrintGCDetails",
			fmt.Sprintf("-Xloggc:%s", this.GetNormalizedPath(home, now)),
		}...)

		if this.MaxSize != "" {
			result = append(result, []string{
				"-XX:+UseGCLogFileRotation",
				fmt.Sprintf("-XX:GCLogFileSize=%s", this.MaxSize),
			}...)

			if this.Backup > 0 {
				result = append(result, fmt.Sprintf("-XX:NumberOfGCLogFiles=%s", this.Backup))
			}
		}
	}

	return result
}

type Dump struct {
	Prefix string
}

func (this Dump) getOpts() []string {
	result := []string{}

	if this.Prefix != "" {
		result = append(result,
			"-XX:+HeapDumpOnOutOfMemoryError",
			fmt.Sprintf("-XX:HeapDumpPath=%s", this.Prefix),
		)
	}

	return result
}

type ErrorLog struct {
	Path string
}

func (this ErrorLog) getOpts() []string {
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
