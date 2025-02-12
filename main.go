package main

import (
	"devcli/pkg"
	"devcli/pkg/commands"
	"fmt"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
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

	if len(argsWithoutProg) > 0 {
		switch argsWithoutProg[0] {
		case "--clean":
			err = commands.Cleanup(outputFolder, true)
			if err != nil {
				fmt.Println(fmt.Errorf("error: failed to create output folder: %w", err))
				os.Exit(1)
			}
			return
		case "--setup":
			err = commands.Cleanup(outputFolder, false)
			if err != nil {
				fmt.Println(fmt.Errorf("error: failed to cleanup resources: %w", err))
				os.Exit(1)
			}
			err = commands.Setup(outputFolder)
			if err != nil {
				fmt.Println(fmt.Errorf("error: failed to setup dev env %w", err))
				os.Exit(1)
			}
		default:
			pkg.Help()
		}
	} else {
		pkg.Help()
	}
}
