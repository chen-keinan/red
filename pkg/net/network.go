package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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

func PortForward(portContainer string, portLocal string, deploymentName string) error {

	var stdoutBuf, stderrBuf bytes.Buffer
	forward := fmt.Sprintf("kubectl port-forward deployment/%s %s:%s -n codefresh", deploymentName, portLocal, portContainer)
	cmd := exec.Command("bash", "-c", forward)
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start() // Starts command asynchronously
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	return nil
}
