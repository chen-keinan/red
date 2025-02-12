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

func PortForward(portInternal string, portExternal string, deploymentName string) error {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl get pods -n codefresh | grep %s | awk '{print $1}'", deploymentName)).Output()
	if err != nil {
		return err
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	forward := fmt.Sprintf("kubectl port-forward pods/%s %s:%s -n codefresh", strings.Trim(string(out), "\n"), portInternal, portExternal)
	cmd := exec.Command("bash", "-c", forward)
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start() // Starts command asynchronously
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	return nil
}
