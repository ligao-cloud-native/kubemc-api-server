package watcher

//var KubeClientSet KubeClient
//var once sync.Once
//
//
//type KubeClient struct {
//	kubeClient *kubernetes.Clientset
//	dynamicClient dynamic.Interface
//}
//
//
//func InitKubeClient() {
//	once.Do(func() {
//		config, err := buildConfig("")
//		if err != nil {
//			klog.Errorf("Failed to build config, err: %v", err)
//			os.Exit(1)
//		}
//
//		kubeClient := kubernetes.NewForConfigOrDie(config)
//		dynamicClient := dynamic.NewForConfigOrDie(config)
//
//		KubeClientSet = KubeClient{
//			kubeClient: kubeClient,
//			dynamicClient: dynamicClient,
//		}
//	})
//
//}
