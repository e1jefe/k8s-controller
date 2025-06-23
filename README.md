# Zerolog Demo

A simple Go application demonstrating the [zerolog](https://github.com/rs/zerolog) structured logging library with different log levels and formatting options.

## Features

- **Multiple Log Levels**: Supports TRACE, DEBUG, INFO, WARN, and ERROR log levels
- **Configurable Output**: Choose between structured JSON logging or pretty console output
- **Level Filtering**: Only displays messages at or above the specified log level
- **Command-line Interface**: Built with Cobra for clean CLI experience

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd go-k8s-controller

# Build the application
go build -o zerolog-demo

# Run the demo
./zerolog-demo
```

## Usage

### Basic Usage

```bash
# Run with default settings (INFO level)
./zerolog-demo

# Run with DEBUG level logging
./zerolog-demo --level debug

# Run with pretty console output
./zerolog-demo --pretty

# Run with TRACE level and pretty output
./zerolog-demo --level trace --pretty
```

### Command Line Options

- `--level, -l`: Set the log level (trace, debug, info, warn, error) - default: "info"
- `--pretty, -p`: Enable pretty console output with colors - default: false

### Log Levels

The application demonstrates five log levels in order of verbosity:

1. **TRACE** - Most detailed logging for debugging
2. **DEBUG** - Debug information
3. **INFO** - General operational messages
4. **WARN** - Warning messages
5. **ERROR** - Error conditions

When you set a log level, only messages at that level and above will be displayed.

## Examples

```bash
# Show only warnings and errors
./zerolog-demo --level warn

# Show all log messages with colorized output
./zerolog-demo --level trace --pretty

# Use short flags
./zerolog-demo -l debug -p
```

## Output Formats

### Structured JSON (default)
```json
{"level":"info","time":"2023-12-07T10:30:45Z","message":"This is an INFO message"}
```

### Pretty Console (with --pretty flag)
```
10:30:45 INF This is an INFO message
```

## Prerequisites

- Go 1.21 or later

## Dependencies

- [zerolog](https://github.com/rs/zerolog) - Structured logging library
- [cobra](https://github.com/spf13/cobra) - CLI framework

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).