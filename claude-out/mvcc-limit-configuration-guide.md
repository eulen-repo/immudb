# How to Configure MVCC Read Set Limit

## Summary
The MVCC read set limit **cannot be configured via environment variables or command line flags**. It must be configured via gRPC using the `UpdateDatabaseV2` method after the database is created.

## Complete Go SDK Solution (Recommended)

### Step 1: Create the Configuration Tool

Create a file called `mvcc-config.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "strconv"
    
    "github.com/codenotary/immudb/pkg/client"
    "github.com/codenotary/immudb/pkg/api/schema"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: ./mvcc-config <new-limit>")
        fmt.Println("Example: ./mvcc-config 500000")
        os.Exit(1)
    }
    
    // Parse the new limit from command line
    newLimit, err := strconv.ParseUint(os.Args[1], 10, 32)
    if err != nil {
        log.Fatalf("Invalid limit value: %v", err)
    }
    
    // Create client
    c := client.NewClient()
    
    // Configure connection options
    opts := client.DefaultOptions().
        WithAddress("localhost").
        WithPort(3322)
    c.WithOptions(opts)
    
    ctx := context.Background()
    
    // Open session with ImmuDB
    fmt.Println("Connecting to ImmuDB...")
    err = c.OpenSession(ctx, []byte("immudb"), []byte("immudb"), "defaultdb")
    if err != nil {
        log.Fatalf("Failed to connect to ImmuDB: %v", err)
    }
    defer c.CloseSession(ctx)
    
    fmt.Println("✅ Connected successfully!")
    
    // Get current database settings (BEFORE)
    fmt.Println("\n🔍 Checking current configuration...")
    currentSettings, err := c.GetDatabaseSettingsV2(ctx)
    if err != nil {
        log.Fatalf("Failed to get current settings: %v", err)
    }
    
    // Show current MVCC limit
    if currentSettings.Settings.MvccReadSetLimit != nil {
        fmt.Printf("📊 Current MVCC Read Set Limit: %d\n", currentSettings.Settings.MvccReadSetLimit.Value)
    } else {
        fmt.Printf("📊 Current MVCC Read Set Limit: 100,000 (default)\n")
    }
    
    // Update MVCC Read Set Limit
    fmt.Printf("\n🔧 Setting new MVCC Read Set Limit to: %d\n", newLimit)
    
    newSettings := &schema.DatabaseNullableSettings{
        MvccReadSetLimit: &schema.NullableUint32{Value: uint32(newLimit)},
    }
    
    response, err := c.UpdateDatabaseV2(ctx, "defaultdb", newSettings)
    if err != nil {
        log.Fatalf("Failed to update MVCC limit: %v", err)
    }
    
    fmt.Printf("✅ MVCC Read Set Limit updated successfully on database: %s\n", response.Database)
    
    // Verify the change (AFTER)
    fmt.Println("\n🔍 Verifying new configuration...")
    updatedSettings, err := c.GetDatabaseSettingsV2(ctx)
    if err != nil {
        log.Fatalf("Failed to verify settings: %v", err)
    }
    
    if updatedSettings.Settings.MvccReadSetLimit != nil {
        fmt.Printf("📊 New MVCC Read Set Limit: %d\n", updatedSettings.Settings.MvccReadSetLimit.Value)
        
        if updatedSettings.Settings.MvccReadSetLimit.Value == uint32(newLimit) {
            fmt.Println("✅ Configuration verified successfully!")
        } else {
            fmt.Println("❌ Configuration verification failed!")
        }
    } else {
        fmt.Println("❌ Failed to verify new configuration")
    }
    
    fmt.Println("\n🎉 MVCC Read Set Limit configuration completed!")
}
```

### Step 2: Setup and Build

```bash
# Create a new directory for the tool
mkdir mvcc-config-tool
cd mvcc-config-tool

# Initialize Go module
go mod init mvcc-config

# Download dependencies
go get github.com/codenotary/immudb/pkg/client
go get github.com/codenotary/immudb/pkg/api/schema

# Copy the code above into mvcc-config.go
# Then build the tool
go build -o mvcc-config mvcc-config.go
```

### Step 3: Start Your ImmuDB Container

```bash
# Your existing command
sudo docker run --detach \
  -v "immudb_pix2depixd_$env:/var/lib/immudb" \
  -v "/home/azure/backup:/home/azure/backup" \
  --net host \
  -it \
  --name immudb \
  codenotary/immudb:latest

# Wait for ImmuDB to be ready
sleep 10
```

### Step 4: Run the Configuration Tool

```bash
# Check current configuration and set new limit to 500,000
./mvcc-config 500000
```

### Expected Output

```
Connecting to ImmuDB...
✅ Connected successfully!

🔍 Checking current configuration...
📊 Current MVCC Read Set Limit: 100,000 (default)

🔧 Setting new MVCC Read Set Limit to: 500000
✅ MVCC Read Set Limit updated successfully on database: defaultdb

🔍 Verifying new configuration...
📊 New MVCC Read Set Limit: 500000
✅ Configuration verified successfully!

🎉 MVCC Read Set Limit configuration completed!
```

### Step 5: Create a Check-Only Tool (Optional)

Create `check-mvcc.go` to just view current settings:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/codenotary/immudb/pkg/client"
)

func main() {
    // Create client
    c := client.NewClient()
    opts := client.DefaultOptions().WithAddress("localhost").WithPort(3322)
    c.WithOptions(opts)
    
    ctx := context.Background()
    
    // Connect
    err := c.OpenSession(ctx, []byte("immudb"), []byte("immudb"), "defaultdb")
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer c.CloseSession(ctx)
    
    // Get current settings
    settings, err := c.GetDatabaseV2(ctx, "defaultdb")
    if err != nil {
        log.Fatalf("Failed to get settings: %v", err)
    }
    
    // Display current MVCC limit
    fmt.Println("📊 ImmuDB Configuration:")
    fmt.Println("========================")
    
    if settings.Settings.MvccReadSetLimit != nil {
        fmt.Printf("MVCC Read Set Limit: %d\n", settings.Settings.MvccReadSetLimit.Value)
    } else {
        fmt.Printf("MVCC Read Set Limit: 100,000 (default)\n")
    }
    
    // Show other relevant settings
    if settings.Settings.MaxConcurrency != nil {
        fmt.Printf("Max Concurrency: %d\n", settings.Settings.MaxConcurrency.Value)
    }
    
    if settings.Settings.MaxIOConcurrency != nil {
        fmt.Printf("Max IO Concurrency: %d\n", settings.Settings.MaxIOConcurrency.Value)
    }
    
    if settings.Settings.SyncFrequency != nil {
        fmt.Printf("Sync Frequency: %s\n", settings.Settings.SyncFrequency.Value)
    }
}
```

```bash
# Build check tool
go build -o check-mvcc check-mvcc.go

# Use it to check current configuration anytime
./check-mvcc
```

## Alternative Methods

### Option 2: REST API (if you enable web server)

### Enable Web Server in Your Docker Command
```bash
sudo docker run --detach \
  -v "immudb_pix2depixd_$env:/var/lib/immudb" \
  -v "/home/azure/backup:/home/azure/backup" \
  --net host \
  -it \
  --name immudb \
  -e IMMUDB_WEB_SERVER=true \
  -e IMMUDB_WEB_SERVER_PORT=8080 \
  codenotary/immudb:latest
```

### Configure via REST API
```bash
# 1. Login to get token
TOKEN=$(curl -s -X POST "http://localhost:8080/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"user": "immudb", "password": "immudb"}' | \
  grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 2. Update MVCC limit
curl -X POST "http://localhost:8080/db/update/v2" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "database": "defaultdb",
    "settings": {
      "mvccReadSetLimit": 500000
    }
  }'
```

### Option 3: Automated Setup Script

Create `setup-immudb-with-mvcc.sh`:

```bash
#!/bin/bash
set -e

ENV=${1:-"dev"}
MVCC_LIMIT=${2:-500000}

echo "Setting up ImmuDB with MVCC limit: $MVCC_LIMIT"

# Start ImmuDB
echo "Starting ImmuDB container..."
sudo docker run --detach \
  -v "immudb_pix2depixd_$ENV:/var/lib/immudb" \
  -v "/home/azure/backup:/home/azure/backup" \
  --net host \
  -it \
  --name immudb \
  codenotary/immudb:latest

# Wait for ImmuDB to be ready
echo "Waiting for ImmuDB to be ready..."
until nc -z localhost 3322; do
    echo "  Still waiting..."
    sleep 2
done

echo "ImmuDB is ready! Configuring MVCC limit..."

# Create temporary Go configuration
cat > /tmp/config-mvcc.go << EOF
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/codenotary/immudb/pkg/client"
    "github.com/codenotary/immudb/pkg/api/schema"
)

func main() {
    c := client.NewClient()
    opts := client.DefaultOptions().WithAddress("localhost").WithPort(3322)
    c.WithOptions(opts)
    
    ctx := context.Background()
    err := c.OpenSession(ctx, []byte("immudb"), []byte("immudb"), "defaultdb")
    if err != nil {
        log.Fatal(err)
    }
    defer c.CloseSession(ctx)
    
    // Get current settings
    current, _ := c.GetDatabaseV2(ctx, "defaultdb")
    if current.Settings.MvccReadSetLimit != nil {
        fmt.Printf("Current MVCC limit: %d\n", current.Settings.MvccReadSetLimit.Value)
    } else {
        fmt.Println("Current MVCC limit: 100,000 (default)")
    }
    
    // Update settings
    settings := &schema.DatabaseNullableSettings{
        MvccReadSetLimit: &schema.NullableUint32{Value: $MVCC_LIMIT},
    }
    
    _, err = c.UpdateDatabaseV2(ctx, "defaultdb", settings)
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
    
    fmt.Println("MVCC Read Set Limit updated to $MVCC_LIMIT successfully!")
}
EOF

# Run configuration
cd /tmp
go mod init config
go get github.com/codenotary/immudb/pkg/client
go get github.com/codenotary/immudb/pkg/api/schema
go run config-mvcc.go

# Cleanup
rm -f config-mvcc.go go.mod go.sum

echo "✅ ImmuDB setup completed with MVCC limit: $MVCC_LIMIT"
```

### Usage
```bash
# Make executable
chmod +x setup-immudb-with-mvcc.sh

# Run with default limit (500,000)
./setup-immudb-with-mvcc.sh dev

# Run with custom limit
./setup-immudb-with-mvcc.sh dev 1000000
```

## Quick Reference Commands

### Build the tools once
```bash
# Create tools directory
mkdir immudb-tools && cd immudb-tools

# Setup Go module
go mod init immudb-tools
go get github.com/codenotary/immudb/pkg/client
go get github.com/codenotary/immudb/pkg/api/schema

# Copy the mvcc-config.go and check-mvcc.go code above
# Build both tools
go build -o mvcc-config mvcc-config.go
go build -o check-mvcc check-mvcc.go
```

### Daily usage
```bash
# Check current configuration
./check-mvcc

# Change MVCC limit
./mvcc-config 500000

# Check again to verify
./check-mvcc
```

### Troubleshooting

**If connection fails:**
```bash
# Check if ImmuDB is running
sudo docker ps | grep immudb

# Check if port 3322 is accessible
nc -zv localhost 3322

# Check ImmuDB logs
sudo docker logs immudb
```

**If authentication fails:**
```bash
# Make sure you're using correct credentials
# Default: username="immudb", password="immudb"

# Check if auth is enabled in ImmuDB
sudo docker logs immudb | grep -i auth
```

## Key Points

1. **No CLI support** - Cannot be done via immuclient/immuadmin
2. **No environment variables** - Must use gRPC/REST API  
3. **Per-database setting** - Each database has its own limit
4. **Requires authentication** - Need admin permissions
5. **Persistent** - Setting is saved in database metadata

## Recommended Values

| Use Case | Recommended Limit | Reasoning |
|----------|------------------|-----------|
| Small datasets (<10k records) | 100,000 (default) | Default is sufficient |
| Medium datasets (10k-100k records) | 250,000 - 500,000 | Buffer for growth |
| Large datasets (100k+ records) | 500,000 - 1,000,000 | Accommodate large queries |
| Very large datasets (1M+ records) | 1,000,000+ | Based on available RAM |

**Remember:** Higher limits use more memory. Monitor your system resources!