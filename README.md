## loggo

**loggo** is a simple logging package for Go applications. It provides structured logging with multiple handlers and context support.

**Features:**

* Structured logging with JSON format
* Multiple handlers for logging to files and stdout
* Configurable log levels
* Context support for adding operation IDs to logs
* Helper function to generate unique operation IDs

**Installation:**

```bash
go get -u github.com/jonesrussell/loggo
```

**Usage:**

1. Create a new logger instance specifying the log file path:

```go
import (
  "github.com/jonesrussell/loggo"
)

logger, err := loggo.NewLogger("/path/to/logfile.log")
if err != nil {
  // handle error
}
```

2. Use the logger methods to log messages with different levels:

```go
logger.Debug("Starting application")
logger.Info("Processing request", "user", "john")
logger.Warn("Unexpected error", "error", err)
logger.Error("Failed to connect to database", err)
```

3. Add operation IDs to logs for context:

```go
operationID := loggo.NewOperationID()
logger = logger.WithOperation(operationID)

logger.Info("Processing payment", "operationID", operationID, "amount", 100)
```

**Dependencies:**

* [github.com/google/uuid](https://github.com/google/uuid)
* [github.com/samber/slog-multi](https://github.com/samber/slog-multi)
