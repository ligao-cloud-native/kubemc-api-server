package watcher

import (
	"context"
	"github.com/ligao-cloud-native/kubemc-api-server/crd/manager"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"time"
)

var XWCWatcher *XwcWatcher

type XwcWatcher struct {
	kubeClient    *kubernetes.Clientset
	dynamicClient dynamic.Interface
	xwcClient     interface{}
	Cache         cache.Store
	Manager       *manager.XwcManager
}

func NewXwcWatcher() *XwcWatcher {
	//config, err := buildConfig("")
	//if err != nil {
	//	klog.Errorf("Failed to build config, err: %v", err)
	//	os.Exit(1)
	//}
	//
	//kubeClient := kubernetes.NewForConfigOrDie(config)
	//dynamicClient := dynamic.NewForConfigOrDie(config)
	//
	//XWCWatcher = &XwcWatcher{
	//	kubeClient:    kubeClient,
	//	dynamicClient: dynamicClient,
	//}

	return XWCWatcher

}

func (w *XwcWatcher) Start() {
	// 添加所有的xwc实例
	//wcs, err := w.kubeClient.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	klog.Warningf("init workloadcluster resource error, %v", err)
	//	return
	//}

	//for _, wc := range wcs.Items {
	//	// if status is success
	//	w.Manager.Add(&wc)
	//
	//}

	// list and watch crd
	w.watch()

	select {}

}

func (w *XwcWatcher) watch() {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (result runtime.Object, err error) {
				return w.kubeClient.AppsV1().Deployments("").List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return w.kubeClient.AppsV1().Deployments("").Watch(context.TODO(), options)
			},
		},
		&appv1.Deployment{},
		1*time.Minute,
		WorkloadClusterEventHandler{})

	go controller.Run(wait.NeverStop)

	w.Cache = store
}

func buildConfig(kubeconfig string) (config *rest.Config, err error) {
	if kubeconfig == "" {
		klog.Info("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		klog.Infof("use configuration from %s", kubeconfig)
		//config, err = clientcmd.BuildConfigFromFlags("https://1.2.3.4:6443", "/root/.kube/config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, err
	}

	config.QPS = float32(100)
	config.Burst = int(200)
	config.ContentType = "application/json"

	return config, nil
}
