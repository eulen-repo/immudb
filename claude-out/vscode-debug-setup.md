# VS Code Debug Setup for ImmuDB

## Prerequisites

1. **Go extension for VS Code** - Install the official Go extension by Google
2. **Delve debugger** - Usually installed automatically with the Go extension
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   ```

## Debug Configurations

A `.vscode/launch.json` file has been created with the following debug configurations:

### 1. Debug ImmuDB Server (Development Mode)
- **Name**: `Debug ImmuDB Server`
- **Args**: `--devmode --dir ./debug-data`
- **Build Flags**: `-tags webconsole` (enables web console)
- **Usage**: Main development debugging - no authentication required

### 2. Debug ImmuDB Server (Production Mode)
- **Name**: `Debug ImmuDB Server (Production Mode)`
- **Args**: `--auth --dir ./debug-data`
- **Build Flags**: `-tags webconsole` (enables web console)
- **Usage**: Debug with authentication enabled

### 3. Debug ImmuDB Server (Custom Port)
- **Name**: `Debug ImmuDB Server (Custom Port)`
- **Args**: `--devmode --port 3323 --dir ./debug-data`
- **Build Flags**: `-tags webconsole` (enables web console)
- **Usage**: Debug on different port to avoid conflicts

### 4. Debug ImmuClient
- **Name**: `Debug ImmuClient`
- **Usage**: Debug the CLI client

### 5. Debug ImmuAdmin
- **Name**: `Debug ImmuAdmin`
- **Usage**: Debug the admin tool

### 6. Debug ImmuTest
- **Name**: `Debug ImmuTest`
- **Usage**: Debug the test utility

## How to Use

1. **Open VS Code** in the immudb project root
2. **Go to Run and Debug panel** (Ctrl+Shift+D)
3. **Select a configuration** from the dropdown
4. **Set breakpoints** in your code
5. **Press F5** or click the green play button

## Tips

- **Data Directory**: All configurations use `./debug-data` to avoid conflicts with production data
- **Development Mode**: Use `Debug ImmuDB Server` configuration for most debugging - it disables authentication
- **Web Console**: All configurations include web console support - access at `http://localhost:8080`
- **Multiple Instances**: Use the custom port configuration if you need to run multiple instances
- **Breakpoints**: Set breakpoints in any Go file before starting debug session
- **Variables**: Inspect variables, call stack, and goroutines in the debug panel
- **Console**: Use the integrated terminal to see server output

## Environment Variables

You can add environment variables to any configuration by adding an `env` property:

```json
{
    "name": "Debug ImmuDB Server",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}/cmd/immudb",
    "args": ["--devmode"],
    "buildFlags": "-tags webconsole",
    "env": {
        "LOG_LEVEL": "debug"
    },
    "console": "integratedTerminal"
}
```

Note: ImmuDB uses port 3322 by default for gRPC, and 8080 for web console when enabled.

## Debugging Tests

To debug specific tests, you can create additional configurations or use the built-in Go test debugging:

1. **Open a test file**
2. **Click "debug test"** above any test function
3. **Or use Command Palette**: `Go: Debug Test At Cursor`