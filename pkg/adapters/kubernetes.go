package adapters

import (
	"context"
	"fmt"
	"github.com/german-muzquiz/factory-crd/pkg/domain"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

const (
	FactoryGVR = "factories.v1.poc.techm.com"
)

type KubeFactoryRepository struct {
	factories  map[string]*unstructured.Unstructured
	stopCh     chan struct{}
	kubeClient dynamic.Interface
}

func (r *KubeFactoryRepository) GetFactories() map[string]domain.Factory {
	result := map[string]domain.Factory{}
	for n, f := range r.factories {
		result[n] = r.factoryFromUnstructured(f)
	}
	return result
}

func (r *KubeFactoryRepository) UpdateCapacity(name string, newCapacity int) {
	f := r.factories[name]
	capString := fmt.Sprintf("%d vehicles per minute", newCapacity)
	err := unstructured.SetNestedField(f.Object, capString, "status", "currentCapacity")
	if err != nil {
		log.WithError(err).Errorf("Error setting new capacity to unstructred object %s", name)
		return
	}

	gvr, _ := schema.ParseResourceArg(FactoryGVR)
	ns := f.GetNamespace()
	if ns == "" {
		ns = "default"
	}

	updated, err := r.kubeClient.Resource(*gvr).Namespace(ns).Update(
		context.Background(), f, metav1.UpdateOptions{})
	if err != nil {
		log.WithError(err).Errorf("Error updating factory %s", name)
		return
	}
	r.factories[updated.GetName()] = updated
}

func (r *KubeFactoryRepository) Init() {
	r.factories = map[string]*unstructured.Unstructured{}
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
	r.kubeClient = dc
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
	u := obj.(*unstructured.Unstructured)
	log.Infof("Factory %s added", u.GetName())
	r.factories[u.GetName()] = u
}

func (r *KubeFactoryRepository) onUpdateFactory(_, newObj interface{}) {
	u := newObj.(*unstructured.Unstructured)
	log.Infof("Factory %s updated", u.GetName())
	r.factories[u.GetName()] = u
}

func (r *KubeFactoryRepository) onDeleteFactory(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	log.Infof("Factory %s deleted", u.GetName())
	delete(r.factories, u.GetName())
}

func (r *KubeFactoryRepository) factoryFromUnstructured(obj *unstructured.Unstructured) domain.Factory {
	f := &domain.Factory{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), f)
	if err != nil {
		log.WithError(err).Error("Error marshaling factory from unstructured to typed")
	}
	return *f
}
