package watch

import (
	kwatch "k8s.io/apimachinery/pkg/watch"
)

type Event struct {
	Type   kwatch.EventType
	Result interface{}
}
