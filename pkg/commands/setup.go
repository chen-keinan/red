package commands

import (
	"devcli/pkg"
	"devcli/pkg/cluster"
	"devcli/pkg/env"
	"devcli/pkg/net"
	"fmt"
)

func Setup(outputFolder string) {
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
	pkg.ReadInput(paramMap, outputFolder)
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
