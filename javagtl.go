package main

import (
	"fmt"
	"os"

	"github.com/ocadaruma/javagtl/setting"
)

func main() {
	fmt.Println("WIP")
	setting, err := setting.LoadSetting(os.Args[1])

	fmt.Printf("setting: %v", setting)
	fmt.Printf("error: %v", err)
}
