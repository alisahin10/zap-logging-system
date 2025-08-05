# Zap Logging System

This project is a high-performance, thread-safe logging system built in Go using the `zap` library. It features a single logger with multiple cores, controlled by an atomic boolean flag (`ActiveFlag`) for conditional logging. Logs are written to timestamped files in separate directories, with dynamic file switching handled in a thread-safe manner. The system is designed to be modular and extensible, making it easy to add new cores or integrate with external systems like Sentry.

## Features

- **Dual-Core Logger**: A single `zap.Logger` with two cores, each writing to distinct directories (`./logs/core1` and `./logs/core2`) based on an atomic flag.
- **Conditional Logging**: Uses `ConditionalLevelEnabler` to enable logging only when the log level and flag state match.
- **Dynamic File Switching**: `dynamicFileSyncer` supports runtime file switching with thread-safe operations, preventing resource leaks.
- **Modular Design**: Structured for easy extension, allowing new cores or conditions to be added with minimal changes.
- **Time-Based Log Rotation**: Logs are written to timestamped files, with Core1 switching to a new file at a specified interval.

## Technologies Used

- **Go**: Version 1.16 or higher.
- **Zap**: High-performance logging library (`go.uber.org/zap`).
- **Atomic Operations**: Thread-safe boolean flag management using `sync/atomic`.
- **File I/O**: Thread-safe file operations with `os` and `sync.Mutex`.

## Project Structure
<img width="824" height="265" alt="Screenshot 2025-08-05 at 21 20 51" src="https://github.com/user-attachments/assets/82d3250d-d6fc-4800-a268-5d788c0ff9cc" />

## Logging Behaviour
The system logs messages in a 15-second cycle with the following behavior:

0-5 seconds: Core1 writes to ./logs/core1/<timestamp>.log (flag = true).
5-10 seconds: Core2 writes to ./logs/core2/<timestamp>.log (flag = false).
10-15 seconds: Core1 writes to a new file ./logs/core1/<new_timestamp>.log (flag = true).

## Example Log Files
```bash
{"level":"info","caller":"main.go:XX","msg":"Log entry","log":0}
{"level":"info","caller":"main.go:XX","msg":"Log entry","log":1}
...
{"level":"info","caller":"main.go:XX","msg":"Log entry","log":4}
```

## Example Terminal Output
```bash
Core1 wrote log 0 to ./logs/core1/2025_08_05_21_06_05.log
Core1 wrote log 1 to ./logs/core1/2025_08_05_21_06_05.log
...
Core1 wrote log 4 to ./logs/core1/2025_08_05_21_06_05.log
Switched to Core2 at log 5, Core1 logs synced to ./logs/core1/2025_08_05_21_06_05.log
Core2 wrote log 5 to ./logs/core2/2025_08_05_21_06_05.log
...
Core2 wrote log 9 to ./logs/core2/2025_08_05_21_06_05.log
Switched to Core1 with new file at log 10, writing to ./logs/core1/2025_08_05_21_06_10.log
Core1 wrote log 10 to ./logs/core1/2025_08_05_21_06_10.log
...
Core1 wrote log 14 to ./logs/core1/2025_08_05_21_06_10.log
All logs synced and file handles closed
```
