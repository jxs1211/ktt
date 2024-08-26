package log

import (
	uberzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/go-logr/logr"
)

func New(options ...uberzap.Option) logr.Logger {
	opts := zap.Options{
		Development:     true,
		StacktraceLevel: zapcore.DPanicLevel,
		ZapOpts:         options,
	}
	return zap.New(zap.UseFlagOptions(&opts))
}
