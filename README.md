# pflag Demo

A simple Go application demonstrating the use of [pflag](https://github.com/spf13/pflag) and [Cobra](https://github.com/spf13/cobra) for handling command-line flags, specifically for setting log levels.

## Features

- **Multiple Log Levels**: Supports DEBUG, INFO, WARN, and ERROR log levels
- **Convenient Flag Shortcuts**: Use `--verbose` for debug level or `--quiet` for error level
- **Conflict Detection**: Prevents using conflicting flags together (e.g., `--verbose` and `--quiet`)
- **Input Validation**: Validates log level input and provides helpful error messages
- **Demo Mode**: Demonstrates different log outputs based on the selected level

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd go-k8s-controller

# Build the application
go build -o pflag-demo

# Run the demo
./pflag-demo
```

## Usage

### Basic Usage

```bash
# Run with default settings (INFO level)
./pflag-demo

# Run with DEBUG level logging
./pflag-demo --log-level debug

# Run with verbose mode (equivalent to debug level)
./pflag-demo --verbose

# Run in quiet mode (error level only)
./pflag-demo --quiet
```

### Command Line Options

- `--log-level, -l`: Set the log level (debug, info, warn, error) - default: "info"
- `--verbose, -v`: Enable verbose logging (equivalent to --log-level=debug)
- `--quiet, -q`: Enable quiet mode (equivalent to --log-level=error)

### Log Levels

The application supports four log levels in order of verbosity:

1. **DEBUG** - Most detailed logging including all message types
2. **INFO** - General operational messages, warnings, and errors
3. **WARN** - Warning messages and errors only
4. **ERROR** - Error messages only

When you set a log level, only messages at that level and above will be displayed in the demo.

## Examples

```bash
# Show only error messages
./pflag-demo --quiet
./pflag-demo --log-level error

# Show all log messages including debug
./pflag-demo --verbose
./pflag-demo --log-level debug

# Show warnings and errors
./pflag-demo --log-level warn

# Use short flags
./pflag-demo -v  # verbose mode
./pflag-demo -q  # quiet mode
./pflag-demo -l info  # info level
```

## Flag Conflicts

The application prevents conflicting flags from being used together:

```bash
# This will result in an error
./pflag-demo --verbose --quiet
```

## Output Example

When run with different log levels, the application will show:

### Debug Level (--verbose)
```
Log level set to: DEBUG

--- Application Demo ---
[DEBUG] Debug message: Application started with detailed logging
[INFO] Info message: Processing data...
[WARN] Warn message: This is a warning
[ERROR] Error message: This is an error
```

### Error Level (--quiet)
```
Log level set to: ERROR

--- Application Demo ---
[ERROR] Error message: This is an error
```

## Prerequisites

- Go 1.21 or later

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [pflag](https://github.com/spf13/pflag) - POSIX/GNU-style command-line flag parsing (included with Cobra)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).