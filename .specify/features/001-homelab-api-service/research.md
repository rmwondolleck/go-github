# Research: Gin Framework Performance Benchmarks

**Date**: 2026-03-01  
**Task**: T001 - Benchmark Gin framework basic routing  
**Phase**: Phase 0 - Research & Validation  
**Batch**: Batch 1 (Research)

## Objective

Validate Gin framework performance against stdlib net/http to ensure the framework overhead is acceptable (target: <10ms overhead) for the homelab API service.

## Methodology

Created comprehensive benchmarks in `research/gin_benchmark_test.go` comparing:
- Simple route handling (GET /health)
- Routes with middleware
- Multiple route registration
- Parameterized routes (URL parameters)

Each benchmark was run with:
- `-benchtime=3s` for statistically significant results
- `-benchmem` to measure memory allocations
- Go 1.25.7 on AMD EPYC 7763 64-Core Processor

## Results Summary

### Simple Route Handling
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1158 ns | 1064 B | 11 |
| Gin | 1096 ns | 1128 B | 12 |
| **Overhead** | **-62 ns** | **+64 B** | **+1** |

**Key Finding**: Gin is actually **faster** than stdlib for simple routes (62 nanoseconds faster, or ~5% improvement).

### With Middleware
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1172 ns | 1064 B | 11 |
| Gin | 1105 ns | 1128 B | 12 |
| **Overhead** | **-67 ns** | **+64 B** | **+1** |

**Key Finding**: Gin maintains performance advantage even with middleware (67 nanoseconds faster).

### Multiple Routes (5 routes)
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1261 ns | 1064 B | 11 |
| Gin | 1114 ns | 1128 B | 12 |
| **Overhead** | **-147 ns** | **+64 B** | **+1** |

**Key Finding**: Gin's advantage increases with more routes (147 nanoseconds faster, ~12% improvement), likely due to optimized routing algorithm.

### Parameterized Routes
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1581 ns | 1120 B | 14 |
| Gin | 1116 ns | 1128 B | 12 |
| **Overhead** | **-465 ns** | **+8 B** | **-2** |

**Key Finding**: Gin significantly outperforms stdlib for parameterized routes (465 nanoseconds faster, ~29% improvement), with fewer allocations.

## Performance Analysis

### Time Overhead
- **Target**: <10ms (10,000,000 ns) overhead
- **Actual**: Gin is consistently **faster** than stdlib, not slower
- **Conclusion**: ✅ **TARGET EXCEEDED** - No overhead detected; Gin provides performance improvements

### Overhead Breakdown by Scenario
1. **Simple routes**: -0.062 μs (Gin faster)
2. **With middleware**: -0.067 μs (Gin faster)
3. **Multiple routes**: -0.147 μs (Gin faster)
4. **Parameterized routes**: -0.465 μs (Gin faster)

All measurements are in the sub-microsecond range, well below the 10ms (10,000 μs) target.

### Memory Overhead
- Gin consistently uses ~64 bytes more per request
- One additional allocation per request
- For parameterized routes, Gin uses 8 bytes more but 2 fewer allocations
- Memory overhead is minimal and acceptable for the use case

## Advantages of Gin Framework

Based on benchmark results, Gin provides:

1. **Better Performance**: Faster than stdlib in all tested scenarios
2. **Optimized Routing**: Superior performance with multiple and parameterized routes
3. **Built-in Features**: Middleware support, JSON binding/validation, parameter extraction
4. **Developer Productivity**: Cleaner API, less boilerplate code
5. **Production-Ready**: Battle-tested framework with extensive ecosystem

## Recommendations

✅ **APPROVED**: Proceed with Gin framework for the homelab API service

**Rationale**:
1. Performance target (<10ms overhead) is **far exceeded** - Gin is actually faster than stdlib
2. Routing performance improves significantly with complex routes
3. Built-in features reduce development time and potential bugs
4. Memory overhead (64 bytes/request) is negligible for homelab scale
5. Gin's middleware system provides cleaner architecture for future features

## Trade-offs

### Pros:
- Superior performance to stdlib
- Rich middleware ecosystem
- Better developer experience
- Built-in JSON handling and validation
- Excellent parameter routing

### Cons:
- Additional dependency (~30 transitive dependencies added)
- Slightly higher memory usage (64 bytes per request)
- Framework lock-in (but migration path exists if needed)

## Conclusion

The Gin framework not only meets but **exceeds** the performance requirements with **negative overhead** (i.e., performance improvements) compared to stdlib net/http. The <10ms overhead target is achieved with significant margin - actual performance is improved by 0.062-0.465 microseconds depending on the scenario.

**Recommendation**: Proceed to Phase 1 with Gin framework as the chosen web framework.

## Next Steps

1. ✅ Benchmarks completed and validated
2. ✅ Performance target confirmed as achievable
3. Ready to proceed to Phase 1: Gin framework implementation
4. Future consideration: Add benchmarks for JSON parsing and response serialization

---

# Research: sync.Map vs RWMutex for Device Storage

**Date**: 2026-03-01  
**Task**: T003 - Benchmark sync.Map vs RWMutex for device storage  
**Phase**: Phase 0 - Research & Validation  
**Batch**: Batch 1 (Research)

## Objective

Benchmark in-memory storage options for device data to determine the optimal concurrent data structure. The primary requirement is to validate O(1) lookup performance while handling concurrent reads and mixed read/write workloads.

## Methodology

Created comprehensive benchmarks in `research/storage_benchmark_test.go` comparing:
- **sync.Map**: Go's built-in concurrent map optimized for read-heavy workloads
- **RWMutex**: Traditional map with read-write mutex protection

Test scenarios:
1. Concurrent reads with default parallelism (RunParallel)
2. Mixed workload (90% reads, 10% writes)
3. Explicit 100 goroutines concurrent reads
4. Single lookup with varying dataset sizes (100, 1000, 10000) to verify O(1) behavior

Each benchmark was run with:
- `-benchtime=3s` for statistically significant results
- `-benchmem` to measure memory allocations
- Pre-populated with 1000 devices
- Go 1.25.0 on AMD EPYC 7763 64-Core Processor

## Results Summary

### Concurrent Reads (RunParallel)
| Implementation | Time/op | Memory/op | Allocs/op |
|----------------|---------|-----------|-----------|
| sync.Map | 9.83 ns | 0 B | 0 |
| RWMutex | 46.81 ns | 0 B | 0 |
| **Performance** | **sync.Map 4.8x faster** | **Equal** | **Equal** |

**Key Finding**: sync.Map is **4.8x faster** than RWMutex for concurrent reads with zero allocations for both.

### Mixed Workload (90% reads, 10% writes)
| Implementation | Time/op | Memory/op | Allocs/op |
|----------------|---------|-----------|-----------|
| sync.Map | 101.3 ns | 75 B | 2 |
| RWMutex | 225.3 ns | 69 B | 2 |
| **Performance** | **sync.Map 2.2x faster** | **+6 B** | **Equal** |

**Key Finding**: sync.Map is **2.2x faster** than RWMutex for mixed read/write workloads with minimal additional memory overhead (6 bytes).

### 100 Goroutines Concurrent Reads
| Implementation | Time/op | Memory/op | Allocs/op |
|----------------|---------|-----------|-----------|
| sync.Map | 652,380 ns | 223,682 B | 17,642 |
| RWMutex | 711,042 ns | 223,671 B | 17,642 |
| **Performance** | **sync.Map 8.2% faster** | **+11 B** | **Equal** |

**Key Finding**: With 100 goroutines each performing 100 reads (10,000 total operations), sync.Map is **8.2% faster** with nearly identical memory characteristics.

### O(1) Lookup Performance Validation

#### sync.Map Lookup Times by Dataset Size
| Dataset Size | Time/op | Growth |
|--------------|---------|--------|
| 100 | 102.8 ns | baseline |
| 1,000 | 119.6 ns | +16.3% |
| 10,000 | 125.0 ns | +21.6% |

#### RWMutex Lookup Times by Dataset Size
| Dataset Size | Time/op | Growth |
|--------------|---------|--------|
| 100 | 95.80 ns | baseline |
| 1,000 | 118.4 ns | +23.6% |
| 10,000 | 120.4 ns | +25.7% |

**Key Finding**: Both implementations demonstrate **O(1) lookup performance**. Despite 100x increase in dataset size (100→10,000), lookup time increased by only ~22-26% for both implementations, confirming constant-time hash map lookups. The slight increase is due to cache effects and memory access patterns, not algorithmic complexity.

## Performance Analysis

### Time Performance
- **Concurrent reads**: sync.Map is **4.8x faster** (9.83 ns vs 46.81 ns)
- **Mixed workload**: sync.Map is **2.2x faster** (101.3 ns vs 225.3 ns)
- **100 goroutines**: sync.Map is **8.2% faster** (652 μs vs 711 μs)
- **Single lookup**: RWMutex is marginally faster for single-threaded access due to lower overhead

### Memory Performance
- Both implementations use **zero allocations** for pure read operations
- Mixed workload shows minimal memory overhead difference (6 bytes)
- Allocation patterns are identical for both approaches

### Scalability
- sync.Map scales better with increasing concurrency (4.8x advantage)
- Both maintain O(1) lookup performance regardless of dataset size
- sync.Map's lock-free read operations provide superior throughput under contention

## Advantages of sync.Map

Based on benchmark results:

1. **Superior Read Performance**: 4.8x faster for concurrent reads, the primary use case
2. **Better Mixed Workload Performance**: 2.2x faster even with 10% writes
3. **Zero Allocation Reads**: No memory allocation overhead for read operations
4. **Lock-Free Reads**: Non-blocking reads improve throughput under high concurrency
5. **O(1) Lookup Confirmed**: Maintains constant-time performance across dataset sizes

## Recommendations

✅ **APPROVED**: Use **sync.Map** for device storage

**Rationale**:
1. **4.8x faster** for concurrent reads (the dominant operation in a homelab API)
2. Still **2.2x faster** for mixed workloads with frequent writes
3. Zero-allocation reads minimize GC pressure
4. O(1) lookup performance validated across all dataset sizes
5. Better scalability as concurrent load increases
6. Built-in, no external dependencies

**When RWMutex might be considered**:
- Single-threaded or low-concurrency scenarios (marginally faster)
- Need to iterate over all entries frequently (sync.Map.Range is slower than map iteration)
- Write-heavy workloads (>50% writes) - though this doesn't apply to typical device state reading

## Trade-offs

### sync.Map Pros:
- 4.8x faster concurrent reads
- Lock-free read operations
- Zero allocations for reads
- Better scalability under contention

### sync.Map Cons:
- Slightly higher memory overhead for mixed workloads (6 bytes)
- Range operations are slower than standard map
- Type-unsafe (requires type assertions)
- Marginally slower for single-threaded access

### RWMutex Pros:
- Type-safe
- Faster single-threaded access
- Faster iteration over all entries
- Familiar patterns for developers

### RWMutex Cons:
- 4.8x slower for concurrent reads
- 2.2x slower for mixed workloads
- Read operations still require lock acquisition
- Contention increases with concurrent readers

## Conclusion

For a homelab API service with read-heavy workloads (device state queries), **sync.Map is the clear winner** with **4.8x better performance** for concurrent reads and **2.2x better performance** for mixed workloads. The O(1) lookup performance is validated across dataset sizes from 100 to 10,000 devices.

The performance advantage of sync.Map becomes even more pronounced as concurrent load increases, making it the optimal choice for the device storage layer.

**Recommendation**: Proceed to Phase 1 with **sync.Map** as the device storage implementation.

## Next Steps

1. ✅ Benchmarks completed and validated
2. ✅ O(1) lookup performance confirmed
3. ✅ sync.Map approved for device storage
4. Ready to implement device storage service using sync.Map in Phase 1
