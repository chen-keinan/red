package cluster

import (
	"fmt"
	"os/exec"
)

func PatchConfigMap(cmName string, key string, value string) error {
	patch := fmt.Sprintf(`kubectl patch configmap/%s -n codefresh \--type merge -p '{"data":{"%s":"%s"}}'`, cmName, key, value)
	_, err := exec.Command("bash", "-c", patch).Output()
	if err != nil {
		return err
	}
	return nil
}

func PatchGitOpsOperatorAppProxyEnvVar(apUrl string) error {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl set env deployment/gitops-operator AP_URL=%s", apUrl)).Output()
	if err != nil {
		return err
	}
	return nil
}
