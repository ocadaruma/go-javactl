package main

import (
	//"path/filepath"
	"fmt"
	"os"
	"testing"
)

func tearDown() {
	fmt.Println("tear -- down")
}

func TestJavagtl(t *testing.T) {
	defer tearDown()

	//wd, _ := os.Getwd()

	//templateFile := filepath.Join(wd, "tests", "resources", "test_01.yml.j2")

	os.Args = append(os.Args, "abc", "def")
	main()
}
