package commands

import "fmt"

func Help() {
	fmt.Println("RuntimeEnvDev")
	fmt.Println("Command Options:")
	fmt.Println("--clean      Clean up resources and delete DevEnv files")
	fmt.Println("--setup      Setting up app-proxy and gitops-operator DevEnv")
	fmt.Println("--no-setup   Loading setup from red.json (not valid on 1st setup)")
	fmt.Println("--help       Show avaliable command options")
}
