package kubeconfig

// GetCurrentContext returns "current-context" value in given
// kubeconfig object Node, or returns "" if not found.
func (k *Kubeconfig) GetCurrentContext() string {
	v := valueOf(k.rootNode, "current-context")
	if v == nil {
		return ""
	}
	return v.Value
}
