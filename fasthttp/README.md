# FastHTTP Server

A high-performance HTTP server built with [fasthttp](https://github.com/valyala/fasthttp) and [Cobra CLI](https://github.com/spf13/cobra).

## Installation

```bash
go mod tidy
go build -o fasthttp-server .
```

## Usage

### Basic Commands

```bash
# Show help
./fasthttp-server --help

# Show server command help
./fasthttp-server server --help
```

### Starting the Server

```bash
# Start server on default port 8080
./fasthttp-server server

# Start server on custom port
./fasthttp-server server --port 3000
./fasthttp-server server -p 3000

# Start server on custom host and port
./fasthttp-server server --host 0.0.0.0 --port 8080
./fasthttp-server server -H 0.0.0.0 -p 8080
```

### Server Flags

- `--port, -p`: Port to run the server on (default: 8080)
- `--host, -H`: Host address to bind the server to (default: localhost)

## License

This project is a demonstration of fasthttp and Cobra CLI integration. 