package main

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	debugLevels := map[string]slog.Level{"DEBUG": -4, "INFO": 0, "WARN": 4, "ERROR": 8}

	logLevel := GetEnvOrDefault("LOG_LEVEL", "INFO")

	options := slog.HandlerOptions{Level: debugLevels[logLevel]}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &options))

	return logger
}
