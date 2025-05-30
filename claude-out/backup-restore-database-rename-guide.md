# Backup & Restore Method for Database Rename

## Overview

This method allows you to migrate data from `defaultdb` to a new database using ImmuDB's official backup and restore functionality. This is **safer than directory renaming** and works with any ImmuDB version.

## Method: Hot Backup + Hot Restore

### Why This Works

- `hot-backup` (dump) exports data from any database 
- `hot-restore` can restore to **any database name**, including new ones
- Database can stay online during backup
- Restore creates new database if it doesn't exist

## Step-by-Step Process

### Step 1: Create Backup from defaultdb

```bash
# Make sure ImmuDB is running
./immudb --devmode --web-server

# Login and backup defaultdb
./immuadmin login immudb
./immuadmin dump defaultdb_backup.bkp
```

This creates a backup file with all data from `defaultdb`.

### Step 2: Create New Database and Restore

```bash
# Restore to a new database name
./immuadmin hot-restore mydatabase -i defaultdb_backup.bkp

# Verify the new database exists
./immuadmin database list
```

### Step 3: Update Your Application

Change your application to use the new database:

```go
// Before
err = c.OpenSession(ctx, []byte("immudb"), []byte("immudb"), "defaultdb")

// After  
err = c.OpenSession(ctx, []byte("immudb"), []byte("immudb"), "mydatabase")
```

### Step 4: Configure MVCC Settings

Now you can modify MVCC settings on the new database:

```bash
cd ~/path/to/mvcc-config
go run . 500000  # This will now work!
```

## Complete Example

```bash
# 1. Start ImmuDB
./immudb --devmode --web-server

# 2. Login and backup
./immuadmin login immudb
./immuadmin dump my_data_backup.bkp

# 3. Restore to new database  
./immuadmin hot-restore mydatabase -i my_data_backup.bkp

# 4. Verify data
./immuadmin database list

# 5. Test MVCC configuration
cd ~/tools/mvcc-config
go run . 500000
```

## Docker Example

```bash
# 1. Backup from running container
docker exec immudb_container ./immuadmin login immudb
docker exec immudb_container ./immuadmin dump /data/backup.bkp

# 2. Restore to new database
docker exec immudb_container ./immuadmin hot-restore mydatabase -i /data/backup.bkp

# 3. Update your application to use "mydatabase"
```

## Advantages Over Directory Rename

✅ **Official ImmuDB functionality**  
✅ **Works with any ImmuDB version**  
✅ **No custom builds needed**  
✅ **Safer data migration**  
✅ **Can verify backup before restore**  
✅ **Database stays online during backup**  
✅ **Supported by Codenotary**  

## Command Reference

### Backup Commands

```bash
# Simple backup
./immuadmin dump backup.bkp

# Backup with progress
./immuadmin dump backup.bkp --progress-bar

# Backup to stdout (for piping)
./immuadmin dump -
```

### Restore Commands

```bash
# Restore to new database
./immuadmin hot-restore newdb -i backup.bkp

# Append to existing database  
./immuadmin hot-restore newdb -i backup.bkp --append

# Verify backup without restoring
./immuadmin hot-restore newdb -i backup.bkp --verify-only

# Force restore (skip transaction checks)
./immuadmin hot-restore newdb -i backup.bkp --force
```

## Cleanup (Optional)

After confirming everything works with the new database:

```bash
# Drop the old defaultdb (if desired)
./immuadmin database unload defaultdb
```

## Rollback

If something goes wrong, you still have:
- Original `defaultdb` intact
- Backup file as additional safety

```bash
# Switch back to defaultdb in your application
# Or restore backup to defaultdb if needed
./immuadmin hot-restore defaultdb -i backup.bkp --force
```

## Notes

- **Database names**: Can use any valid database name
- **Data integrity**: All transactions and history preserved  
- **Indexes**: Rebuilt automatically during restore
- **Performance**: May be slower than directory rename but much safer
- **Permissions**: User permissions need to be recreated for new database

This method is the **recommended approach** for production environments!