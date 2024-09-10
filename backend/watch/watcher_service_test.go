package watch

import (
	"ktt/backend/types"
	"reflect"
	"testing"
	"time"
)

func TestWatcherManager_addWatcher(t *testing.T) {
	tests := []struct {
		name       string
		watcherMap map[string]*Watcher
		watcher    *Watcher
		wantErr    bool
	}{
		{
			name:       "test",
			watcherMap: map[string]*Watcher{},
			watcher: &Watcher{
				Name:         "test",
				ResourceType: "pod",
				WatchingTypes: []watchingType{
					watchingTypeCreation,
					watchingTypeDeletion,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &WatcherManager{
				watcherMap: tt.watcherMap,
			}
			err := wm.addWatcher(tt.watcher)
			if tt.wantErr {
				if err == nil {
					t.Errorf("WatcherManager.addWatcher() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("WatcherManager.addWatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestWatcherManager_AddWatcher(t *testing.T) {
// 	type fields struct {
// 		watcherMap map[string]*Watcher
// 	}
// 	type args struct {
// 		watcher Watcher
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   types.JSResp
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			wm := &WatcherManager{
// 				watcherMap: tt.fields.watcherMap,
// 			}
// 			if got := wm.AddWatcher(tt.args.watcher); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("WatcherManager.AddWatcher() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestWatcherManager_RunWatcher(t *testing.T) {
	tests := []struct {
		name       string
		watcherMap map[string]*Watcher
		want       types.JSResp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &WatcherManager{
				watcherMap: tt.watcherMap,
			}
			if got := wm.RunWatcher(tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WatcherManager.RunWatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWatcherManager_runWatcher(t *testing.T) {
	tests := []struct {
		name    string
		watcher *Watcher
		wantErr bool
	}{
		{
			name: "test",
			watcher: &Watcher{
				Name:         "test",
				ResourceType: "pod",
				ResourceNS:   "default",
				ResourceName: "test",
				WatchingTypes: []watchingType{
					watchingTypeCreation,
					watchingTypeDeletion,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := NewWatcherManager()
			if err := wm.addWatcher(tt.watcher); err != nil {
				t.Errorf("WatcherManager.addWatcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := wm.runWatcher(tt.name); (err != nil) != tt.wantErr {
				t.Errorf("WatcherManager.runWatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
			// create resource to be watched
			// wm.watcherMap[tt.name].kclient.CoreV1().Pods("default").Create(context.TODO(), &v1.Pod{
			// 	ObjectMeta: metav1.ObjectMeta{
			// 		Name: "test",
			// 	},
			// }, metav1.CreateOptions{})
			time.Sleep(4 * time.Second)
			// stop
			wm.stopWatcher(tt.name)
			if err := wm.runWatcher(tt.name); (err != nil) != tt.wantErr {
				t.Errorf("WatcherManager.runWatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(4 * time.Second)
			// stop
			wm.stopWatcher(tt.name)
		})
	}
}

func TestWatcherManager_validate(t *testing.T) {
	tests := []struct {
		name       string
		watcherMap map[string]*Watcher
		watcher    *Watcher
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &WatcherManager{
				watcherMap: tt.watcherMap,
			}
			if err := wm.validate(tt.watcher); (err != nil) != tt.wantErr {
				t.Errorf("WatcherManager.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
