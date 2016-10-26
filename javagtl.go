package main

import (
	"fmt"
	"os"

	"github.com/ocadaruma/javagtl/setting"
)

func main() {
	var err error
	fmt.Println("WIP")

	var s *setting.Setting
	s, err = setting.LoadSetting(os.Args[1])

	fmt.Printf("setting: %v\n", s)
	fmt.Printf("error: %v\n", err)

	fmt.Printf("args: %v", os.Args)

}
