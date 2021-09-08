package handler

import (
	"encoding/json"
	xwcv1 "github.com/ligao-cloud-native/kubemc/pkg/apis/xwc/v1"
	"io/ioutil"
	"net/http"
)


type ClusterScaleType string

const (
	ClusterScaleTypeUp ClusterScaleType = "up"
	ClusterScaleTypeDown ClusterScaleType = "Down"
)


type ClusterScale struct {
	Scale ClusterScaleType `json:"scale"`
	Workers []xwcv1.Node `json:"workers"`
}

// TODO: check k8s cluster in xwc
func ValidateCluster(wc *xwcv1.WorkloadCluster) error {
	return nil
}


func AddOrRemoveWorkerNode(xwcNode, toScaledNode []xwcv1.Node) {}


func ParseRequestBody(req *http.Request, v interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	if err := json.Unmarshal(buf, v); err != nil {
		return err
	}

	return nil
}

