package env

import (
	"encoding/json"
	"fmt"
	"os"
)

func GenerateEnvVarForAppProxyDev(paramMap map[string]string, outputFolder string) error {
	appProxyMap := map[string]string{
		"NODE_TLS_REJECT_UNAUTHORIZED": "0",
		"ARGO_CD_URL":                  "http://localhost:8080",
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
		return err
	}
	filePath := fmt.Sprintf("%s/app-proxy-dev-env.json", outputFolder)
	err = os.WriteFile(filePath, gitAppProxyData, 0755)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", filePath)
	return nil
}

func GenerateEnvVarForGitOpsOpertorDev(paramMap map[string]string, outputFolder string) error {
	gitOpsOperatorMap := map[string]string{
		"AP_URL":                    "<app-proxy-local-ip>",
		"ARGO_CD_URL":               "localhost:8080",
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
	gitOpsOperatorMap["ARGO_WF_TOKEN"] = paramMap["ARGO_WORKFLOWS_SA_TOKEN"]
	gitOpsOperatorMap["CF_TOKEN"] = paramMap["RUNTIME_TOKEN"]
	gitOpsOperatorMap["CF_URL"] = paramMap["CF_HOST"]

	gitOpsOperatorData, err := json.MarshalIndent(gitOpsOperatorMap, "", "    ")
	if err != nil {
		return err
	}
	filePath := fmt.Sprintf("%s/gitops-dev-env.json", outputFolder)
	err = os.WriteFile(filePath, gitOpsOperatorData, 0755)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", filePath)
	return nil
}
