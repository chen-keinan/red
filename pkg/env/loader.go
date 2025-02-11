package env

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func AddEnvParams(envVar map[string]string) {
	cmd := fmt.Sprintf("%s %s %s", envVar["Environment Variable Extract Script Path"], envVar["Codefresh Namespace"], envVar["Cluster Name"])
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
