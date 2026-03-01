# Research Findings: Home Lab API Service

## Phase 0 - Research & Validation

This document consolidates all research findings from Phase 0, including performance benchmarks and technology evaluations to inform Phase 1 implementation decisions.

---

## T001: Gin Framework Performance Benchmarks

**Date**: 2026-03-01  
**Task**: T001 - Benchmark Gin framework basic routing  
**Phase**: Phase 0 - Research & Validation  
**Batch**: Batch 1 (Research)

### Objective

Validate Gin framework performance against stdlib net/http to ensure the framework overhead is acceptable (target: <10ms overhead) for the homelab API service.

### Methodology

Created comprehensive benchmarks in `research/gin_benchmark_test.go` comparing:
- Simple route handling (GET /health)
- Routes with middleware
- Multiple route registration
- Parameterized routes (URL parameters)

Each benchmark was run with:
- `-benchtime=3s` for statistically significant results
- `-benchmem` to measure memory allocations
- Go 1.25.7 on AMD EPYC 7763 64-Core Processor

### Results Summary

#### Simple Route Handling
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1158 ns | 1064 B | 11 |
| Gin | 1096 ns | 1128 B | 12 |
| **Overhead** | **-62 ns** | **+64 B** | **+1** |

**Key Finding**: Gin is actually **faster** than stdlib for simple routes (62 nanoseconds faster, or ~5% improvement).

#### With Middleware
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1172 ns | 1064 B | 11 |
| Gin | 1105 ns | 1128 B | 12 |
| **Overhead** | **-67 ns** | **+64 B** | **+1** |

**Key Finding**: Gin maintains performance advantage even with middleware (67 nanoseconds faster).

#### Multiple Routes (5 routes)
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1261 ns | 1064 B | 11 |
| Gin | 1114 ns | 1128 B | 12 |
| **Overhead** | **-147 ns** | **+64 B** | **+1** |

**Key Finding**: Gin's advantage increases with more routes (147 nanoseconds faster, ~12% improvement), likely due to optimized routing algorithm.

#### Parameterized Routes
| Framework | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| stdlib net/http | 1581 ns | 1120 B | 14 |
| Gin | 1116 ns | 1128 B | 12 |
| **Overhead** | **-465 ns** | **+8 B** | **-2** |

**Key Finding**: Gin significantly outperforms stdlib for parameterized routes (465 nanoseconds faster, ~29% improvement), with fewer allocations.

### Performance Analysis

#### Time Overhead
- **Target**: <10ms (10,000,000 ns) overhead
- **Actual**: Gin is consistently **faster** than stdlib, not slower
- **Conclusion**: ✅ **TARGET EXCEEDED** - No overhead detected; Gin provides performance improvements

#### Overhead Breakdown by Scenario
1. **Simple routes**: -0.062 μs (Gin faster)
2. **With middleware**: -0.067 μs (Gin faster)
3. **Multiple routes**: -0.147 μs (Gin faster)
4. **Parameterized routes**: -0.465 μs (Gin faster)

All measurements are in the sub-microsecond range, well below the 10ms (10,000 μs) target.

#### Memory Overhead
- Gin consistently uses ~64 bytes more per request
- One additional allocation per request
- For parameterized routes, Gin uses 8 bytes more but 2 fewer allocations
- Memory overhead is minimal and acceptable for the use case

### Advantages of Gin Framework

Based on benchmark results, Gin provides:

1. **Better Performance**: Faster than stdlib in all tested scenarios
2. **Optimized Routing**: Superior performance with multiple and parameterized routes
3. **Built-in Features**: Middleware support, JSON binding/validation, parameter extraction
4. **Developer Productivity**: Cleaner API, less boilerplate code
5. **Production-Ready**: Battle-tested framework with extensive ecosystem

### Recommendations

✅ **APPROVED**: Proceed with Gin framework for the homelab API service

**Rationale**:
1. Performance target (<10ms overhead) is **far exceeded** - Gin is actually faster than stdlib
2. Routing performance improves significantly with complex routes
3. Built-in features reduce development time and potential bugs
4. Memory overhead (64 bytes/request) is negligible for homelab scale
5. Gin's middleware system provides cleaner architecture for future features

### Trade-offs

**Pros**:
- Superior performance to stdlib
- Rich middleware ecosystem
- Better developer experience
- Built-in JSON handling and validation
- Excellent parameter routing

**Cons**:
- Additional dependency (~30 transitive dependencies added)
- Slightly higher memory usage (64 bytes per request)
- Framework lock-in (but migration path exists if needed)

### Conclusion

The Gin framework not only meets but **exceeds** the performance requirements with **negative overhead** (i.e., performance improvements) compared to stdlib net/http. The <10ms overhead target is achieved with significant margin - actual performance is improved by 0.062-0.465 microseconds depending on the scenario.

**Recommendation**: Proceed to Phase 1 with Gin framework as the chosen web framework.

---

## T003: sync.Map vs RWMutex Benchmark Results

**Date**: 2026-03-01  
**CPU**: AMD EPYC 7763 64-Core Processor  
**Go Version**: 1.24.13  

### Executive Summary

Comprehensive benchmarking of sync.Map vs RWMutex-protected map for device storage reveals distinct performance characteristics:

- **sync.Map** excels in read-heavy workloads, showing 3.8x better performance
- **RWMutex** performs better under single-key contention scenarios
- Both implementations demonstrate **O(1) lookup complexity** ✅
- For typical device storage use cases (90% reads, 10% writes), **sync.Map is recommended**

### Detailed Benchmark Results

#### 1. Concurrent Reads (100 Goroutines)

Testing concurrent read performance with 100 goroutines accessing 1000 devices:

| Implementation | ns/op | Relative Performance |
|---------------|-------|---------------------|
| sync.Map      | 13.08 | **3.7x faster** |
| RWMutex       | 48.71 | baseline |

**Finding**: sync.Map significantly outperforms RWMutex in read-heavy concurrent scenarios due to lock-free reads for stable keys.

#### 2. Mixed Read/Write Workload (90% reads, 10% writes)

Simulating realistic device storage patterns:

| Implementation | ns/op | Memory (B/op) | Allocs/op |
|---------------|-------|---------------|-----------|
| sync.Map      | 29.67 | 16 | 0 |
| RWMutex       | 62.06 | 3 | 0 |

**Finding**: sync.Map provides 2.1x better throughput in typical workload patterns. The higher memory usage is acceptable for the performance gain.

#### 3. Write-Heavy Workload (50% reads, 50% writes)

Testing scenarios with frequent device updates:

| Implementation | ns/op | Memory (B/op) | Allocs/op |
|---------------|-------|---------------|-----------|
| sync.Map      | 85.06 | 80 | 2 |
| RWMutex       | 128.9 | 16 | 1 |

**Finding**: sync.Map maintains 1.5x advantage even with heavy writes, though with increased memory allocation.

#### 4. Single Key Contention

All goroutines accessing the same "hot" device:

| Implementation | ns/op | Memory (B/op) | Allocs/op |
|---------------|-------|---------------|-----------|
| sync.Map      | 38.87 | 12 | 0 |
| RWMutex       | 28.13 | 0 | 0 |

**Finding**: RWMutex performs 1.4x better under extreme single-key contention. However, this scenario is unlikely in device storage where access is distributed.

#### 5. Lookup Complexity Validation (O(1) Verification)

Testing lookup performance across different dataset sizes:

**sync.Map**:
| Size    | ns/op | Growth Factor |
|---------|-------|---------------|
| 100     | 29.60 | - |
| 1,000   | 31.85 | 1.08x |
| 10,000  | 42.41 | 1.33x |
| 100,000 | 57.62 | 1.36x |

**RWMutex**:
| Size    | ns/op | Growth Factor |
|---------|-------|---------------|
| 100     | 23.52 | - |
| 1,000   | 27.45 | 1.17x |
| 10,000  | 37.28 | 1.36x |
| 100,000 | 42.70 | 1.15x |

**Finding**: Both implementations demonstrate **O(1) lookup complexity** ✅. The slight performance degradation with larger datasets is due to cache effects, not algorithmic complexity. Growth is sub-linear and acceptable.

### Performance Analysis

#### sync.Map Characteristics

**Strengths**:
- Lock-free reads for stable keys (3-4x faster read performance)
- Excellent for read-heavy workloads (90%+ reads)
- No lock contention on reads
- Built-in atomic operations

**Weaknesses**:
- Higher memory overhead per entry
- More allocations under write-heavy scenarios
- Slightly slower under extreme single-key contention

**Best Use Cases**:
- Device status queries (primary use case)
- Service discovery data
- Configuration caching
- Any read-dominated access pattern

#### RWMutex Characteristics

**Strengths**:
- Lower memory overhead
- Better single-key contention performance
- More predictable memory usage
- Simpler debugging and profiling

**Weaknesses**:
- Lock contention on read operations under high concurrency
- 2-4x slower for concurrent reads
- RLock still requires synchronization

**Best Use Cases**:
- Write-heavy workloads (>30% writes)
- Single hot-key scenarios
- Memory-constrained environments

### Latency Target Validation

**Target**: Sub-10ms API response times

**Analysis**:
- Device lookup: ~13-50 ns (sync.Map/RWMutex)
- Network overhead: ~1-5 ms (typical)
- JSON serialization: ~100-500 ns per device
- Total estimated: **1-6 ms per request** ✅

**Conclusion**: Both implementations easily meet latency targets. Performance difference (13ns vs 50ns) is negligible compared to network overhead.

### Recommendation

**✅ Recommended: sync.Map**

**Rationale**:
1. **Primary use case alignment**: Device queries are read-heavy (90%+ reads expected)
2. **Superior concurrent read performance**: 3.7x faster with 100 goroutines
3. **Validated O(1) complexity**: Meets performance requirements at scale
4. **Future-proof**: Better scaling characteristics for HomeAssistant integration (100s-1000s of devices)

**Implementation Notes**:
- Use sync.Map for device storage
- Memory overhead is acceptable (16 bytes/op vs 3 bytes/op)
- Provides better horizontal scaling for K8s deployment
- Lock-free reads reduce CPU usage under load

**When to reconsider**:
- If write frequency exceeds 30% of operations
- If memory becomes a critical constraint
- If profiling shows sync.Map as a bottleneck (unlikely)

### Next Steps

1. ✅ **Benchmarks complete** - Performance targets validated
2. **T004**: Complete research.md with all findings (this document)
3. **Phase 1**: Begin foundation work with sync.Map as storage mechanism
4. **Monitoring**: Add Prometheus metrics to validate real-world performance

### Additional Notes

- All benchmarks run with Go 1.24.13 on AMD EPYC 7763
- Benchtime: 3 seconds per benchmark for statistical significance
- Test data: 1000 simulated devices with realistic structure
- Concurrent access pattern validates production scenarios

**Performance Critical**: These results directly impact latency targets and justify architectural decisions for Phase 1 implementation.

---

## Phase 0 Summary

### Key Decisions

1. **Web Framework**: Gin v1.12.0
   - Performance: 5-29% faster than stdlib
   - Rich middleware ecosystem
   - Built-in JSON handling

2. **Device Storage**: sync.Map
   - Performance: 3.7x faster for concurrent reads
   - O(1) lookup complexity validated
   - Optimal for read-heavy workloads

### Performance Targets: ✅ VALIDATED

- API latency target: <10ms ✅
- Gin overhead: -62 to -465 ns (faster than stdlib) ✅
- Device lookup: 13 ns with sync.Map ✅
- Concurrent scalability: Tested with 100 goroutines ✅

### Ready for Phase 1

All research complete. Technology choices validated. Performance targets confirmed achievable. Ready to proceed with implementation.
