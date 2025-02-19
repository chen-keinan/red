package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetNgrokPublicUrl(port string, tunnelPort string) (string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ngrok http %s &", port))
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start() // Starts command asynchronously
	if err != nil {
		return "", err
	}
	time.Sleep(time.Second * 2)
	res, err := http.Get(fmt.Sprintf("http://localhost:%s/api/tunnels", tunnelPort))
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var tmapOne = make(map[string]interface{})
	err = json.Unmarshal(b, &tmapOne)
	if err != nil {
		return "", err
	}
	gnroks := tmapOne["tunnels"]
	aa := gnroks.([]interface{})
	n := aa[0].(map[string]interface{})
	return n["public_url"].(string), nil
}

func PortForwardString(portContainer string, portLocal string, deploymentName string) (string, error) {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl get pods -n codefresh | grep %s | awk '{print $1}'", deploymentName)).Output()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("kubectl port-forward pod/%s %s:%s -n codefresh", strings.Trim(string(out), "\n"), portLocal, portContainer), nil

}

func PortForward(pf []string) error {
	var buffer bytes.Buffer

	for _, pfStr := range pf {
		buffer.WriteString(fmt.Sprintf("%s & ", pfStr))
	}
	_, err := exec.Command("bash", "-c", buffer.String()).Output()
	if err != nil {
		return err
	}
	return nil
}
