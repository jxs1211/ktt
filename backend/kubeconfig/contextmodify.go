package kubeconfig

import (
	"gopkg.in/yaml.v3"
)

func (k *Kubeconfig) ModifyCurrentContext(name string) error {
	currentCtxNode := valueOf(k.rootNode, "current-context")
	if currentCtxNode != nil {
		currentCtxNode.Value = name
		return nil
	}

	// if current-context field doesn't exist, create new field
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "current-context",
		Tag:   "!!str"}
	valueNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: name,
		Tag:   "!!str"}
	k.rootNode.Content = append(k.rootNode.Content, keyNode, valueNode)
	return nil
}
