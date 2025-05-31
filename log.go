package puppy

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func isConsole(logFile string) bool {
	if logFile == "" || logFile == "console" {
		return true
	}
	return false
}

func RootLog(id string, logFile string, level slog.Level) *slog.Logger {
	var out io.Writer = os.Stdout
	if !isConsole(logFile) {
		out = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10, // megabytes
			MaxBackups: 10,
			MaxAge:     10, //days
		}
	}
	h := slog.NewJSONHandler(out, &slog.HandlerOptions{Level: level}).WithAttrs([]slog.Attr{slog.String("id", id)})
	logger := slog.New(h)
	slog.SetDefault(logger)
	return logger
}

/*
func BuildRootLog(conf AppConf) *slog.Logger {
	logFile := conf.LoggerConf.LogFile
	if logFile == "" {
		logFile = "dcp.log"
	}

	var level slog.Level

	l := strings.ToUpper(conf.LoggerConf.Level)
	switch l {
	case "DEBUG", "debug":
		level = slog.LevelDebug
	case "WARN", "warn":
		level = slog.LevelWarn
	case "ERROR", "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return RootLog(conf.ID, logFile, level)
}
*/
