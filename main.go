package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	data := map[string]string{
		"codefresh.url":             "https://<your-platform-url>.ngrok.io",
		"codefresh.userToken.token": "",
		"codefresh.accountId":       "",
		"runtime.ingressUrl":        "http://host.docker.internal:8080",
		"service.type":              "NodePort",
		"service.nodePort":          "31243",
		"app-proxy.config.cors":     "http://local.codefresh.io,https://<your-platform-url>.ngrok.io",
	}
	keys := []string{"codefresh.url", "codefresh.userToken.token", "codefresh.accountId", "runtime.name", "runtime.ingressUrl", "service.type", "service.nodePort", "app-proxy.config.cors"}
	inputScanner := bufio.NewScanner(os.Stdin)
	count := 1
	for _, key := range keys {
		fmt.Printf("%d. Enter %s (default:%s):", count, key, data[key])
		for inputScanner.Scan() {
			val := inputScanner.Text()
			data[key] = val
			break
		}
		count++
	}
	keys = []string{"codefreshNamespace", "clusterName", "environmentVariableExtractScriptPath"}
	envVar := map[string]string{
		"codefreshNamespace":                   "codefresh",
		"clusterName":                          "kind-codefresh-local-cluster",
		"environmentVariableExtractScriptPath": "/Users/chenkeinan/workspace/codefresh-values/env.sh",
	}
	for _, key := range keys {
		fmt.Printf("%d. Enter %s :", count, key)
		for inputScanner.Scan() {
			input := inputScanner.Text()
			if input == "" && len(envVar[key]) == 0 {
				fmt.Print("you must enter a value\n")
				fmt.Printf("%d. Enter %s", count, key)
				continue
			}
			if input != "" {
				envVar[key] = input
			}
			break
		}
		count++
	}
	cmd := fmt.Sprintf("%s %s %s", envVar["environmentVariableExtractScriptPath"], envVar["codefreshNamespace"], envVar["clusterName"])
	fmt.Print(cmd)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic("some error found")
	}
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(0)
	}
	file := strings.NewReader(string(out))
	scanner := bufio.NewScanner(file)
	fmt.Println(string(out))
}
