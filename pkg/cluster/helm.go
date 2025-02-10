package cluster

import (
	"os"

	"k8s.io/helm/pkg/chartutil"
)

func AddHelmValues(paramMap map[string]string) {
	file, err := os.ReadFile(paramMap["Helm Values Path"])
	if err != nil {
		panic(err.Error())
	}
	doc, err := chartutil.ReadValues(file)
	if err != nil {
		panic(err.Error())
	}
	userToken, err := doc.Table("global.codefresh.userToken")
	if err != nil {
		panic(err.Error())
	}
	paramMap["codefreshUserToken"] = userToken["token"].(string)
	gitData, err := doc.Table("global.runtime.gitCredentials.password")
	if err != nil {
		panic(err.Error())
	}
	paramMap["gitPassword"] = gitData["value"].(string)
	corsData, err := doc.Table("app-proxy.config")
	if err != nil {
		panic(err.Error())
	}
	paramMap["cors"] = corsData["cors"].(string)
	hostData, err := doc.Table("global.codefresh")
	if err != nil {
		panic(err.Error())
	}
	paramMap["CF_HOST"] = hostData["url"].(string)
}