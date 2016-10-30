package setting

import (
	"reflect"
	"testing"
	"time"
)

func TestSetting_GetArgs(t *testing.T) {
	s, err := LoadConfig("../testdata/example.yml")

	if err != nil {
		t.Errorf("load must not failed. err : %v", err)
	}

	args := s.GetArgs(time.Date(2015, 9, 10, 12, 34, 56, 789, time.Local))

	testCases0 := []struct{
		expect, actual string
	} {
		{args[0], "/usr/java/latest/bin/java"},
		{args[len(args)-2], "-jar"},
		{args[len(args)-1], "/path/to/your-app/bin/your-app-assembly-0.1.0.jar"},
	}

	for i, c := range testCases0 {
		if c.actual != c.expect {
			t.Errorf("case %v : %v must equal to %v", i, c.actual, c.expect)
		}
	}

	otherArgsExpect := make(map[string]struct{})
	otherArgsActual := make(map[string]struct{})


	for _, arg := range args[1:len(args)-2] {
		otherArgsActual[arg] = struct{}{}
	}

	for _, arg := range []string {
		"-server",
		"-Xms64M",
		"-Xmx2G",
		"-XX:MetaspaceSize=1G",
		"-XX:MaxMetaspaceSize=2G",
		"-Xmn256M",
		"-XX:MaxNewSize=256M",
		"-XX:SurvivorRatio=8",
		"-XX:TargetSurvivorRatio=50",
		"-Dcom.sun.management.jmxremote",
		"-Dcom.sun.management.jmxremote.port=20001",
		"-Dcom.sun.management.jmxremote.ssl=false",
		"-Dcom.sun.management.jmxremote.authenticate=false",
		"-Dcom.amazonaws.sdk.disableCertChecking=True",
		"-Dfile.encoding=UTF-8",
		"-Dhttp.netty.maxInitialLineLength=8192",
		"-Dhttp.port=9000",
		"-XX:+UseConcMarkSweepGC",
		"-XX:+CMSParallelRemarkEnabled",
		"-XX:+UseCMSInitiatingOccupancyOnly",
		"-XX:CMSInitiatingOccupancyFraction=70",
		"-XX:+ScavengeBeforeFullGC",
		"-XX:+CMSScavengeBeforeRemark",
		"-verbose:gc",
		"-XX:+PrintGCDateStamps",
		"-XX:+PrintGCDetails",
		"-Xloggc:/path/to/your-app/logs/gc/gc_20150910_123456.log",
		"-XX:+UseGCLogFileRotation",
		"-XX:GCLogFileSize=10M",
		"-XX:NumberOfGCLogFiles=10",
		"-XX:+HeapDumpOnOutOfMemoryError",
		"-XX:HeapDumpPath=/path/to/your-app/logs/dump",
		"-XX:ErrorFile=/path/to/your-app/logs/hs_error_pid%p.log",
	} {
		otherArgsExpect[arg] = struct{}{}
	}

	if !reflect.DeepEqual(otherArgsActual, otherArgsExpect) {
		t.Errorf("%v must equal to %v", otherArgsActual, otherArgsExpect)
	}
}
