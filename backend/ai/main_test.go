package ai

import (
	"path/filepath"
	"testing"

	"ktt/backend/utils/log"
	strutil "ktt/backend/utils/string"
)

var (
	wrongNginxErrMsg = `
	{
		"kind": "Pod",
		"name": "default/nginx-deployment-7c79c4bf97-kspnz",
		"error": [
		{
			"Text": "Back-off pulling image \"nginx2:latest\"",
			"KubernetesDoc": "",
			"Sensitive": []
		}
		],
		"parentObject": "Deployment/nginx-deployment"
	}
	`
)

func TestMain(m *testing.M) {
	log.Init(filepath.Join(strutil.RootPath(), "logs"))
}
