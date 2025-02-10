package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

func CreateOutputFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}
}

func Cleanup(folder string, silent bool) {
	if silent {
		fmt.Println("- Clean up ngrok tunnels")
	}
	_, err := exec.Command("bash", "-c", "pgrep -f ngrok | xargs kill -9").Output()
	if err != nil {
		panic(err.Error())
	}
	if silent {
		fmt.Println("- Clean up port forwards")
	}
	_, err = exec.Command("bash", "-c", "pgrep -f port-forward | xargs kill -9").Output()
	if err != nil {
		panic(err.Error())
	}
	if silent {
		fmt.Printf("- Clean up output folder: %s\n", folder)
	}
	os.Remove(fmt.Sprintf("%s/app-proxy-dev-env.json", folder))
	os.Remove(fmt.Sprintf("%s/gitops-dev-env.json", folder))
}

func Help() {
	fmt.Println("cf-cli")
	fmt.Println("Command Options:")
	fmt.Println("-- clean      Clean up resources and delete DevEnv files")
	fmt.Println("-- setup      Setting up app-proxy and gitops-operator DevEnv")
}

func GetOutputFolder() string {
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s/.cf-cli", usr.HomeDir)
}