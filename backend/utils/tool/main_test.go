package tool

import (
	"path/filepath"
	"testing"

	"ktt/backend/utils/log"
	strutil "ktt/backend/utils/string"
)

func TestMain(m *testing.M) {
	log.Init(filepath.Join(strutil.RootPath(), "logs"))
}
