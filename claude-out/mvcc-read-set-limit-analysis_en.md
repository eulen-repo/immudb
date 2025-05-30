# MVCC Read-Set Limit Analysis in ImmuDB

## Executive Summary

The "MVCC read-set limit exceeded" error occurs when a transaction performs more than **100,000 read operations**. This limit exists to prevent memory exhaustion and ensure MVCC validation remains feasible. Large tables (100k+ records) commonly trigger this error due to table scans that exceed the read operation limit.

## 1. Error Definition and Location

**File:** `embedded/store/immustore.go`
```go
var ErrMVCCReadSetLimitExceeded = errors.New("MVCC read-set limit exceeded")
```

**File:** `embedded/store/options.go`
```go
const DefaultMVCCReadSetLimit = 100_000
```

## 2. Understanding MVCC Read-Set

### What is the Read-Set?

The MVCC read-set tracks every read operation within a transaction to ensure **Snapshot Isolation**. It maintains:

**File:** `embedded/store/ongoing_tx.go`
```go
type mvccReadSet struct {
    expectedGets           []expectedGet            // Individual key reads
    expectedGetsWithPrefix []expectedGetWithPrefix  // Prefix-based reads  
    expectedReaders        []*expectedReader        // Key readers (scanners)
    readsetSize            int                      // Total count of operations
}
```

### What Counts as a Read Operation?

1. **Single key reads** (`Get` operations)
2. **Prefix-based reads** (`GetWithPrefix` operations)
3. **Key reader operations** (table scans/iterations)
4. **Reader resets** (when a scanner is reset during iteration)

Each operation consumes **one slot** in the read-set counter.

## 3. Where the Limit is Enforced

### Primary Check Function
**File:** `embedded/store/ongoing_tx.go`
```go
func (tx *OngoingTx) mvccReadSetLimitReached() bool {
    return tx.mvccReadSet.readsetSize == tx.st.mvccReadSetLimit
}
```

### Enforcement Points
1. **Get operations** (lines 457, 476 in `ongoing_tx.go`)
2. **GetWithPrefix operations** (lines 516, 537 in `ongoing_tx.go`)
3. **NewKeyReader creation** (line 60 in `ongoing_tx_keyreader.go`)
4. **Key reader reads** (lines 115, 134 in `ongoing_tx_keyreader.go`)
5. **Reader resets** (line 170 in `ongoing_tx_keyreader.go`)

## 4. Why Large Tables Trigger This Error

### Common Scenarios with 100k+ Record Tables

1. **Full Table Scans**
   - `SELECT * FROM large_table` reads every row
   - Each row access = 1 read operation
   - 100k+ rows > 100k limit

2. **Unindexed WHERE Clauses**
   - `SELECT * FROM table WHERE non_indexed_column = value`
   - Forces table scan to find matching rows

3. **Complex Queries**
   - JOINs between large tables
   - Aggregations over large datasets
   - DISTINCT operations on large result sets

4. **UPDATE/DELETE with Complex WHERE**
   - Even with PK, if query plan performs scan
   - Multiple index lookups in complex conditions

### Your Specific Case: UPDATE with PK

Even an `UPDATE WHERE pk = value` can exceed the limit if:
1. **Query planner chooses suboptimal path**
2. **Multiple indexes are consulted**
3. **Transaction has accumulated reads from previous operations**
4. **Table statistics are outdated**

## 5. Configuration Options

### Database Level Configuration

**File:** `pkg/api/schema/schema.proto`
```protobuf
// Limit the number of read entries per transaction
NullableUint32 mvccReadSetLimit = 27;
```

**File:** `pkg/server/db_options.go`
```go
MVCCReadSetLimit: store.DefaultMVCCReadSetLimit,  // Default: 100,000
```

### Runtime Configuration

**File:** `embedded/store/options.go`
```go
func (opts *Options) WithMVCCReadSetLimit(mvccReadSetLimit int) *Options {
    opts.MVCCReadSetLimit = mvccReadSetLimit
    return opts
}
```

### How to Increase the Limit

**⚠️ IMPORTANT: MVCC Read Set Limit is NOT configurable via:**
- ❌ Environment variables
- ❌ Command line flags  
- ❌ Configuration files
- ❌ SQL DDL statements

**✅ Must be configured via gRPC using `UpdateDatabaseV2`:**

```go
// Via Go SDK after database creation
settings := &schema.DatabaseNullableSettings{
    MvccReadSetLimit: &schema.NullableUint32{Value: 500000},
}
response, err := client.UpdateDatabaseV2(ctx, "defaultdb", settings)
```

See `mvcc-limit-configuration-guide.md` for complete examples.

## 6. Side Effects of Increasing the Limit

### Memory Impact

| Limit Value | Estimated Memory per Transaction | Risk Level |
|-------------|----------------------------------|------------|
| 100,000 (default) | ~10-50 MB | Low |
| 500,000 | ~50-250 MB | Medium |
| 1,000,000 | ~100-500 MB | High |
| 10,000,000 | ~1-5 GB | Very High |

### Performance Impact

1. **Commit Time**: Longer MVCC validation with larger read-sets
2. **Memory Usage**: Each read operation stores metadata
3. **Garbage Collection**: More objects to clean up
4. **Concurrency**: More memory pressure affects other transactions

### Stability Risks

1. **Out of Memory**: Very large read-sets can exhaust RAM
2. **Timeout Issues**: Validation may exceed timeouts
3. **System Responsiveness**: High memory usage affects overall performance

## 7. Can You Increase Infinitely?

**No.** Practical limits exist:

### Hard Limits
1. **Available RAM**: Each read consumes memory
2. **Go runtime limits**: Slice size limitations
3. **Validation time**: MVCC validation becomes prohibitively expensive

### Recommended Approaches

**Instead of infinite increases:**

1. **Optimize Queries**
   - Use proper indexes
   - Add LIMIT clauses
   - Break large operations into chunks

2. **Transaction Design**
   - Smaller transaction scopes
   - Read-only transactions for large scans
   - Batch processing patterns

3. **Table Design**
   - Proper indexing strategy
   - Table partitioning (if supported)
   - Archive old data

## 8. Root Cause Analysis: Your 100k Table Issue

### Why It Started Happening

1. **Threshold Effect**: Read-set tracking overhead becomes significant
2. **Query Plan Changes**: Optimizer may choose different paths
3. **Index Fragmentation**: As table grows, index efficiency may degrade
4. **Memory Pressure**: System under higher load affects query planning

### Why Archiving Fixed It

1. **Reduced Table Size**: Fewer rows to potentially scan
2. **Better Index Efficiency**: Smaller indexes are more cache-friendly
3. **Fresh Statistics**: New table has accurate cardinality estimates
4. **Reduced Memory Pressure**: Less data in memory

## 9. Recommended Solutions

### Immediate Fixes

1. **Increase Limit Moderately**
   ```go
   // Increase to 500,000 (5x default)
   dbOptions.WithMVCCReadSetLimit(500000)
   ```

2. **Query Optimization**
   - Ensure proper indexes on WHERE clause columns
   - Use EXPLAIN to understand query plans
   - Add LIMIT clauses where appropriate

### Long-term Solutions

1. **Monitoring Setup**
   - Track transaction read-set sizes
   - Monitor memory usage patterns
   - Set up alerts for approaching limits

2. **Data Management Strategy**
   - Regular archiving of old data
   - Table partitioning if supported
   - Proper maintenance procedures

3. **Application Design**
   - Implement read-only transaction patterns for large scans
   - Use streaming/pagination for large result sets
   - Design transactions to minimize read operations

## 10. Monitoring and Debugging

### Key Metrics to Track

1. **Transaction read-set sizes**
2. **Memory usage per transaction**
3. **MVCC validation times**
4. **Query execution plans**

### Debug Information

When the error occurs, log:
- Current transaction read-set size
- Query being executed
- Table sizes involved
- Index usage statistics

This analysis should help you understand and resolve the MVCC read-set limit issues while making informed decisions about configuration changes.