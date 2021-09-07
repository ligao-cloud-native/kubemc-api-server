package watcher

type WorkloadClusterEventHandler struct{}

func (h WorkloadClusterEventHandler) OnAdd(obj interface{}) {}

func (h WorkloadClusterEventHandler) OnUpdate(oldObj, newObj interface{}) {}

func (h WorkloadClusterEventHandler) OnDelete(obj interface{}) {}
