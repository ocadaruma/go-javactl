package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func tearDown(dir string) {
	os.RemoveAll(dir)
}

func TestJavactl(t *testing.T) {
	//================
	// create temporary directory
	//================
	curdir, _ := os.Getwd()
	templateFile := filepath.Join(curdir, "testdata", "test_01.yml.template")
	tempDir, _ := ioutil.TempDir(filepath.Dir(templateFile), "temp")

	defer tearDown(tempDir)

	//================
	// create configuration file from template
	//================
	confFile := filepath.Join(tempDir, "test_01.yml")
	tmpl, _ := template.ParseFiles(templateFile)

	var b bytes.Buffer
	user, _ := user.Current()
	tmpl.Execute(&b, map[string]string{"tempdir": tempDir, "curdir": curdir, "os_user": user.Username})

	ioutil.WriteFile(confFile, b.Bytes(), 0644)

	//================
	// run main and capture stdout, stderr
	//================
	origStdout := os.Stdout
	origStderr := os.Stderr

	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()

	os.Stdout = wStdout
	os.Stderr = wStderr

	os.Args = []string{"javactl", confFile}
	main()

	os.Stdout = origStdout
	os.Stderr = origStderr

	wStdout.Close()
	wStderr.Close()

	var bufStdout, bufStderr bytes.Buffer
	io.Copy(&bufStdout, rStdout)
	io.Copy(&bufStderr, rStderr)

	//================
	// assertion
	//================

	// assert log directory existence
	logDir := filepath.Join(tempDir, "logs")
	logDirs := []string{
		filepath.Join(logDir, "console"),
		filepath.Join(logDir, "gc"),
		filepath.Join(logDir, "dump"),
	}
	logDirLen := len(logDir)

	for _, dir := range logDirs {
		_, err := os.Stat(dir)
		if err != nil {
			t.Errorf("directory %s must exist", dir)
		}
	}

	// assert stdout, stderr
	stdout := string(bufStdout.Bytes())
	stderr := string(bufStderr.Bytes())
	expectedStdout := fmt.Sprintf("Creating directory: %s\nCreating directory: %s\nCreating directory: %s\n",
		logDirs[0], logDirs[1], logDirs[2])

	if stdout != expectedStdout {
		t.Errorf("stdout(%s) must equal to %s", stdout, expectedStdout)
	}
	if stderr != "" {
		t.Errorf("stderr(%s) must be empty", stderr)
	}

	// assert log file existence
	files, _ := ioutil.ReadDir(logDirs[0])
	if len(files) != 1 {
		t.Errorf("count of log files must be %d", len(files))
	}

	// assert log content
	bytes, _ := ioutil.ReadFile(filepath.Join(logDirs[0], files[0].Name()))
	logContent := string(bytes)

	expected1 := strings.Join([]string{
		" -server -Xms64M -Xmx2G -XX:MetaspaceSize=1G -XX:MaxMetaspaceSize=2G -Xmn256M -XX:MaxNewSize=256M ",
		"-XX:SurvivorRatio=8 -XX:TargetSurvivorRatio=50 -Dcom.sun.management.jmxremote ",
		"-Dcom.sun.management.jmxremote.port=20001 -Dcom.sun.management.jmxremote.ssl=false ",
		"-Dcom.sun.management.jmxremote.authenticate=false -Dcom.amazonaws.sdk.disableCertChecking=true ",
		"-Dfile.encoding=UTF-8 -Dhttp.netty.maxInitialLineLength=8192 -Dhttp.port=9000 ",
		"-XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:+UseCMSInitiatingOccupancyOnly ",
		"-XX:CMSInitiatingOccupancyFraction=70 -XX:+ScavengeBeforeFullGC -XX:+CMSScavengeBeforeRemark ",
		"-verbose:gc -XX:+PrintGCDateStamps -XX:+PrintGCDetails ",
		fmt.Sprintf("-Xloggc:%s/gc_", logDirs[1]),
	}, "")
	expected2 := strings.Join([]string{
		" -XX:+UseGCLogFileRotation -XX:GCLogFileSize=10M -XX:NumberOfGCLogFiles=10 ",
		"-XX:+HeapDumpOnOutOfMemoryError ",
		fmt.Sprintf("-XX:HeapDumpPath=%s ", logDirs[2]),
		fmt.Sprintf("-XX:ErrorFile=%s/hs_error_pid%%p.log ", logDir),
		fmt.Sprintf("-jar %s/bin/your-app-assembly-0.1.0.jar", tempDir),
	}, "")

	if logContent[25:709+logDirLen] != expected1 {
		t.Errorf("log content(%s) must equal to (%s)", logContent[25:709+logDirLen], expected1)
	}
	if logContent[728+logDirLen:] != expected2 {
		t.Errorf("log content(%s) must equal to (%s)", logContent[728+logDirLen:], expected2)
	}
}
