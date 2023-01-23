package main

import (
	log "github.com/snirkop89/simplelogger"
)

// Examples
func main() {
	log.Println("this is a regular log")

	// Initialize a logger with specific options
	logger := log.New(log.FormatJSON, log.LevelInfo)

	logger.Info("Info message")
	logger.Warnf("Warning: %d", 10)
	logger.Error("This is an error")

	// With additional properties
	logger.WithFields("stage", "cleanup", "priority", "high").Info("Cleanup")

	// Providing invalid fields
	logger.WithFields("oneonly").Error("bad")
	logger.WithFields("oneonly", "two", "three").Error("bad")

	// #################################
	// Using console writer - human friendly logging
	logger = log.New(log.FormatHuman, log.LevelWarn)

	logger.Info("This won't be printed")
	logger.Warn("Only warning and above")
	logger.WithFields("stage", "cleanup").Error("failed cleanup")
}
