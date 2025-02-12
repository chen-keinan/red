package main

import (
	"devcli/pkg"
	"devcli/pkg/commands"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	outputFolder := pkg.GetOutputFolder()
	pkg.CreateOutputFolder(outputFolder)

	if len(argsWithoutProg) > 0 {
		switch argsWithoutProg[0] {
		case "--clean":
			commands.Cleanup(outputFolder, true)
			return
		case "--setup":
			commands.Cleanup(outputFolder, false)
			commands.Setup(outputFolder)
		default:
			pkg.Help()
		}
	} else {
		pkg.Help()
	}
}
