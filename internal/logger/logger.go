package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hvarillas/smbsync/internal/notification"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Sugar *zap.SugaredLogger

func Init(logPath, logLevel string) {
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %v", err))
	}

	level := zap.InfoLevel
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
	
	useColors := true
	if runtime.GOOS == "windows" {
		useColors = enableVirtualTerminal()
	}
	
	if useColors {
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	
	consoleEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	consoleEncoderConfig.EncodeCaller = nil

	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(fileEncoderConfig),
		zapcore.AddSync(logFile),
		level,
	)

	core := zapcore.NewTee(consoleCore, fileCore)

	telegramHook := func(entry zapcore.Entry) error {
		if entry.Level >= zapcore.ErrorLevel {
			hostname, err := os.Hostname()
			if err != nil {
				hostname = "unknown"
			}
			message := fmt.Sprintf("<b>🚨 Alerta de Error [SMBSync] (%s) 🚨</b>\n\n%s\n", hostname, entry.Message)
			if err := notification.SendTelegramMessage(message); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to send Telegram notification: %v\n", err)
			}
		}
		return nil
	}
	logger := zap.New(core, zap.Hooks(telegramHook))
	Sugar = logger.Sugar()
}
