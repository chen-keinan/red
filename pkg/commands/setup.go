package commands

import (
	"devcli/pkg"
	"devcli/pkg/cluster"
	"devcli/pkg/env"
	"devcli/pkg/net"
	"fmt"
)

func Setup(outputFolder string) error {
	paramMap := map[string]string{
		"helm_values_path":                 "/helm_values/file/path",
		"codefresh_namespace":              "codefresh",
		"cluster_name":                     "clusterName",
		"environment_variable_script_path": "/env/shell/script/file/path",
		"debug_app_proxy":                  "y",
		"debug_gitops_operator":            "y",
	}

	config := env.LoadConfigfile(outputFolder)
	if config != nil {
		paramMap["helm_values_path"] = config.HelmValuesPath
		paramMap["codefresh_namespace"] = config.CodefreshNamespace
		paramMap["cluster_name"] = config.CodefreshClusterName
		paramMap["environment_variable_script_path"] = config.EnvironmentVariableScriptPath
		paramMap["debug_app_proxy"] = config.DebugAppProxy
		paramMap["debug_gitops_operator"] = config.DebugGitopsOperator
	}

	// read user input
	err := pkg.ReadInput(paramMap, outputFolder)
	if err != nil {
		return err
	}
	// add params from values yaml
	fmt.Println("- Reading Helm Values")
	err = cluster.AddHelmValues(paramMap)
	if err != nil {
		return err
	}
	// add params from envVar
	fmt.Println("- Extracting Values from EnvVar script")
	err = env.AddEnvParams(paramMap)
	if err != nil {
		return err
	}
	var argoServerPortForward bool
	_, err = net.GetNgrokPublicUrl("2020", "4040")
	if err != nil {
		return err
	}
	if paramMap["debug_app_proxy"] == "y" {
		fmt.Println("- Tunneling 3017 --> Localhost")
		appProxyLocalIp, err := net.GetNgrokPublicUrl("3017", "4041")
		if err != nil {
			return err
		}
		paramMap["app-proxy-local-ip"] = appProxyLocalIp
		err = net.PortForward("2746", "2746", "argo-server")
		if err != nil {
			return err
		}
		err = net.PortForward("8080", "8080", "argo-cd-server")
		if err != nil {
			return err
		}
		fmt.Println("- Updating codefresh-cm")
		err = cluster.PatchConfigMap("codefresh-cm", "ingressHost", paramMap["app-proxy-local-ip"])
		if err != nil {
			return err
		}
		argoServerPortForward = true
	}

	if paramMap["debug_gitops_operator"] == "y" {
		fmt.Println("- Tunneling 8082 --> Localhost")
		gitopsOperatorLocalIp, err := net.GetNgrokPublicUrl("8082", "4042")
		if err != nil {
			return err
		}
		paramMap["gitops-operator-local-ip"] = gitopsOperatorLocalIp
		if !argoServerPortForward {
			err = net.PortForward("2746", "2746", "argo-server")
			if err != nil {
				return err
			}
		}
		fmt.Println("- Scalling down gitops operator to 0")
		err = cluster.PatchGitOpsDeployment()
		if err != nil {
			return err
		}
		fmt.Println("- Updating gitops-operator-notifications cm")
		err = cluster.PatchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-degraded-notifier", fmt.Sprintf("url: %s/app-degraded\\nheaders:\\n- name: Content-Type\\n  value: application/json\\n", paramMap["gitops-operator-local-ip"]))
		if err != nil {
			return err
		}
		err = cluster.PatchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-revision-changed-notifier", fmt.Sprintf("url: %s/app-revision-changed\\nheaders:\\n- name: Content-Type\\n  value: application/json\\n", paramMap["gitops-operator-local-ip"]))
		if err != nil {
			return err
		}
	}
	if paramMap["debug_app_proxy"] == "y" || paramMap["debug_gitops_operator"] == "y" {
		fmt.Println("********************************************************")
		fmt.Println("-- output files:")

		if paramMap["debug_app_proxy"] == "y" {
			err := env.GenerateEnvVarForAppProxyDev(paramMap, outputFolder)
			if err != nil {
				return err
			}
		}
		if paramMap["debug_gitops_operator"] == "y" {
			err := env.GenerateEnvVarForGitOpsOpertorDev(paramMap, outputFolder)
			if err != nil {
				return err
			}
		}
		fmt.Println("\n******************************************************")
	}
	return nil
}
