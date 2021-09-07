package main

import (
	"github.com/ligao-cloud-native/kubemc-api-server/api"
	"k8s.io/component-base/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	s := api.NewAPIServer()
	s.Run()

}
