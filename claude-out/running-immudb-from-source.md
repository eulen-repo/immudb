# Running ImmuDB from Source

## Dependencies

**Go 1.17+** (Go 1.23.8 confirmed working)
- Project requires Go 1.17 as specified in `go.mod`
- All other dependencies are managed by Go modules
- No external dependencies required

## How to Run

### Method 1: Build and Run (Production-like)

1. **Build the server:**
```bash
# Basic build (without web console)
make immudb

# Build with web console support
make WEBCONSOLE=default
```

2. **Run the server** (multiple options):

**Development mode (simplest):**
```bash
./immudb --devmode
```
- Accepts remote connections without authentication
- Ideal for local development

**Foreground mode (to see logs):**
```bash
./immudb
```

**Background mode:**
```bash
./immudb -d
```

**With custom configuration:**
```bash
./immudb --port 3322 --dir ./mydata --auth
```

### Method 2: Direct Go Run (Better for Debugging)

You can run ImmuDB directly with `go run` for easier debugging:

**From the main directory:**
```bash
go run ./cmd/immudb
```

**With development mode:**
```bash
go run ./cmd/immudb --devmode
```

**With web console support (requires build first):**
```bash
# First build with webconsole, then run with go run
make WEBCONSOLE=default
go run -tags webconsole ./cmd/immudb --devmode
```

**With custom flags:**
```bash
go run ./cmd/immudb --port 3322 --dir ./mydata --devmode
```

**For debugging with delve:**
```bash
dlv debug ./cmd/immudb -- --devmode
```

## Default Ports

- **3322**: Main gRPC API
- **8080**: Web console
- **5432**: PostgreSQL wire protocol 
- **9497**: Prometheus metrics

## Testing the Connection

After running the server, you can test with:

```bash
# Build and test the client
make immuclient
./immuclient
```

Or run the client directly:
```bash
go run ./cmd/immuclient
```

## Notes

- ImmuDB is self-contained with no external dependencies beyond Go
- Creates its own data files in `./data` directory by default
- Use `--devmode` for development to skip authentication
- Use `go run` method for easier debugging and development iterations

## Web Console

The web console is available at `http://localhost:8080` when enabled:

- **To enable**: Build with `make WEBCONSOLE=default`
- **Default credentials**: username: `immudb`, password: `immudb`
- **Without webconsole**: You'll get "immudb was built without web console support" error

## Build Tags

- `webconsole`: Embeds the web console in the binary
- `swagger`: Enables Swagger UI
- Example: `go run -tags "webconsole,swagger" ./cmd/immudb --devmode`

## gRPC Access

ImmuDB primarily uses gRPC for client communication:

- **Default port**: 3322
- **Go SDK**: `github.com/codenotary/immudb/pkg/client`
- **CLI tools**: `immuclient`, `immuadmin`
- **Authentication**: Session-based with username/password