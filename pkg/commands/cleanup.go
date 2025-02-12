package commands

import (
	"fmt"
	"os"
	"os/exec"
)

func Cleanup(folder string, silent bool) error {
	if silent {
		fmt.Println("- Clean up ngrok tunnels")
	}
	_, err := exec.Command("bash", "-c", "pgrep -f ngrok | xargs kill -9").Output()
	if err != nil {
		return err
	}
	if silent {
		fmt.Println("- Clean up port forwards")
	}
	_, err = exec.Command("bash", "-c", "pgrep -f port-forward | xargs kill -9").Output()
	if err != nil {
		return err
	}
	if silent {
		fmt.Printf("- Clean up output folder: %s\n", folder)
	}
	os.Remove(fmt.Sprintf("%s/app-proxy-dev-env.json", folder))
	os.Remove(fmt.Sprintf("%s/gitops-dev-env.json", folder))
	return nil
}
