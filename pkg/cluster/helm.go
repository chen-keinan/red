package cluster

import (
	"os"

	"k8s.io/helm/pkg/chartutil"
)

func AddHelmValues(paramMap map[string]string) error {
	file, err := os.ReadFile(paramMap["helm_values_path"])
	if err != nil {
		return err
	}
	doc, err := chartutil.ReadValues(file)
	if err != nil {
		return err
	}
	userToken, err := doc.Table("global.codefresh.userToken")
	if err != nil {
		return err
	}
	paramMap["codefreshUserToken"] = userToken["token"].(string)
	gitData, err := doc.Table("global.runtime.gitCredentials.password")
	if err != nil {
		return err
	}
	paramMap["gitPassword"] = gitData["value"].(string)
	corsData, err := doc.Table("app-proxy.config")
	if err != nil {
		return err
	}
	paramMap["cors"] = corsData["cors"].(string)
	hostData, err := doc.Table("global.codefresh")
	if err != nil {
		return err
	}
	paramMap["CF_HOST"] = hostData["url"].(string)
	return nil
}

func GetIngressUrl(filePath string) (string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	doc, err := chartutil.ReadValues(file)
	if err != nil {
		return "", err
	}
	runtime, err := doc.Table("global.runtime")
	if err != nil {
		return "", err
	}
	return runtime["ingressUrl"].(string), nil
}
