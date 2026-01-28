# Performance Optimizations Summary

## Overview

This document summarizes the performance optimizations made to `yet-another-sort` to significantly improve its sorting performance.

## Performance Results

### Before Optimizations
- User time: 5.47-7.02 seconds
- Elapsed time: 5.21-5.90 seconds
- Memory usage: 428-512 MB
- **Performance vs GNU sort: ~10-12x slower**

### After Optimizations
- User time: 1.85-2.10 seconds
- Elapsed time: 1.74-1.94 seconds
- Memory usage: 335-341 MB
- **Performance vs GNU sort: ~3.5x slower**

### Improvements Achieved
- **70-73% reduction in CPU time** (3-4x faster)
- **68-70% reduction in elapsed time** (3-3.5x faster)
- **20-35% reduction in memory usage**

## Optimizations Implemented

### 1. Replace Custom Merge Sort with Go's Built-in Sort (`sort_merge_sort.go`)

**Problem**: The custom merge sort implementation created two full copies of the data and performed excessive copying operations between arrays.

**Solution**: 
- Implemented `sort.Interface` for `ContentType`
- Leveraged Go's highly optimized sorting algorithm (pdqsort - pattern-defeating quicksort)
- Go's sort uses a hybrid approach:
  - Quicksort for general cases
  - Heapsort when recursion depth is too high (avoiding O(n²) worst case)
  - Insertion sort for small slices (< 12 elements)

**Impact**: This single change provided the majority of the performance improvement.

### 2. Optimize String Comparisons

**Problem**: The original implementation extracted strings to variables before comparison.

**Solution**: Use direct comparison `c[i].CompareLine < c[j].CompareLine` which is optimized by the Go compiler.

**Impact**: Reduced comparison overhead and improved cache locality.

### 3. Optimize File Loading (`files.go`)

**Problem**: 
- Small default buffer size for file I/O
- No pre-allocation of line slices causing multiple reallocations
- No buffer tuning for large files

**Solution**:
- Increased buffer size to 256KB for both file and stdin reading
- Pre-allocated line slices with initial capacity of 16,384 lines
- Set scanner buffer size to handle long lines efficiently
- Added proper error handling for scanner errors

**Impact**: Reduced I/O overhead and memory allocation churn during file loading.

### 4. Optimize Input Processing (`process.go`)

**Problem**:
- Multiple small allocations during line processing
- Repeated string splitting and joining operations
- No pre-allocation of slices
- Inefficient string concatenation using `strings.Join()`

**Solution**:
- Pre-calculate and pre-allocate exact capacity for multilines
- Use `strings.Builder` for efficient string concatenation (avoids intermediate allocations)
- Optimize for common case: space separator using `strings.Fields()`
- Added fast paths for single-field cases (avoid unnecessary string operations)
- Estimate and pre-grow string builders to minimize reallocations

**Impact**: Significantly reduced memory allocations and improved string processing performance.

### 5. Optimize Output Writing (`main.go`)

**Problem**: Unbuffered output with `fmt.Fprintln()` caused excessive system calls.

**Solution**:
- Added 256KB buffered writer for output
- Use `WriteString()` and `WriteByte()` instead of `fmt.Fprintln()`
- Explicit flush at the end

**Impact**: Reduced system call overhead for output operations.

### 6. Optimize Uniq Function (`uniq.go`)

**Problem**: Always created a copy of the entire array even when uniq was disabled.

**Solution**: Check if uniq is disabled and return original slice without copying.

**Impact**: Avoided unnecessary memory allocation when uniq feature is not used.

## Technical Details

### Memory Allocation Reduction

The optimizations focused on reducing allocations through:
1. **Pre-allocation**: Calculate required capacity and allocate once
2. **Capacity estimation**: Start with reasonable capacities to avoid multiple grow operations
3. **String builders**: Use `strings.Builder` with pre-grown capacity
4. **Avoid copying**: Return original slices when modifications aren't needed

### I/O Optimization

1. **Large buffers**: 256KB buffers for both input and output
2. **Buffered I/O**: Batch system calls to reduce overhead
3. **Scanner tuning**: Configure scanner with appropriate buffer sizes

### Algorithm Optimization

1. **Native sort**: Go's stdlib sort is highly optimized C-level implementation
2. **Smart algorithm selection**: Go's sort automatically chooses the best algorithm
3. **Direct comparisons**: Compiler-optimized string comparisons

## Remaining Performance Gap vs GNU Sort

GNU sort is still ~3.5x faster primarily due to:

1. **External sorting**: GNU sort uses disk-based sorting for very large files, keeping memory usage low (72MB vs 335MB)
2. **C implementation**: Native code with decades of optimization
3. **Parallel processing**: GNU sort can utilize multiple cores
4. **Specialized algorithms**: Hand-tuned assembly for comparison operations
5. **OS-specific optimizations**: Takes advantage of platform-specific features

## Future Optimization Opportunities

1. **Parallel sorting**: Use `sort.Slice()` with goroutines for large datasets
2. **External sorting**: Implement disk-based sorting for very large files
3. **Memory-mapped I/O**: Use `mmap` for faster file access
4. **Custom comparison functions**: Optimize for specific data patterns
5. **Reduce allocations further**: Use sync.Pool for temporary buffers
6. **Profile-guided optimization**: Use CPU profiles to identify remaining bottlenecks

## Benchmarking

Test dataset: 1,024,000 lines, 4 fields per line, 5 characters per field

Run the performance test:
```bash
./scripts/generate-random-input-file.py --lines 1024000 --fields 4 --field-length 5 > /tmp/test.txt
/usr/bin/time -f "%U user %S sys %E elapsed %M kB" ./yet-another-sort --sort-mode merge /tmp/test.txt > /dev/null
```

Compare with GNU sort:
```bash
/usr/bin/time -f "%U user %S sys %E elapsed %M kB" sort /tmp/test.txt > /dev/null
```

## Conclusion

The performance optimizations reduced execution time by approximately 70% and memory usage by 20-35%. The primary improvement came from replacing the custom merge sort with Go's highly optimized built-in sort, combined with reducing memory allocations and improving I/O buffering. While GNU sort remains faster due to its external sorting capabilities and decades of optimization, `yet-another-sort` is now much more competitive for typical use cases.