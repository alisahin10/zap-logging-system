package logger

import (
	"logging-system/logswitcher"

	"go.uber.org/zap/zapcore"
)

// ConditionalLevelEnabler implements zapcore.LevelEnabler interface to provide conditional logging based on external flag states.
type ConditionalLevelEnabler struct {
	enabledLevel zapcore.Level           // Minimum log level that this enabler accepts
	flag         *logswitcher.ActiveFlag // Shared flag that controls activation
	desiredState bool                    // The flag state that enables this logger
}

// NewConditionalLevelEnabler creates a new conditional level enabler that will only enable logging when both the log level is sufficient AND the flag matches the desired state.
func NewConditionalLevelEnabler(level zapcore.Level, flag *logswitcher.ActiveFlag, desiredState bool) *ConditionalLevelEnabler {
	return &ConditionalLevelEnabler{
		enabledLevel: level,
		flag:         flag,
		desiredState: desiredState,
	}
}

// Enabled determines whether a log entry at the given level should be processed.
func (e *ConditionalLevelEnabler) Enabled(lvl zapcore.Level) bool {
	// Check both level requirement and flag state
	// Example: If desiredState=true and flag.IsActive()=true, then logging is enabled
	// Example: If desiredState=false and flag.IsActive()=false, then logging is enabled
	return lvl >= e.enabledLevel && e.flag.IsActive() == e.desiredState
}
