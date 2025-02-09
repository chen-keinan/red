package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"helm.sh/helm/pkg/chartutil"
)

func main() {
	inputScanner := bufio.NewScanner(os.Stdin)
	count := 1
	keys := []string{"valuesFilePath", "codefreshNamespace", "clusterName", "environmentVariableExtractScriptPath"}
	paramMap := map[string]string{
		"valuesFilePath":                       "/Users/chenkeinan/workspace/codefresh-values/sandbox.values.yaml",
		"codefreshNamespace":                   "codefresh",
		"clusterName":                          "kind-codefresh-local-cluster",
		"environmentVariableExtractScriptPath": "/Users/chenkeinan/workspace/codefresh-values/env.sh",
		"debug-app-proxy":                      "y",
		"debug-gitops-operator":                "y",
	}
	for _, key := range keys {
		fmt.Printf("%d. Enter %s (default:%s):", count, key, paramMap[key])
		for inputScanner.Scan() {
			input := inputScanner.Text()
			if input == "" && len(paramMap[key]) == 0 {
				fmt.Print("you must enter a value\n")
				fmt.Printf("%d. Enter %s", count, key)
				continue
			}
			if input != "" {
				paramMap[key] = input
			}
			break
		}
		count++
	}
	file, err := os.ReadFile(paramMap["valuesFilePath"])
	if err != nil {
		panic(err.Error())
	}
	doc, err := chartutil.ReadValues(file)
	if err != nil {
		panic(err.Error())
	}
	userToken, err := doc.Table("global.codefresh.userToken")
	if err != nil {
		panic(err.Error())
	}
	paramMap["codefreshUserToken"] = userToken["token"].(string)
	gitData, err := doc.Table("global.runtime.gitCredentials.password")
	if err != nil {
		panic(err.Error())
	}
	paramMap["gitPassword"] = gitData["value"].(string)
	AddEnvParams(paramMap)

	if paramMap["debug-app-proxy"] == "y" {
		paramMap["app-proxy-local-ip"] = getNgrokPublicUrl("3017", "4040")
	}

	if paramMap["debug-app-proxy"] == "y" {
		paramMap["gitops-operator-local-ip"] = getNgrokPublicUrl("8082", "4041")
	}
	KillProcess("gnrok")
	fmt.Println(paramMap)
}

func KillProcess(name string) {
	_, err := exec.Command("bash", "-c", "pgrep -f ngrok | xargs kill -9").Output()
	if err != nil {
		panic(err.Error())
	}
}

func trimValues(val string) string {
	sepArray := []string{`\n`, `\r`, `,`, `"`, `\`}
	for _, sep := range sepArray {
		val = strings.ReplaceAll(val, sep, "")
	}
	return strings.TrimSpace(val)
}

func getNgrokPublicUrl(port string, tunnelPort string) string {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ngrok http %s &", port))
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start() // Starts command asynchronously
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(time.Second * 2)
	res, err := http.Get(fmt.Sprintf("http://localhost:%s/api/tunnels", tunnelPort))
	if err != nil {
		panic(err.Error())
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var tmapOne = make(map[string]interface{})
	err = json.Unmarshal(b, &tmapOne)
	if err != nil {
		panic(err.Error())
	}
	gnroks := tmapOne["tunnels"]
	aa := gnroks.([]interface{})
	n := aa[0].(map[string]interface{})
	return n["public_url"].(string)
}

func AddEnvParams(envVar map[string]string) {
	cmd := fmt.Sprintf("%s %s %s", envVar["environmentVariableExtractScriptPath"], envVar["codefreshNamespace"], envVar["clusterName"])
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
