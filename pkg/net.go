package pkg

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

func GetNgrokPublicUrl(port string, tunnelPort string) string {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ngrok http %s &", port))
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start() // Starts command asynchronously
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(time.Second * 2)
	res, err := http.Get(fmt.Sprintf("http://localhost:%s/api/tunnels", tunnelPort))
	if err != nil {
		panic(err.Error())
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var tmapOne = make(map[string]interface{})
	err = json.Unmarshal(b, &tmapOne)
	if err != nil {
		panic(err.Error())
	}
	gnroks := tmapOne["tunnels"]
	aa := gnroks.([]interface{})
	n := aa[0].(map[string]interface{})
	return n["public_url"].(string)
}

func PortForward(portInternal string, portExternal string, deploymentName string) {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl get pods -n codefresh | grep %s | awk '{print $1}'", deploymentName)).Output()
	if err != nil {
		panic(err.Error())
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	forward := fmt.Sprintf("kubectl port-forward pods/%s %s:%s -n codefresh", strings.Trim(string(out), "\n"), portInternal, portExternal)
	cmd := exec.Command("bash", "-c", forward)
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err = cmd.Start() // Starts command asynchronously
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(time.Second * 2)
}