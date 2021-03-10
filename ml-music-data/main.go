package main

import (
	"fmt"
	"github.com/sta-golang/go-lib-utils/cmd"
)

func main() {
	command, err := cmd.ExecCmd("python","./music_download/crawler.py")
	fmt.Println(err)
	fmt.Println(command.OutInfo())
}