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

const DevCliConfigFile = "devcli.json"

func AddEnvParams(envVar map[string]string) {
	cmd := fmt.Sprintf("%s %s %s", envVar["environment_variable_script_path"], envVar["codefresh_namespace"], envVar["cluster_name"])
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic(err.Error())
	}
	if err != nil {
		panic(err.Error())
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
}

func trimValues(val string) string {
	sepArray := []string{`\n`, `\r`, `,`, `"`, `\`}
	for _, sep := range sepArray {
		val = strings.ReplaceAll(val, sep, "")
	}
	return strings.TrimSpace(val)
}

type DevCliConfig struct {
	HelmValuesPath                string `json:"helm_Values_path"`
	CodefreshNamespace            string `json:"codefresh_namespace"`
	CodefreshClusterName          string `json:"cluster_name"`
	EnvironmentVariableScriptPath string `json:"environment_variable_script_path"`
	DebugAppProxy                 string `json:"debug_app_proxy"`
	DebugGitopsOperator           string `json:"debug_gitops_operator"`
}

func LoadConfigfile(folderPath string) *DevCliConfig {
	devcli, err := os.Open(fmt.Sprintf("%s/%s", folderPath, DevCliConfigFile))
	if err != nil {
		return nil
	}
	byteValue, _ := io.ReadAll(devcli)

	// we initialize our Users array
	var devCliConfig DevCliConfig

	err = json.Unmarshal(byteValue, &devCliConfig)
	if err != nil {
		return nil
	}
	return &devCliConfig
}
