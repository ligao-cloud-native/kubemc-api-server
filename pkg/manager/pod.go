package manager

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"reflect"
	"sync"
)

// CachePod is the struct save pod data for check pod is really changed
type CachePod struct {
	metav1.ObjectMeta
	Spec v1.PodSpec
}

// PodManager is a manager watch pod change event
type PodManager struct {
	// events from watch kubernetes api server
	realEvents chan watch.Event

	// events merged
	mergedEvents chan watch.Event

	// pods, key is UID, value is *v1.Pod
	pods sync.Map
}

// Events return a channel, can receive all pod event
func (pm *PodManager) Events() chan watch.Event {
	return pm.mergedEvents
}

// NewPodManager create PodManager from config
func NewPodManager(kubeClient *kubernetes.Clientset, namespace, nodeName string) (*PodManager, error) {

	// 定义一个pod对象的listwatch
	var lw *cache.ListWatch

	if "" == nodeName {
		lw = cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "pods", namespace, fields.Everything())
	} else {
		selector := fields.OneTermEqualSelector("spec.nodeName", nodeName)
		lw = cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "pods", namespace, selector)
	}

	//定义一个事件接收通道
	realEvents := make(chan watch.Event, 1)
	mergedEvents := make(chan watch.Event, 1)

	// 定义一个事件处理Handler
	rh := NewCommonResourceEventHandler(realEvents)

	// listwatch pod且处理事件
	si := cache.NewSharedInformer(lw, &v1.Pod{}, 0)
	si.AddEventHandler(rh) // rh需要实现ResourceEventHandler接口

	pm := &PodManager{realEvents: realEvents, mergedEvents: mergedEvents}

	// starts and runs the shared informer
	stopNever := make(chan struct{})
	go si.Run(stopNever)
	go pm.merge()

	return pm, nil
}

// 对接收pod的真实事件进行合处理
func (pm *PodManager) merge() {
	for re := range pm.realEvents {
		pod := re.Object.(*v1.Pod)
		switch re.Type {
		case watch.Added:
			pm.pods.Store(pod.UID, &CachePod{ObjectMeta: pod.ObjectMeta, Spec: pod.Spec})
			if pod.DeletionTimestamp == nil {
				pm.mergedEvents <- re
			} else {
				re.Type = watch.Modified
				pm.mergedEvents <- re
			}
		case watch.Deleted:
			pm.pods.Delete(pod.UID)
			pm.mergedEvents <- re
		case watch.Modified:
			value, ok := pm.pods.Load(pod.UID)
			pm.pods.Store(pod.UID, &CachePod{ObjectMeta: pod.ObjectMeta, Spec: pod.Spec})
			if ok {
				cachedPod := value.(*CachePod)
				if pm.isPodUpdated(cachedPod, pod) {
					pm.mergedEvents <- re
				}
			} else {
				pm.mergedEvents <- re
			}
		default:
			klog.Warningf("event type: %s unsupported", re.Type)
		}

	}
}

func (pm *PodManager) isPodUpdated(old *CachePod, new *v1.Pod) bool {
	// does not care fields
	old.ObjectMeta.ResourceVersion = new.ObjectMeta.ResourceVersion
	old.ObjectMeta.Generation = new.ObjectMeta.Generation

	// return true if ObjectMeta or Spec changed, else false
	return !reflect.DeepEqual(old.ObjectMeta, new.ObjectMeta) || !reflect.DeepEqual(old.Spec, new.Spec)
}
