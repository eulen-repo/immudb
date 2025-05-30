# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Main Build Commands
- `make all` - Build all main binaries (immudb, immuclient, immuadmin, immutest)
- `make clean` - Remove built binaries and generated files
- `make rebuild` - Clean and rebuild with code generation

### Individual Component Builds
- `make immudb` - Build the database server (with webconsole and swagger if enabled)
- `make immuclient` - Build the CLI client
- `make immuadmin` - Build the admin tool
- `make immutest` - Build the test utility

### Static Builds (cross-platform, no external dependencies)
- `make immudb-static` - Static build of database server
- `make immuclient-static` - Static build of client
- `make immuadmin-static` - Static build of admin tool

### Testing and Quality
- `make test` - Run full test suite with vet checks (sets LOG_LEVEL=error)
- `make test-client` - Run client package tests only
- `make coverage` - Generate test coverage report (requires go-acc tool)

### Code Generation
- `make build/codegen` - Generate protobuf files and gRPC code
- `make build/codegenv2` - Generate v2 protobuf files for authorization/documents APIs

### Web Console and Swagger
- `make webconsole` - Generate webconsole (use WEBCONSOLE=default for embedding)
- `make swagger` - Generate swagger UI (use SWAGGER=true to enable)

### Development Setup
- `make prerequisites` - Install required tools and tidy modules
- `make vendor` - Create vendor directory

## Project Architecture

### Main Executables (cmd/)
- **immudb** - Core database server with gRPC API, web console, and PostgreSQL compatibility
- **immuclient** - Interactive CLI client for database operations
- **immuadmin** - Administrative tool for server management, backups, and user management
- **immutest** - Testing and benchmarking utility

### Core Packages (pkg/)
- **database/** - Main database engine, multi-database management, SQL/KV operations
- **server/** - gRPC server implementation, session management, authentication/authorization
- **client/** - Client SDK with verification, caching, and streaming capabilities
- **auth/** - Authentication, authorization, user management, and security interceptors

### Embedded Storage Engine (embedded/)
- **store/** - Core immutable storage engine with cryptographic proofs
- **sql/** - SQL engine with parser, query execution, and transaction support
- **tbtree/** - Time-based B-tree for indexing with history
- **ahtree/** - Authenticated hash tree for cryptographic verification
- **appendable/** - Append-only file storage with multi-app and remote storage support

### Key Features
- **Immutable Storage** - All data is cryptographically verifiable and tamper-evident
- **Multi-Model** - Supports key-value, document, and relational (SQL) data models
- **PostgreSQL Wire Protocol** - Compatible with PostgreSQL clients via pkg/pgsql/
- **Streaming Replication** - Master-follower replication with pkg/replication/
- **Remote Storage** - S3-compatible storage backends via embedded/remotestorage/
- **Audit and Compliance** - Built-in auditing and tamper detection

### Protocol Definitions (pkg/api/)
- **schema/** - Main gRPC API definitions and protobuf schemas
- **proto/** - V2 API for authorization and document operations

### Configuration
- Default config files in configs/ directory
- Environment variable support for all settings
- S3 storage configuration via IMMUDB_S3_* variables

## Test Commands

Run single test file:
```bash
go test ./path/to/package -run TestName
```

Run tests with verbose output:
```bash
go test -v ./pkg/client ${GO_TEST_FLAGS}
```

Run integration tests:
```bash
go test -v ./pkg/integration/...
```

## Claude Code Guidelines

When generating files with Claude Code, place output files in the `claude-out/` directory whenever it makes sense (e.g., scripts, documentation, analysis reports, or any generated content that isn't part of the core codebase).