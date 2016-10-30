package util

import (
	"testing"
)

func TestNewMemSize(t *testing.T) {
	testCases := []struct {
		actual, expect MemSize
	} {
		{*NewMemSize(0), MemSize{0, 1}},
		{*NewMemSize(123), MemSize{123, 1}},
		{*NewMemSize("0"), MemSize{0, 1}},
		{*NewMemSize("123"), MemSize{123, 1}},
		{*NewMemSize("123k"), MemSize{123, 1 << 10}},
		{*NewMemSize("123m"), MemSize{123, 1 << 20}},
		{*NewMemSize("123g"), MemSize{123, 1 << 30}},
		{*NewMemSize("123G"), MemSize{123, 1 << 30}},
	}

	for i, c := range testCases {
		if c.actual != c.expect {
			t.Errorf("case %v : %v must equal to %v", i, c.actual, c.expect)
		}
	}

	invalids := []*MemSize {
		NewMemSize(""),
		NewMemSize("1.23"),
		NewMemSize("p"),
		NewMemSize("gg"),
		NewMemSize("1.23k"),
	}

	for i, c := range invalids {
		if c != nil {
			t.Errorf("case %v : %v must be nil", i, c)
		}
	}
}

func TestMemSize_Bytes(t *testing.T) {
	testCases := []struct {
		actual, expect int64
	} {
		{NewMemSize("0m").Bytes(), 0},
		{NewMemSize(123).Bytes(), 123},
		{NewMemSize("123").Bytes(), 123},
		{NewMemSize("10k").Bytes(), 10240},
		{NewMemSize("10m").Bytes(), 10 * 1024 * 1024},
		{NewMemSize("10G").Bytes(), 10 * 1024 * 1024 * 1024},
	}

	for i, c := range testCases {
		if c.actual != c.expect {
			t.Errorf("case %v : %v must equal to %v", i, c.actual, c.expect)
		}
	}
}

func TestEmptyIfZero(t *testing.T) {
	testCases := []struct {
		actual, expect string
	} {
		{EmptyIfZero("foo:%s", ""), ""},
		{EmptyIfZero("foo:%s", "bar"), "foo:bar"},
		{EmptyIfZero("foo:%d", 0), ""},
		{EmptyIfZero("foo:%d", 1), "foo:1"},
		{EmptyIfZero("foo:%d", -1), "foo:-1"},
		{EmptyIfZero("foo:%v", (*string)(nil)), ""},
	}

	for i, c := range testCases {
		if c.actual != c.expect {
			t.Errorf("case %v : %v must equal to %v", i, c.actual, c.expect)
		}
	}
}

func TestEmptyIfNil(t *testing.T) {
	str := "bar"
	strPtr := &str

	i := 42
	iPtr := &i

	b := false
	bPtr := &b

	testCases := []struct {
		actual, expect string
	} {
		{EmptyIfNilString("foo:%s", (*string)(nil)), ""},
		{EmptyIfNilString("foo:%s", strPtr), "foo:bar"},
		{EmptyIfNilInt("foo:%d", (*int)(nil)), ""},
		{EmptyIfNilInt("foo:%d", iPtr), "foo:42"},
		{EmptyIfNilBool("foo:%t", (*bool)(nil)), ""},
		{EmptyIfNilBool("foo:%t", bPtr), "foo:false"},
	}

	for i, c := range testCases {
		if c.actual != c.expect {
			t.Errorf("case %v : %v must equal to %v", i, c.actual, c.expect)
		}
	}
}
