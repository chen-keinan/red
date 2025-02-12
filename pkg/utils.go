package pkg

import (
	"bufio"
	"devcli/pkg/env"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"
)

func CreateOutputFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}
}

func Help() {
	fmt.Println("devcli")
	fmt.Println("Command Options:")
	fmt.Println("-- clean      Clean up resources and delete DevEnv files")
	fmt.Println("-- setup      Setting up app-proxy and gitops-operator DevEnv")
}

func GetOutputFolder() string {
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s/.devcli", usr.HomeDir)
}

func ReadInput(paramMap map[string]string, configFolder string) error {
	inputScanner := bufio.NewScanner(os.Stdin)
	count := 1
	keys := []string{"Helm Values Path", "Codefresh Namespace", "Cluster Name", "Environment Variable Script Path", "debug-app-proxy", "debug-gitops-operator"}
	fmt.Println("***************************************************************************************************************************")
	fmt.Println()
	for _, key := range keys {
		realKey := key
		if strings.Contains(key, "-") {
			realKey = strings.ReplaceAll(key, "-", "_")
		} else {
			realKey = strings.ReplaceAll(strings.ToLower(key), " ", "_")
		}
		fmt.Printf("%d. Enter %s (default: %s):", count, key, paramMap[realKey])
		for inputScanner.Scan() {
			input := inputScanner.Text()
			if input == "" && len(paramMap[realKey]) == 0 {
				fmt.Print("you must enter a value\n")
				fmt.Printf("%d. Enter %s", count, key)
				continue
			}
			if input != "" {
				paramMap[realKey] = input
			}
			break
		}
		count++
	}
	data, err := json.MarshalIndent(paramMap, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s", configFolder, env.DevCliConfigFile), data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("\n****************************************************************************************************************************")
	fmt.Println()
	return nil
}
