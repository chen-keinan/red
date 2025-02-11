package cluster

import (
	"fmt"
	"os/exec"
)

func PatchConfigMap(cmName string, key string, value string) {
	patch := fmt.Sprintf(`kubectl patch configmap/%s -n codefresh \--type merge -p '{"data":{"%s":"%s"}}'`, cmName, key, value)
	_, err := exec.Command("bash", "-c", patch).Output()
	if err != nil {
		panic(err.Error())
	}
}

func PatchGitOpsDeployment() {
	_, err := exec.Command("bash", "-c", "kubectl scale deployment gitops-operator -n codefresh --replicas=0").Output()
	if err != nil {
		panic(err.Error())
	}
}
