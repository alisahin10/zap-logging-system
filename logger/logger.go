package logger

import (
	"fmt"
	"logging-system/logswitcher"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a single logger instance with two cores writing to specified directories based on a shared flag.
func NewLogger(directory1, directory2 string, flag *logswitcher.ActiveFlag) (*zap.Logger, *dynamicFileSyncer, *dynamicFileSyncer, error) {
	// Generate unique filenames based on current timestamp for both cores
	filename1 := filepath.Join(directory1, time.Now().Format("2006_01_02_15_04_05")+".log")
	filename2 := filepath.Join(directory2, time.Now().Format("2006_01_02_15_04_05")+".log")

	// Create or open log files with append, create, and write-only permissions
	file1, err := os.OpenFile(filename1, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create log file for core1: %w", err)
	}

	file2, err := os.OpenFile(filename2, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file1.Close()
		return nil, nil, nil, fmt.Errorf("failed to create log file for core2: %w", err)
	}

	// Create thread-safe file syncers
	syncer1 := newDynamicFileSyncer(file1)
	syncer2 := newDynamicFileSyncer(file2)

	// Configure the two cores with conditional level enablers
	core1 := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(syncer1),
		NewConditionalLevelEnabler(zapcore.InfoLevel, flag, true),
	)

	core2 := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(syncer2),
		NewConditionalLevelEnabler(zapcore.InfoLevel, flag, false),
	)

	// Combine cores into a single logger using Tee
	tee := zapcore.NewTee(core1, core2)

	// Create the logger with caller information
	logger := zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger, syncer1, syncer2, nil
}

// CreateNewFileForCore creates a new file for core1 and updates its syncer.
func CreateNewFileForCore(directory string, syncer *dynamicFileSyncer) error {
	filename := filepath.Join(directory, time.Now().Format("2006_01_02_15_04_05")+".log")
	newFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file for core1: %w", err)
	}
	syncer.ChangeFile(newFile)
	return nil
}
