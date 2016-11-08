package logger

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"sort"
	"strings"
	"time"
)

func TestConsoleLogger(t *testing.T) {
	dir, _ := ioutil.TempDir("", "temp")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "__test_console_logger_1.log")

	writer := NewConsoleLogger(path, 1024 * 1024 * 1, 2)

	zeros := make([]byte, 200000)
	for i := 0; i < 200000; i ++ { zeros[i] = '0' }

	ones := make([]byte, 400000)
	for i := 0; i < 400000; i ++ { ones[i] = '1' }

	// rotation 1
	twos := make([]byte, 600000)
	for i := 0; i < 600000; i ++ { twos[i] = '2' }

	// rotation 2
	threes := make([]byte, 500000)
	for i := 0; i < 500000; i ++ { threes[i] = '3' }

	fours := make([]byte, 250000)
	for i := 0; i < 250000; i ++ { fours[i] = '4' }

	// rotation 3
	fives := make([]byte, 1000000)
	for i := 0; i < 1000000; i ++ { fives[i] = '5' }

	for _, message := range []string{string(zeros), string(ones), string(twos), string(threes), string(fours), string(fives)} {
		io.WriteString(writer, message + "\n")
	}
	writer.Close()

	time.Sleep(time.Second)

	files,_ := ioutil.ReadDir(dir)
	if len(files) != 3 {
		t.Errorf("backups must be 2. files: %v", files)
	}

	withoutExt := []string{}
	for _, file := range files {
		withoutExt = append(withoutExt, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
	}
	sort.Strings(withoutExt)

	currentLog := filepath.Join(dir, withoutExt[0] + ".log")
	content, _ := ioutil.ReadFile(currentLog)
	if string(content)[26:1000026] != string(fives) {
		t.Error("current log must consists of 5.")
	}

	oldLog1 := filepath.Join(dir, withoutExt[1] + ".log")
	content, _ = ioutil.ReadFile(oldLog1)
	if string(content)[26:600026] != string(twos) {
		t.Error("old log 1 must consists of 2.")
	}

	oldLog2 := filepath.Join(dir, withoutExt[2] + ".log")
	content, _ = ioutil.ReadFile(oldLog2)
	split := strings.Split(string(content), "\n")
	if split[0][26:500026] != string(threes) {
		t.Error("line 1 of old log 2 must consists of 3.")
	}
	if split[1][26:250026] != string(fours) {
		t.Error("line 2 of old log 2 must consists of 4.")
	}
}

func TestConsoleLogger_Unicode(t *testing.T) {
	dir, _ := ioutil.TempDir("", "temp")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "__test_console_logger_2.log")

	writer := NewConsoleLogger(path, 1024 * 1024 * 1, 1)

	io.WriteString(writer, "あいうえお")
	writer.Close()

	files,_ := ioutil.ReadDir(dir)
	if len(files) != 1 {
		t.Error("backups must not exists.")
	}

	content, _ := ioutil.ReadFile(path)
	if string(content)[26:] != "あいうえお" {
		t.Errorf("log content(%s) must equal to あいうえお", string(content)[26:])
	}
}
