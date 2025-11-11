package logger

import (
	"io"
	defLog "log"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	logsDir = "logs"

	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var (
	MsgCantCreateLogsDir = "can not create logs directory"
	MsgCantOpenLogsFile  = "can not open logs file"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

func New(env string) *Logger {
	mustPrepareLogsDir()

	logFile := mustOpenLogsFile()

	var handler slog.Handler
	switch env {
	case envLocal:
		handler = slog.NewTextHandler(
			io.MultiWriter(os.Stdout, logFile),
			&slog.HandlerOptions{Level: slog.LevelDebug},
		)
	case envDev:
		handler = slog.NewJSONHandler(
			io.MultiWriter(os.Stdout, logFile),
			&slog.HandlerOptions{Level: slog.LevelDebug},
		)
	case envProd:
		handler = slog.NewJSONHandler(
			logFile,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)
	default:
		handler = slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)
	}

	return &Logger{
		Logger: slog.New(handler),
		file:   logFile,
	}
}

func (l *Logger) MustClose() {
	_ = l.file.Close()
}

func mustPrepareLogsDir() {
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		defLog.Fatal(MsgCantCreateLogsDir)
	}
}

func mustOpenLogsFile() *os.File {
	logFile, err := os.OpenFile(
		logsFilepath(),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		defLog.Fatal(MsgCantOpenLogsFile)
	}

	return logFile
}

func logsFilepath() string {
	logFileName := time.Now().Format("15-04-05_02-01-2006") + ".log"

	return filepath.Join(logsDir, logFileName)
}
