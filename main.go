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
	"os/user"
	"strings"
	"time"

	"helm.sh/helm/pkg/chartutil"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "cleanup" {
		cleanup()
		return
	}
	cleanup()
	paramMap := map[string]string{
		"valuesFilePath":                       "/Users/chenkeinan/workspace/codefresh-values/sandbox.values.yaml",
		"codefreshNamespace":                   "codefresh",
		"clusterName":                          "kind-codefresh-local-cluster",
		"environmentVariableExtractScriptPath": "/Users/chenkeinan/workspace/codefresh-values/env.sh",
		"debug-app-proxy":                      "y",
		"debug-gitops-operator":                "y",
	}
	// read user input
	readInput(paramMap)
	// add params from values yaml
	AddHelmValues(paramMap)
	// add params from envVar
	AddEnvParams(paramMap)
	var argoServerPortForward bool
	if paramMap["debug-app-proxy"] == "y" {
		fmt.Println("setting ngrok 3017")
		paramMap["app-proxy-local-ip"] = getNgrokPublicUrl("3017", "4040")
		portForward("2746", "2746", "argo-server")
		portForward("8080", "8081", "argo-cd-server")
		patchConfigMap("codefresh-cm", "ingressHost", paramMap["app-proxy-local-ip"])
		argoServerPortForward = true
	}

	if paramMap["debug-gitops-operator"] == "y" {
		fmt.Println("setting ngrok 8082")
		paramMap["gitops-operator-local-ip"] = getNgrokPublicUrl("8082", "4041")
		if !argoServerPortForward {
			fmt.Println("setting port forward 2746")
			portForward("2746", "2746", "argo-server")
		}
		fmt.Println("scalling gitops operator down to 0")
		patchGitOpsDeployment()
		fmt.Println("updating gitops-operator-notifications cm with gitops local dev ip")
		patchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-degraded-notifier", fmt.Sprintf("url: %s/app-degraded", paramMap["gitops-operator-local-ip"]))
		patchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-revision-changed-notifier", fmt.Sprintf("url: %s/app-revision-changed", paramMap["gitops-operator-local-ip"]))
	}
	outputFolder := getOutputFolder()
	createOutputFolder(outputFolder)
	generateEnvVarForGitOpsOpertorDev(paramMap, outputFolder)
	generateEnvVarForAppProxyDev(paramMap, outputFolder)
}

func getOutputFolder() string {
	usr, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s/dev-output", usr.HomeDir)
}

func createOutputFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}
}

func generateEnvVarForGitOpsOpertorDev(paramMap map[string]string, outputFolder string) {
	gitOpsOperatorMap := map[string]string{
		"AP_URL":                    "<app-proxy-local-ip>",
		"ARGO_CD_URL":               "localhost:8081",
		"ARGO_WF_URL":               "http://localhost:2746",
		"ARGO_WF_TOKEN":             "ARGO_WF_TOKEN",
		"CF_TOKEN":                  "RUNTIME_TOKEN",
		"CF_URL":                    "CF_HOST",
		"GITOPS_OPERATOR_VERSION":   "v0.3.17",
		"HEALTH_PROBE_BIND_ADDRESS": ":8081",
		"LEADER_ELECT":              "false",
		"METRICS_BIND_ADDRESS":      "127.0.0.1: 8085",
		"RUNTIME":                   "codefresh",
	}
	gitOpsOperatorMap["AP_URL"] = paramMap["app-proxy-local-ip"]
	gitOpsOperatorMap["ARGO_WF_TOKEN"] = paramMap["ARGO_WF_TOKEN"]
	gitOpsOperatorMap["CF_TOKEN"] = paramMap["RUNTIME_TOKEN"]
	gitOpsOperatorMap["CF_URL"] = paramMap["CF_HOST"]

	gitOpsOperatorData, err := json.MarshalIndent(gitOpsOperatorMap, "", "    ")
	if err != nil {
		panic(err.Error())
	}
	err = os.WriteFile(fmt.Sprintf("%s/gitops-dev-env.json", outputFolder), gitOpsOperatorData, 0755)
	if err != nil {
		panic(err.Error())
	}
}

func generateEnvVarForAppProxyDev(paramMap map[string]string, outputFolder string) {
	appProxyMap := map[string]string{
		"NODE_TLS_REJECT_UNAUTHORIZED": "0",
		"ARGO_CD_URL":                  "http://localhost:8089",
		"ARGO_WORKFLOWS_URL":           "http://localhost:2746",
		"GIT_PASSWORD":                 "gitPassword",
		"ARGO_CD_PASSWORD":             "ARGO_CD_PASSWORD",
		"ARGO_CD_USERNAME":             "admin",
		"ARGO_WORKFLOWS_INSECURE":      "true",
		"ARGO_WORKFLOWS_SA_TOKEN":      "ARGO_WORKFLOWS_SA_TOKEN",
		"CF_HOST":                      "CF_HOST",
		"CHART_VERSION":                "0.15.0",
		"CLUSTER":                      "https://kubernetes.default.svc",
		"CORS":                         "cors",
		"DEPLOYMENT_NAME":              "cap-app-proxy",
		"HELM_RELEASE_NAME":            "cf-gitops-runtime",
		"INSTALLATION_TYPE":            "HELM",
		"LOG_LEVEL":                    "info",
		"MANAGED":                      "false",
		"NAMESPACE":                    "codefresh",
		"PART_OF_VALUE":                "app-proxy",
		"ROLLOUTS_HELM_REPOSITORY":     "https://codefresh-io.github.io/argo-helm",
		"ROLLOUTS_HELM_VERSION":        "2.37.3-1-v1.7.1-CR-24605",
		"RUNTIME_NAME":                 "gitlab-runtime",
		"RUNTIME_STORE_IV":             "RUNTIME_STORE_IV",
		"RUNTIME_TOKEN":                "RUNTIME_TOKEN",
		"RUNTIME_VERSION":              "0.1.65",
		"SKIP_PERMISSIONS_VALIDATION":  "false",
		"USER_TOKEN":                   "USER_TOKEN",
	}
	appProxyMap["GIT_PASSWORD"] = paramMap["gitPassword"]
	appProxyMap["ARGO_CD_USERNAME"] = paramMap["ARGO_CD_USERNAME"]
	appProxyMap["ARGO_CD_PASSWORD"] = paramMap["ARGO_CD_PASSWORD"]
	appProxyMap["ARGO_WORKFLOWS_SA_TOKEN"] = paramMap["ARGO_WORKFLOWS_SA_TOKEN"]
	appProxyMap["CF_HOST"] = paramMap["CF_HOST"]
	appProxyMap["RUNTIME_NAME"] = paramMap["RUNTIME_NAME"]
	appProxyMap["RUNTIME_STORE_IV"] = paramMap["RUNTIME_STORE_IV"]
	appProxyMap["RUNTIME_TOKEN"] = paramMap["RUNTIME_TOKEN"]
	appProxyMap["USER_TOKEN"] = paramMap["USER_TOKEN"]
	appProxyMap["CORS"] = paramMap["cors"]

	gitAppProxyData, err := json.MarshalIndent(appProxyMap, "", "    ")
	if err != nil {
		panic(err.Error())
	}
	err = os.WriteFile(fmt.Sprintf("%s/app-proxy-dev-env.json", outputFolder), gitAppProxyData, 0755)
	if err != nil {
		panic(err.Error())
	}
}

func patchConfigMap(cmName string, key string, value string) {
	patch := fmt.Sprintf(`kubectl patch configmap/%s -n codefresh \--type merge -p '{"data":{"%s":"%s"}}'`, cmName, key, value)
	_, err := exec.Command("bash", "-c", patch).Output()
	if err != nil {
		panic(err.Error())
	}
}

func patchGitOpsDeployment() {
	_, err := exec.Command("bash", "-c", "kubectl scale deployment gitops-operator -n codefresh --replicas=0").Output()
	if err != nil {
		panic(err.Error())
	}
}

func cleanup() {
	_, err := exec.Command("bash", "-c", "pgrep -f ngrok | xargs kill -9").Output()
	if err != nil {
		panic(err.Error())
	}
	_, err = exec.Command("bash", "-c", "pgrep -f port-forward | xargs kill -9").Output()
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

func AddHelmValues(paramMap map[string]string) {
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
	corsData, err := doc.Table("app-proxy.config")
	if err != nil {
		panic(err.Error())
	}
	paramMap["cors"] = corsData["cors"].(string)
	hostData, err := doc.Table("global.codefresh")
	if err != nil {
		panic(err.Error())
	}
	paramMap["CF_HOST"] = hostData["url"].(string)
}

func readInput(paramMap map[string]string) {
	inputScanner := bufio.NewScanner(os.Stdin)
	count := 1
	keys := []string{"valuesFilePath", "codefreshNamespace", "clusterName", "environmentVariableExtractScriptPath"}
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

}

func portForward(portInternal string, portExternal string, deploymentName string) {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl get pods -n codefresh | grep %s | awk '{print $1}'", deploymentName)).Output()
	if err != nil {
		panic(err.Error())
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	forward := fmt.Sprintf("kubectl port-forward pods/%s %s:%s -n codefresh", strings.Trim(string(out), "\n"), portInternal, portExternal)
	cmd := exec.Command("bash", "-c", forward)
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start() // Starts command asynchronously
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(time.Second * 2)
}
