package commands

import (
	"fmt"
	"os"
	"os/exec"
	"red/pkg/cluster"
	"red/pkg/env"
)

func Cleanup(folder string, notSilent bool) error {
	config, err := env.LoadConfigfile(folder)
	if err != nil {
		return err
	}
	if config == nil {
		return nil
	}
	ingressUrl, err := cluster.GetIngressUrl(config.HelmValuesPath)
	if err != nil {
		return err
	}
	if config.DebugAppProxy == "y" {
		if notSilent {
			fmt.Println("- Revert codefresh-cm configmap")
		}
		err = cluster.PatchConfigMap("codefresh-cm", "ingressHost", ingressUrl)
		if err != nil {
			return err
		}
		if config.DebugGitopsOperator == "n" {
			err = cluster.PatchGitOpsOperatorAppProxyEnvVar("http://cap-app-proxy:3017")
			if err != nil {
				return err
			}
		}

	}
	if config.DebugGitopsOperator == "y" {
		if notSilent {
			fmt.Println("- Revert gitops-operator-notifications-cm configmap")
		}
		gitOpsOperatorUrl := "http://localhost:8082"
		err = cluster.PatchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-degraded-notifier", fmt.Sprintf("url: %s/app-degraded\\nheaders:\\n- name: Content-Type\\n  value: application/json\\n", gitOpsOperatorUrl))
		if err != nil {
			return err
		}
		err = cluster.PatchConfigMap("gitops-operator-notifications-cm", "service.webhook.cf-promotion-app-revision-changed-notifier", fmt.Sprintf("url: %s/app-revision-changed\\nheaders:\\n- name: Content-Type\\n  value: application/json\\n", gitOpsOperatorUrl))
		if err != nil {
			return err
		}
	}
	if notSilent {
		fmt.Println("- Clean up ngrok tunnels")
	}
	_, err = exec.Command("bash", "-c", "pgrep -f ngrok | xargs kill -9").Output()
	if err != nil {
		return err
	}
	if notSilent {
		fmt.Println("- Clean up port forwards")
	}
	_, err = exec.Command("bash", "-c", "pgrep -f port-forward | xargs kill -9").Output()
	if err != nil {
		return err
	}
	if notSilent {
		fmt.Printf("- Clean up output folder: %s\n", folder)
	}
	os.Remove(fmt.Sprintf("%s/app-proxy-dev-env.json", folder))
	os.Remove(fmt.Sprintf("%s/gitops-operator-dev-env.json", folder))
	return nil
}
