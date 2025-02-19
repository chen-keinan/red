package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const RedConfigFile = "red.json"

func AddEnvParams(envVar map[string]string) error {
	cmd := fmt.Sprintf("%s %s %s", envVar["environment_variable_script_path"], envVar["codefresh_namespace"], envVar["cluster_name"])
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	env := string(out)
	envScanner := bufio.NewScanner(strings.NewReader(env))
	for envScanner.Scan() {
		input := envScanner.Text()
		values := strings.Split(input, ":")
		if len(values) == 2 {
			envVar[trimValues(values[0])] = trimValues(values[1])
		}
	}
	return nil
}

func trimValues(val string) string {
	sepArray := []string{`\n`, `\r`, `,`, `"`, `\`}
	for _, sep := range sepArray {
		val = strings.ReplaceAll(val, sep, "")
	}
	return strings.TrimSpace(val)
}

type RedConfig struct {
	HelmValuesPath                string `json:"helm_Values_path"`
	CodefreshNamespace            string `json:"codefresh_namespace"`
	CodefreshClusterName          string `json:"cluster_name"`
	EnvironmentVariableScriptPath string `json:"environment_variable_script_path"`
	DebugAppProxy                 string `json:"debug_app_proxy"`
	DebugGitopsOperator           string `json:"debug_gitops_operator"`
}

func LoadConfigfile(folderPath string) (*RedConfig, error) {
	red, err := os.Open(fmt.Sprintf("%s/%s", folderPath, RedConfigFile))
	if err != nil {
		return nil, nil
	}
	byteValue, err := io.ReadAll(red)
	if err != nil {
		return nil, err
	}

	var redConfig RedConfig

	err = json.Unmarshal(byteValue, &redConfig)
	if err != nil {
		return nil, err
	}
	return &redConfig, nil
}

func ConfigFileExist(folderPath string) bool {
	_, err := os.Open(fmt.Sprintf("%s/%s", folderPath, RedConfigFile))
	return err == nil
}
