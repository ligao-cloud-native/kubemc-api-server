package manager

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"testing"
)

// DownstreamController watch kubernetes api server and send change to edge
type DownstreamController struct {
	podManager Manager
}

// 处理pod事件
func (dc *DownstreamController) syncPod() {
	for {
		select {
		case e := <-dc.podManager.Events():
			pod, ok := e.Object.(*v1.Pod)
			if !ok {
				klog.Warningf("object type: %T unsupported", pod)
				continue
			}

			// new a message and send
		}
	}
}

func TestPodManager(t *testing.T) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("Master",
		"KubeConfig")
	if err != nil {
		klog.Warningf("get kube config failed with error: %s", err)
	}
	kubeConfig.QPS = float32(1)
	kubeConfig.Burst = int(1)

	cli, _ := kubernetes.NewForConfig(kubeConfig)

	podManager, err := NewPodManager(cli, v1.NamespaceAll, "nodeName")
	if err != nil {
		klog.Warningf("create pod manager failed with error: %s", err)
	}

	dc := DownstreamController{podManager: podManager}

	go dc.syncPod()
}
