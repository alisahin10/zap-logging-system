package logger

import (
	"os"
	"sync"
)

// dynamicFileSyncer provides a thread-safe wrapper around os.File that supports dynamic file switching during runtime.
type dynamicFileSyncer struct {
	mu   sync.Mutex // Mutex to ensure thread-safe access to the file handle
	file *os.File   // Current file handle (can be nil)
}

// newDynamicFileSyncer creates a new file syncer with the provided initial file.
func newDynamicFileSyncer(file *os.File) *dynamicFileSyncer {
	return &dynamicFileSyncer{file: file}
}

// Write implements the io.Writer interface by writing data to the current file.
func (d *dynamicFileSyncer) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// If no file is currently set, silently succeed without writing
	// This allows the logger to continue functioning even during file transitions
	if d.file == nil {
		return 0, nil
	}

	// Delegate the actual write operation to the underlying file
	return d.file.Write(p)
}

// Sync implements the zapcore.WriteSyncer interface by flushing any buffered data to the underlying storage.
func (d *dynamicFileSyncer) Sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// If no file is currently set, consider sync successful
	if d.file == nil {
		return nil
	}

	// Flush any buffered data to the underlying file system
	return d.file.Sync()
}

// ChangeFile atomically switches to a new file handle, closing the previous one.
func (d *dynamicFileSyncer) ChangeFile(newFile *os.File) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Close the current file if it exists
	// This ensures proper resource cleanup and prevents file handle leaks
	if d.file != nil {
		d.file.Close()
	}

	// Atomically switch to the new file
	d.file = newFile
}

// Close cleanly shuts down the syncer by closing the current file and setting the internal handle to nil.
func (d *dynamicFileSyncer) Close() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Close the file if it exists and reset the handle
	if d.file != nil {
		d.file.Close()
		d.file = nil // Prevent further operations on the closed file
	}
}

// GetCurrentFileName returns the current file name and a boolean indicating if a file is set.
func (d *dynamicFileSyncer) GetCurrentFileName() (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.file == nil {
		return "", false
	}
	return d.file.Name(), true
}
