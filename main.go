package main

import (
	"fmt"
	"logging-system/logger"
	"logging-system/logswitcher"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func main() {
	// Initialize the shared flag (true means core1 is active initially)
	flag := logswitcher.NewActiveFlag(true)

	// Create directories if they don't exist
	os.MkdirAll("./logs/core1", 0755)
	os.MkdirAll("./logs/core2", 0755)

	// Create a single logger with two cores
	log, syncer1, syncer2, err := logger.NewLogger("./logs/core1", "./logs/core2", flag)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}

	// Ensure proper cleanup
	defer func() {
		log.Sync()      // Flush pending logs
		syncer1.Close() // Close core1 file handle
		syncer2.Close() // Close core2 file handle
		fmt.Println("All logs synced and file handles closed")
	}()

	// Ticker for logging every second
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	// Main logging loop: runs for 15 seconds
	for i := 0; i < 15; i++ {
		// Switch cores at specified intervals
		if i == 5 {
			log.Sync()    // Ensure logs are flushed before switching
			flag.Toggle() // Switch to core2 (false)
			if fileName, ok := syncer1.GetCurrentFileName(); ok {
				fmt.Printf("Switched to Core2 at log %d, Core1 logs synced to %s\n", i, fileName)
			} else {
				fmt.Printf("Switched to Core2 at log %d, Core1 had no file\n", i)
			}
		} else if i == 10 {
			log.Sync()    // Ensure logs are flushed before switching
			flag.Toggle() // Switch back to core1 (true)
			// Create a new file for core1
			newFileName := filepath.Join("./logs/core1", time.Now().Format("2006_01_02_15_04_05")+".log")
			if err := logger.CreateNewFileForCore("./logs/core1", syncer1); err != nil {
				fmt.Printf("Failed to switch file for Core1 at log %d: %v\n", i, err)
			} else {
				fmt.Printf("Switched to Core1 with new file at log %d, writing to %s\n", i, newFileName)
			}
		}

		// Wait for the next tick
		<-ticker.C

		// Log to the single logger (active core writes based on flag)
		log.Info("Log entry", zap.Int("log", i))

		// Print which core is writing and to which file
		if flag.IsActive() {
			if fileName, ok := syncer1.GetCurrentFileName(); ok {
				fmt.Printf("Core1 wrote log %d to %s\n", i, fileName)
			} else {
				fmt.Printf("Core1 wrote log %d but no file is set\n", i)
			}
		} else {
			if fileName, ok := syncer2.GetCurrentFileName(); ok {
				fmt.Printf("Core2 wrote log %d to %s\n", i, fileName)
			} else {
				fmt.Printf("Core2 wrote log %d but no file is set\n", i)
			}
		}
	}
}
