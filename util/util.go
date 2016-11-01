package util

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	units = map[string]int64{
		"k": k,
		"m": m,
		"g": g,
	}
)

const (
	_1 int64 = 1
	k int64 = 1 << 10
	m int64 = 1 << 20
	g int64 = 1 << 30
)

type MemSize struct {
	Value int
	Unit int64
}

func (this MemSize) Bytes() int64 {
	return int64(this.Value) * this.Unit
}

func NewMemSize(size interface{}) *MemSize {
	switch s := size.(type) {
	case string:
		i, err := strconv.Atoi(s)
		if err == nil {
			return &MemSize{i, _1}
		} else if m, _ := regexp.MatchString("^([0-9]|[1-9][0-9]*)[kmgKMG]$", s); m {
			value, _ := strconv.Atoi(s[:len(s)-1])
			unit, exists := units[strings.ToLower(s[len(s)-1:])]
			if exists {
				return &MemSize{value, unit}
			} else {
				return &MemSize{value, _1}
			}
		} else {
			return nil
		}
	case int:
		return &MemSize{s, _1}
	default:
		return nil
	}
}

func NormalizePath(path string, baseDir string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(baseDir, path)
	}
}

func FmtIfNonZero(format string, value interface{}) (result string) {
	if value != reflect.Zero(reflect.TypeOf(value)).Interface() {
		result = fmt.Sprintf(format, value)
	}
	return
}

func FmtIfNonNilInt(format string, value *int) (result string) {
	if value != nil {
		result = fmt.Sprintf(format, *value)
	}
	return
}
