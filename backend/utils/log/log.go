package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.SugaredLogger

func Init(logDir string) {
	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	logPath := filepath.Join(logDir, "app.log")
	// Check if the log file exists
	_, err := os.Stat(logPath)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		// File doesn't exist, create it
		file, err := os.Create(logPath)
		if err != nil {
			panic(err)
		}
		file.Close()
	} else {
		// File exists, append a separator
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		separator := fmt.Sprintf("\n\n---------- restarted at %s ----------\n\n", time.Now().Format(time.RFC3339))
		if _, err := file.WriteString(separator); err != nil {
			panic(err)
		}
	}

	// Configure log rotation
	logFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Create custom encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	encoderConfig.StacktraceKey = "stacktrace"

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		zap.InfoLevel,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.FatalLevel))
	Log = logger.Sugar()
}

func Info(msg string, keysAndValues ...interface{}) {
	Log.WithOptions(zap.AddCallerSkip(1)).Infow(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	Log.WithOptions(zap.AddCallerSkip(1)).Errorw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	Log.WithOptions(zap.AddCallerSkip(1)).Fatalw(msg, keysAndValues...)
}
