package adapters

import (
	"fmt"
	"github.com/german-muzquiz/factory-crd/pkg/domain"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

const FactoryGVR = "factories.v1.poc.techm.com"

type KubeFactoryRepository struct {
	factories map[string]domain.Factory
	stopCh    chan struct{}
}

func (r *KubeFactoryRepository) GetFactories() map[string]domain.Factory {
	return r.factories
}

func (r *KubeFactoryRepository) Init() {
	r.factories = map[string]domain.Factory{}
	informer := r.createInformer()

	r.stopCh = make(chan struct{})

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.onAddFactory,
		UpdateFunc: r.onUpdateFactory,
		DeleteFunc: r.onDeleteFactory,
	})

	log.Info("Connecting to cluster API")
	go informer.Run(make(chan struct{}))
}

func (r *KubeFactoryRepository) Shutdown() {
	r.stopCh <- struct{}{}
}

func (r *KubeFactoryRepository) createInformer() cache.SharedIndexInformer {
	config := r.createConfig()
	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, corev1.NamespaceAll, nil)
	gvr, _ := schema.ParseResourceArg(FactoryGVR)
	return factory.ForResource(*gvr).Informer()
}

func (r *KubeFactoryRepository) createConfig() *rest.Config {
	var config *rest.Config
	var err error

	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		// running outside kubernetes
		h := os.Getenv("HOME")
		config, err = clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", h))

	} else {
		// running inside kubernetes
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		panic(err.Error())
	}

	return config
}

func (r *KubeFactoryRepository) onAddFactory(obj interface{}) {
	log.Info("Factory added")
	u := obj.(*unstructured.Unstructured)
	f := r.factoryFromUnstructured(u)
	r.factories[u.GetName()] = f
}

func (r *KubeFactoryRepository) onUpdateFactory(_, newObj interface{}) {
	log.Info("Factory updated")
	u := newObj.(*unstructured.Unstructured)
	f := r.factoryFromUnstructured(u)
	r.factories[u.GetName()] = f
}

func (r *KubeFactoryRepository) onDeleteFactory(obj interface{}) {
	log.Info("Factory deleted")
	u := obj.(*unstructured.Unstructured)
	delete(r.factories, u.GetName())
}

func (r *KubeFactoryRepository) factoryFromUnstructured(obj *unstructured.Unstructured) domain.Factory {
	f := domain.Factory{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &f)
	if err != nil {
		log.WithError(err).Error("Error marshaling factory from unstructured to typed")
	}
	return f
}
