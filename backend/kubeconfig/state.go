package kubeconfig

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func HomeDir() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE") // windows
	}
	return home
}

func kubectxPrevCtxFile() (string, error) {
	home := HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "kubectx"), nil
}

// writeLastContext saves the specified value to the state file.
// It creates missing parent directories.
func writeLastContext(path, value string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrap(err, "failed to create parent directories")
	}
	return os.WriteFile(path, []byte(value), 0644)
}
