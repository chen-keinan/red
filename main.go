package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"os"
	"runtime-cli/pkg"
	"runtime-cli/pkg/cluster"
	"runtime-cli/pkg/env"
	"runtime-cli/pkg/net"
)

func main() {
	argsWithoutProg := os.Args[1:]
	outputFolder := pkg.GetOutputFolder()
	pkg.CreateOutputFolder(outputFolder)

	if len(argsWithoutProg) > 0 {
		switch argsWithoutProg[0] {
		case "--clean":
			pkg.Cleanup(outputFolder, true)
			return
		case "--setup":
			pkg.Cleanup(outputFolder, false)
			setup(outputFolder)
		default:
			pkg.Help()
		}
	} else {
		pkg.Help()
	}
}

func setup(outputFolder string) {
	paramMap := map[string]string{
		"helm_values_path":                 "/Users/chenkeinan/workspace/codefresh-values/local.values.yaml",
		"codefresh_namespace":              "codefresh",
		"cluster_name":                     "kind-codefresh-local-cluster",
		"environment_variable_script_path": "/Users/chenkeinan/workspace/codefresh-values/env.sh",
		"debug_app_proxy":                  "y",
		"debug_gitops_operator":            "y",
	}

	config := env.LoadConfigfile(outputFolder)
	if config != nil {
		paramMap["helm_Values_path"] = config.HelmValuesPath
		paramMap["codefresh_namespace"] = config.CodefreshNamespace
		paramMap["cluster_name"] = config.CodefreshClusterName
		paramMap["environment_variable_script_path"] = config.EnvironmentVariableScriptPath
		paramMap["debug_app_proxy"] = config.DebugAppProxy
		paramMap["debug_gitops_operator"] = config.DebugGitopsOperator
	}

	// read user input
	readInput(paramMap, outputFolder)
	// add params from values yaml
	fmt.Println("- Reading Helm Values")
	cluster.AddHelmValues(paramMap)
	// add params from envVar
	fmt.Println("- Extracting Values from EnvVar script")
	env.AddEnvParams(paramMap)
	var argoServerPortForward bool
	net.GetNgrokPublicUrl("2020", "4040")
	if paramMap["debug_app_proxy"] == "y" {
		fmt.Println("- Tunneling 3017 --> Localhost")
		paramMap["app-proxy-local-ip"] = net.GetNgrokPublicUrl("3017", "4041")
		net.PortForward("2746", "2746", "argo-server")
		net.PortForward("8080", "8080", "argo-cd-server")
		cluster.PatchConfigMap("codefresh-cm", "ingressHost", paramMap["app-proxy-local-ip"])
		argoServerPortForward = true
	}

	if paramMap["debug_gitops_operator"] == "y" {
		fmt.Println("- Tunneling 8082 --> Localhost")
		paramMap["gitops-operator-local-ip"] = net.GetNgrokPublicUrl("8082", "4042")
		if !argoServerPortForward {
			net.PortForward("2746", "2746", "argo-server")
		}
		fmt.Println("- Scalling down gitops operator to 0")
		cluster.PatchGitOpsDeployment()
		fmt.Println("- Updating gitops-operator-notifications cm with gitops local dev ip")
		cluster.PatchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-degraded-notifier", fmt.Sprintf("url: %s/app-degraded\\nheaders:\\n- name: Content-Type\\n  value: application/json\\n", paramMap["gitops-operator-local-ip"]))
		cluster.PatchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-revision-changed-notifier", fmt.Sprintf("url: %s/app-revision-changed\\nheaders:\\n- name: Content-Type\\n  value: application/json\\n", paramMap["gitops-operator-local-ip"]))
	}
	if paramMap["debug_app_proxy"] == "y" || paramMap["debug_gitops_operator"] == "y" {
		pkg.CreateOutputFolder(outputFolder)
		fmt.Println("********************************************************")
		fmt.Println("-- output files:")

		if paramMap["debug_app_proxy"] == "y" {
			env.GenerateEnvVarForAppProxyDev(paramMap, outputFolder)
		}
		if paramMap["debug_gitops_operator"] == "y" {
			env.GenerateEnvVarForGitOpsOpertorDev(paramMap, outputFolder)
		}
		fmt.Println("\n******************************************************")
	}
}

func readInput(paramMap map[string]string, configFolder string) error {
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
		fmt.Printf("%d. Enter %s (default:%s):", count, key, paramMap[realKey])
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
