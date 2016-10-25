package main

import (
	//"os"
	"testing"
	"fmt"
)

func tearDown() {
	fmt.Println("tear -- down")
}

func TestJavagtl(t *testing.T) {
	defer tearDown()

	//wd, _ := os.Getwd()

	fmt.Println("testing")


}
