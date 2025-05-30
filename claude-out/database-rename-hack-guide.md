# Database Rename Guide

## What This Does

Renames the physical database directory to bypass ImmuDB's protection of the `defaultdb` database, allowing you to modify its MVCC settings.

**Simple approach: Just rename the directory** - no source code changes needed!

## Simple Directory Rename Method

**For Docker Volumes:**
```bash
# Stop ImmuDB container
docker stop your_immudb_container

# Rename the database directory
sudo docker run --rm -v your_volume_name:/data alpine mv /data/defaultdb /data/mydatabase

# Verify the rename
sudo docker run --rm -v your_volume_name:/data alpine ls -la /data
```

**For Local Filesystem:**
```bash
# Stop ImmuDB server
./immudb stop  # or kill the process

# Rename directory
mv /path/to/immudb/data/defaultdb /path/to/immudb/data/mydatabase

# Restart ImmuDB
./immudb --devmode
```

### Step 2: Update Your Code

In your MVCC configuration tool, change the database name:

```go
// Before
response, err := c.UpdateDatabaseV2(ctx, "defaultdb", newSettings)

// After  
response, err := c.UpdateDatabaseV2(ctx, "mydatabase", newSettings)
```

## Example: Complete Process

```bash
# 1. Stop container
docker stop my_immudb

# 2. Rename database
sudo docker run --rm -v immudb_volume:/data alpine mv /data/defaultdb /data/mydatabase

# 3. Start container
docker start my_immudb

# 4. Test MVCC configuration
go run mvcc-config.go 500000
```

## Notes

- **Data preservation**: All your data remains intact during rename
- **Simple operation**: Just a filesystem rename, no data corruption risk
- **Reversible**: Can rename back to `defaultdb` if needed
- **No recompilation**: Works with standard ImmuDB binaries

## Rollback

If needed, just rename back:

```bash
sudo docker run --rm -v your_volume_name:/data alpine mv /data/mydatabase /data/defaultdb
```