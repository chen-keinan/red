package main

import (
	"fmt"
	"os"
	"red/pkg"
	"red/pkg/commands"
	"red/pkg/env"
)

func main() {
	commandsParam := os.Args[1:]
	outputFolder, err := pkg.GetOutputFolder()
	if err != nil {
		fmt.Println(fmt.Errorf("error: failed to get output folder name: %w", err))
		os.Exit(1)
	}
	pkg.CreateOutputFolder(outputFolder)
	if err != nil {
		fmt.Println(fmt.Errorf("error: failed to create output folder: %w", err))
		os.Exit(1)
	}
	if len(commandsParam) == 0 {
		commandsParam = append(commandsParam, "no-command")
	}

	switch commandsParam[0] {
	case "--clean":
		err = commands.Cleanup(outputFolder, true)
		if err != nil {
			fmt.Println(fmt.Errorf("error: failed to cleanup resource: %w", err))
			os.Exit(1)
		}
		return
	case "--setup":
		err = commands.Cleanup(outputFolder, false)
		if err != nil {
			fmt.Println(fmt.Errorf("error: failed to cleanup resources: %w", err))
			os.Exit(1)
		}
		err = commands.Setup(outputFolder, true)
		if err != nil {
			fmt.Println(fmt.Errorf("error: failed to setup dev env %w", err))
			os.Exit(1)
		}
	case "--no-setup":
		if !env.ConfigFileExist(outputFolder) {
			fmt.Println("red.config file do not exist yet, you must run: `red --setup` on first time")
			return
		}
		err = commands.Cleanup(outputFolder, false)
		if err != nil {
			fmt.Println(fmt.Errorf("error: failed to cleanup resources: %w", err))
			os.Exit(1)
		}
		err = commands.Setup(outputFolder, false)
		if err != nil {
			fmt.Println(fmt.Errorf("error: failed to setup dev env %w", err))
			os.Exit(1)
		}
	default:
		pkg.Help()
	}
}
