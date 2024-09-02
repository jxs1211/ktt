package watch

import (
	"context"
	"errors"
	"strings"
	"sync"

	"ktt/backend/types"
	k8sutil "ktt/backend/utils/k8s"

	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type WatcherManager struct {
	ctx        context.Context
	watcherMap map[string]*Watcher
}

func NewWatcherManager() *WatcherManager {
	return &WatcherManager{
		watcherMap: make(map[string]*Watcher),
	}
}

func (wm *WatcherManager) Start(ctx context.Context) {
	wm.ctx = ctx
}

func (wm *WatcherManager) validate(watcher *Watcher) error {
	if len(watcher.Name) == 0 {
		return errors.New("watcher name must not be empty")
	}
	if len(watcher.ResourceName) == 0 {
		return errors.New("watcher resourceName must not be empty")
	}
	// check the resourceNS is valid, the resourceNS can't be empty unless the resourceType is not namespaced
	// if len(watcher.ResourceNS) == 0 {
	// 	if !strings.Contains(NonNamespaecedResourceTypes, watcher.ResourceType) {
	// 		return errors.New("watcher resourceNS must not be empty")
	// 	}
	// }
	if len(watcher.ResourceType) <= 0 {
		return errors.New("watcher resourceType must not be empty")
	}
	if len(watcher.WatchingTypes) == 0 {
		return errors.New("watcher watchingTypes must not be empty")
	}
	for _, t := range watcher.WatchingTypes {
		if t == "" {
			return errors.New("watcher watchingTypes must not be empty")
		}
		if !k8sutil.Contains(watchingTypesString(), string(t)) {
			return errors.New("watcher watchingTypes must be in " + strings.Join(watchingTypesString(), ", "))
		}
	}
	if _, ok := wm.watcherMap[watcher.Name]; ok {
		return errors.New("watcher already exists")
	}
	return nil
}

func (wm *WatcherManager) addWatcher(watcher *Watcher) error {
	err := wm.validate(watcher)
	if err != nil {
		return err
	}
	config, err := k8sutil.NewConfig()
	if err != nil {
		return err
	}
	client, err := k8sutil.NewClient(config)
	if err != nil {
		return err
	}
	logr, err := zap.NewProduction()
	if err != nil {
		return err
	}
	var mu sync.Mutex
	watcher.kclient = client
	watcher.result = make(chan Event)
	watcher.logger = logr
	watcher.Mutex = &mu
	// add error watcher if
	wm.watcherMap[watcher.Name] = watcher
	return nil
}

func (wm *WatcherManager) AddWatcher(
	name, resourceType, resourceNS,
	resourceName string, watchingTypes []string) types.JSResp {
	wTypes := make([]watchingType, 0, len(watchingTypes))
	for _, t := range watchingTypes {
		wTypes = append(wTypes, watchingType(t))
	}
	caser := cases.Title(language.English)
	watcher := &Watcher{
		Name:          name,
		ResourceType:  caser.String(resourceType),
		ResourceNS:    resourceNS,
		ResourceName:  resourceName,
		WatchingTypes: wTypes,
	}
	err := wm.addWatcher(watcher)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
	}
}

func (wm *WatcherManager) updateWatcher(watcher *Watcher) error {
	err := wm.validate(watcher)
	if err != nil {
		return err
	}
	w, err := wm.GetWatcher(watcher.Name)
	if err != nil {
		return err
	}
	w.Name = watcher.Name
	w.WatchingTypes = watcher.WatchingTypes
	wm.watcherMap[watcher.Name] = w
	return nil
}

// func (wm *WatcherManager) UpdateWatcher(watcher Watcher) types.JSResp {
// 	err := wm.updateWatcher(watcher)
// 	if err != nil {
// 		return types.FailedResp(err.Error())
// 	}
// 	return types.JSResp{
// 		Success: true,
// 	}
// }

func (wm *WatcherManager) RemoveWatcher(watcher *Watcher) types.JSResp {
	delete(wm.watcherMap, watcher.Name)
	return types.JSResp{
		Success: true,
	}
}

func (wm *WatcherManager) GetWatcher(name string) (*Watcher, error) {
	v, ok := wm.watcherMap[name]
	if !ok {
		return nil, errors.New("watcher not found")
	}
	return v, nil
}

func (wm *WatcherManager) GetWatcherList() []string {
	var res []string
	for k := range wm.watcherMap {
		res = append(res, k)
	}
	return res
}

func (wm *WatcherManager) RunWatcher(name string) types.JSResp {
	err := wm.runWatcher(name)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
	}
}

func (wm *WatcherManager) runWatcher(name string) error {
	v, ok := wm.watcherMap[name]
	if !ok {
		return errors.New("watcher not found")
	}
	go v.Run(1)
	return nil
}

func (wm *WatcherManager) StopWatcher(name string) types.JSResp {
	err := wm.stopWatcher(name)
	if err != nil {
		return types.FailedResp(err.Error())
	}
	return types.JSResp{
		Success: true,
	}
}

func (wm *WatcherManager) stopWatcher(name string) error {
	v, ok := wm.watcherMap[name]
	if !ok {
		return errors.New("watcher not found")
	}
	v.Stop()
	return nil
}
