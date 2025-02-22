package logger

import (
	"context"
	"log"
	"os"
	"wallet-api/internal/pkg/env"

	"github.com/sirupsen/logrus"
)

func Init() *os.File {
	logMode := env.GetWithDefault("LOG_MODE", "stdout")

	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)

	if logMode == "file" {
		logPath := env.GetWithDefault("LOG_PATH", "./application.log")

		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}

		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		})
		logrus.SetOutput(logFile)

		return logFile
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	return nil
}

// Info is used to logging level info.
func Info(ctx context.Context, metadata map[string]any, message string) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Info(message)
}

// Infof is used to logging level info with formatted args.
func Infof(ctx context.Context, metadata map[string]any, message string, args ...any) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Infof(message, args...)
}

// Warning is used to logging level warning.
func Warning(ctx context.Context, metadata map[string]any, message string) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Warning(message)
}

// Warningf is used to logging level warning with formatted args.
func Warningf(ctx context.Context, metadata map[string]any, message string, args ...any) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Warningf(message, args...)
}

// Error is used to logging level error.
func Error(ctx context.Context, metadata map[string]any, message string) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Error(message)
}

// Errorf is used to logging level error with formatted args.
func Errorf(ctx context.Context, metadata map[string]any, message string, args ...any) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Errorf(message, args...)
}

// Fatal is used to logging level fatal and program will exit.
func Fatal(ctx context.Context, metadata map[string]any, message string) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Fatal(message)
}

// Fatalf is used to logging level fatal with formatted args.
func Fatalf(ctx context.Context, metadata map[string]any, message string, args ...any) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Fatalf(message, args...)
}

// Trace is used to logging level trace.
func Trace(ctx context.Context, metadata map[string]any, message string) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Trace(message)
}

// Tracef is used to logging level trace with formatted args.
func Tracef(ctx context.Context, metadata map[string]any, message string, args ...any) {
	logrus.WithFields(logrus.Fields(wrapper(ctx, metadata))).Tracef(message, args...)
}

func wrapper(ctx context.Context, metadata map[string]any) map[string]any {
	if v, ok := ctx.Value("requestid").(string); ok {
		metadata["trace_id"] = v
	}

	return metadata
}
