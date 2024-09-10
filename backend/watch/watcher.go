package watch

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/k8sgpt-ai/k8sgpt/pkg/analyzer"
	"github.com/k8sgpt-ai/k8sgpt/pkg/common"
	k8sgpt "github.com/k8sgpt-ai/k8sgpt/pkg/kubernetes"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kwatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const resyncPeriod = 10 * time.Second

var (
	caser = cases.Title(language.English)
)

type watchingType string

const (
	watchingTypeCreation watchingType = "creation"
	watchingTypeDeletion watchingType = "deletion"
	watchingTypeUpdate   watchingType = "update"
	watchingTypeError    watchingType = "error"
)

var watchingTypes = []watchingType{
	watchingTypeCreation,
	watchingTypeDeletion,
	watchingTypeUpdate,
	watchingTypeError,
}

func typesContains(types []watchingType, targetType watchingType) bool {
	for _, t := range types {
		if t == targetType {
			return true
		}
	}
	return false
}

func watchingTypesString() []string {
	types := make([]string, 0, len(watchingTypes))
	for _, t := range watchingTypes {
		types = append(types, string(t))
	}
	return types
}

type Watcher struct {
	Name                                   string
	kclient                                kubernetes.Interface
	ResourceType, ResourceNS, ResourceName string
	WatchingTypes                          []watchingType
	result                                 chan Event
	Ctx                                    context.Context
	Cancel                                 context.CancelFunc
	// lister that can list resource from a shared cache
	// returns true when the resource cache is ready
	ListerSynced    cache.InformerSynced
	informerFactory informers.SharedInformerFactory
	informer        cache.SharedIndexInformer
	queue           workqueue.RateLimitingInterface
	logger          *zap.Logger
	*sync.Mutex
	started bool
	stopped bool
}

func NewWatcher(
	name string,
	client *kubernetes.Clientset,
	resourceType, resourceNS, resourceName string,
	// ctx context.Context,
	// ctxCancel context.CancelFunc,
	watchingTypes []string,
	logr *zap.Logger,
) *Watcher {
	wTypes := make([]watchingType, 0, len(watchingTypes))
	for _, t := range watchingTypes {
		wTypes = append(wTypes, watchingType(t))
	}
	// ctx, cancel := context.WithCancel(context.Background())
	resType := caser.String(resourceType)
	var mu sync.Mutex
	w := &Watcher{
		Name:          name,
		kclient:       client,
		ResourceType:  resType,
		ResourceNS:    resourceNS,
		ResourceName:  resourceName,
		WatchingTypes: wTypes,
		result:        make(chan Event),
		logger:        logr,
		Mutex:         &mu,
	}
	return w
}

func (c *Watcher) withInformer(name string) error {
	c.informerFactory = informers.NewSharedInformerFactory(c.kclient, 1*time.Second)
	switch caser.String(name) {
	case "Node":
		c.informer = c.informerFactory.Core().V1().Nodes().Informer()
	case "Service":
		c.informer = c.informerFactory.Core().V1().Services().Informer()
	case "Pod":
		c.informer = c.informerFactory.Core().V1().Pods().Informer()
	default:
		return errors.New("informer type not found")
	}
	c.ListerSynced = c.informer.HasSynced
	c.informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.handleAdd,
			DeleteFunc: c.handleDel,
		},
	)
	return nil
}

func (c *Watcher) handleAdd(obj interface{}) {
	// early return if the resource is not matched
	if !c.isMatched(obj) {
		return
	}
	c.queue.Add(obj)
}

func (c *Watcher) handleDel(obj interface{}) {
	if !c.isMatched(obj) {
		return
	}
	c.queue.Add(obj)
}

func (c *Watcher) isMatched(item interface{}) bool {
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("getting key from cahce %s\n", err.Error())
		return false
	}

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("splitting key into namespace and name %s\n", err.Error())
		return false
	}
	if ns != c.ResourceNS {
		return false
	}
	if name != c.ResourceName {
		return false
	}
	return true
}

func (c *Watcher) Run(workers int) {
	c.Lock()
	if c.started {
		c.logger.Sugar().Errorf("watcher %s already started", c.Name)
		return
	}
	c.queue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), c.Name)
	c.Ctx, c.Cancel = context.WithCancel(context.Background())
	err := c.withInformer(c.ResourceType)
	if err != nil {
		c.logger.Sugar().Errorf("error creating informer: %s", err.Error())
		return
	}
	c.informerFactory.Start(c.Ctx.Done())
	c.logger.Info("Starting watcher")
	defer c.logger.Info("Shutting down watcher")
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()
	if !cache.WaitForNamedCacheSync(c.Name, c.Ctx.Done(), c.ListerSynced) {
		return
	}
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(c.Ctx, c.worker, resyncPeriod)
	}
	// go run error wacther
	c.started = true
	c.Unlock()

	if typesContains(c.WatchingTypes, watchingTypeError) {
		go c.runWatchingError()
	}

	for {
		select {
		case <-c.Ctx.Done():
			c.logger.Sugar().Infof("stopped watcher: %s", c.Name)
			return
		case result, ok := <-c.result:
			if !ok {
				continue
			}
			c.logger.Sugar().Infof("received event: %s", result.Type)
			// wailsruntime.Events.EmitEvent(ctx, &v1.Event{}
		}
	}
}

func (c *Watcher) Stop() {
	defer c.Unlock()
	c.Lock()
	// tryAcquire result before stopping the watcher
	select {
	case <-c.Ctx.Done():
	default:
		c.Cancel()
		c.queue.ShutDown()
		c.stopped = true
		c.started = false
	}
}

func (c *Watcher) worker(ctx context.Context) {
	defer c.logger.Info("stop worker")
	c.logger.Sugar().Infof("started watcher to watch: %s, %s, %s", c.ResourceType, c.ResourceNS, c.ResourceName)
	for c.processItem(ctx) {
	}
}

func (c *Watcher) processItem(ctx context.Context) bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Forget(item)
	key, err := cache.MetaNamespaceKeyFunc(item)
	if err != nil {
		fmt.Printf("getting key from cahce %s\n", err.Error())
	}
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		fmt.Printf("splitting key into namespace and name %s\n", err.Error())
		return true
	}
	err = c.syncResource(ctx, ns, name)
	if err != nil {
		c.logger.Sugar().Errorf("Failed to synchronize nodes: %s", err)
		return true
	}
	return true
}

func (c *Watcher) listerBy(ctx context.Context, typ, ns, name string) (runtime.Object, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		switch typ {
		case "Node":
			return c.informerFactory.Core().V1().Nodes().Lister().Get(name)
		case "Service":
			return c.informerFactory.Core().V1().Services().Lister().Services(ns).Get(name)
		case "Pod":
			return c.informerFactory.Core().V1().Pods().Lister().Pods(ns).Get(name)
		default:
			return nil, fmt.Errorf("resource type [%s/%s] not found", ns, name)
		}
	}
}

func (c *Watcher) syncResource(ctx context.Context, ns, name string) error {
	obj, err := c.listerBy(ctx, c.ResourceType, ns, name)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			c.logger.Sugar().Infof("get svc failed: ", err)
			return err
		}
		// deletion
		c.logger.Info("resource deleted")
		c.result <- Event{
			Type:   kwatch.Deleted,
			Result: obj,
		}
		return nil
	}
	// creation
	c.logger.Info("resource created")
	c.result <- Event{
		Type:   kwatch.Added,
		Result: obj,
	}
	return nil
}

// func (w *Watcher) resultChan() <-chan kwatch.Event {
// 	return w.result
// }

func (w *Watcher) runWatchingError() error {
	// call the run function using time.Ticker with select
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	coreAnalyzerMap, _ := analyzer.GetAnalyzerMap()
	analyzer, ok := coreAnalyzerMap[w.ResourceType]
	if !ok {
		w.logger.Error("analyzer not found")
		return errors.New("analyzer not found")
	}
	// client, err := k8sgpt.NewClient()
	analyzerConfig := common.Analyzer{
		Client:        &k8sgpt.Client{Client: w.kclient},
		Context:       w.Ctx,
		Namespace:     w.ResourceNS,
		LabelSelector: "",
		AIClient:      nil,
		OpenapiSchema: nil,
	}
	for {
		select {
		case <-w.Ctx.Done():
			w.logger.Info("stopped watching error")
			return w.Ctx.Err()
		case <-ticker.C:
			w.logger.Info("watching error")
			results, err := analyzer.Analyze(analyzerConfig)
			if err != nil {
				return err
			}
			w.logger.Sugar().Infof("watching err results: %+v", results)
			if len(results) > 0 {
				w.result <- Event{
					Type:   kwatch.Error,
					Result: results,
				}
				w.logger.Sugar().Infof("runWatchingError: %+v", results)
				// emit event
				// wailsruntime.Events.EmitEvent(ctx, &v1.Event{}
			}
		}
	}
}
