package main

import (
	"os"
	"fmt"
	"ocr-pub/internal/cmd"
)

func main() {
	command := cmd.NewServerCommand()
	if err := command.GetCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}