package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"red/pkg/env"
	"strings"
)

func CreateOutputFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetOutputFolder() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.red", usr.HomeDir), nil
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
	err = os.WriteFile(fmt.Sprintf("%s/%s", configFolder, env.RedConfigFile), data, 0644)
	if err != nil {
		return err
	}

	fmt.Println("\n****************************************************************************************************************************")
	fmt.Println()
	return nil
}
