package research

import (
	"fmt"
	"sync"
	"testing"
)

// Device represents a simple device structure for benchmarking
type Device struct {
	ID     string
	Name   string
	Status string
	Value  float64
}

// SyncMapStorage implements device storage using sync.Map
type SyncMapStorage struct {
	data sync.Map
}

func NewSyncMapStorage() *SyncMapStorage {
	return &SyncMapStorage{}
}

func (s *SyncMapStorage) Store(id string, device Device) {
	s.data.Store(id, device)
}

func (s *SyncMapStorage) Load(id string) (Device, bool) {
	val, ok := s.data.Load(id)
	if !ok {
		return Device{}, false
	}
	return val.(Device), true
}

func (s *SyncMapStorage) Delete(id string) {
	s.data.Delete(id)
}

// RWMutexStorage implements device storage using RWMutex and map
type RWMutexStorage struct {
	mu   sync.RWMutex
	data map[string]Device
}

func NewRWMutexStorage() *RWMutexStorage {
	return &RWMutexStorage{
		data: make(map[string]Device),
	}
}

func (s *RWMutexStorage) Store(id string, device Device) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = device
}

func (s *RWMutexStorage) Load(id string) (Device, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, ok := s.data[id]
	return device, ok
}

func (s *RWMutexStorage) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
}

// Helper function to populate storage with test data
func populateStorage(count int) ([]Device, []string) {
	devices := make([]Device, count)
	ids := make([]string, count)

	for i := 0; i < count; i++ {
		id := fmt.Sprintf("device_%d", i)
		devices[i] = Device{
			ID:     id,
			Name:   fmt.Sprintf("Device %d", i),
			Status: "online",
			Value:  float64(i) * 1.5,
		}
		ids[i] = id
	}

	return devices, ids
}

// BenchmarkSyncMapConcurrentReads tests concurrent read performance with sync.Map
// RunParallel automatically uses GOMAXPROCS goroutines for realistic concurrent testing
func BenchmarkSyncMapConcurrentReads(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Cycle through device IDs
			id := ids[i%len(ids)]
			_, _ = storage.Load(id)
			i++
		}
	})
}

// BenchmarkRWMutexConcurrentReads tests concurrent read performance with RWMutex
// RunParallel automatically uses GOMAXPROCS goroutines for realistic concurrent testing
func BenchmarkRWMutexConcurrentReads(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Cycle through device IDs
			id := ids[i%len(ids)]
			_, _ = storage.Load(id)
			i++
		}
	})
}

// BenchmarkSyncMapConcurrentReads100Goroutines tests with exactly 100 goroutines
func BenchmarkSyncMapConcurrentReads100Goroutines(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	goroutines := 100
	readsPerGoroutine := b.N / goroutines

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < readsPerGoroutine; i++ {
				id := ids[(goroutineID*readsPerGoroutine+i)%len(ids)]
				_, _ = storage.Load(id)
			}
		}(g)
	}

	wg.Wait()
}

// BenchmarkRWMutexConcurrentReads100Goroutines tests with exactly 100 goroutines
func BenchmarkRWMutexConcurrentReads100Goroutines(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	goroutines := 100
	readsPerGoroutine := b.N / goroutines

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < readsPerGoroutine; i++ {
				id := ids[(goroutineID*readsPerGoroutine+i)%len(ids)]
				_, _ = storage.Load(id)
			}
		}(g)
	}

	wg.Wait()
}

// BenchmarkSyncMapMixedWorkload tests mixed read/write with sync.Map (90% reads, 10% writes)
func BenchmarkSyncMapMixedWorkload(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			id := ids[i%len(ids)]
			// 90% reads, 10% writes
			if i%10 == 0 {
				device := Device{
					ID:     id,
					Name:   fmt.Sprintf("Updated Device %d", i),
					Status: "online",
					Value:  float64(i) * 2.0,
				}
				storage.Store(id, device)
			} else {
				_, _ = storage.Load(id)
			}
			i++
		}
	})
}

// BenchmarkRWMutexMixedWorkload tests mixed read/write with RWMutex (90% reads, 10% writes)
func BenchmarkRWMutexMixedWorkload(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			id := ids[i%len(ids)]
			// 90% reads, 10% writes
			if i%10 == 0 {
				device := Device{
					ID:     id,
					Name:   fmt.Sprintf("Updated Device %d", i),
					Status: "online",
					Value:  float64(i) * 2.0,
				}
				storage.Store(id, device)
			} else {
				_, _ = storage.Load(id)
			}
			i++
		}
	})
}

// BenchmarkSyncMapWriteHeavy tests write-heavy workload with sync.Map (50% reads, 50% writes)
func BenchmarkSyncMapWriteHeavy(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			id := ids[i%len(ids)]
			// 50% reads, 50% writes
			if i%2 == 0 {
				device := Device{
					ID:     id,
					Name:   fmt.Sprintf("Updated Device %d", i),
					Status: "online",
					Value:  float64(i) * 2.0,
				}
				storage.Store(id, device)
			} else {
				_, _ = storage.Load(id)
			}
			i++
		}
	})
}

// BenchmarkRWMutexWriteHeavy tests write-heavy workload with RWMutex (50% reads, 50% writes)
func BenchmarkRWMutexWriteHeavy(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			id := ids[i%len(ids)]
			// 50% reads, 50% writes
			if i%2 == 0 {
				device := Device{
					ID:     id,
					Name:   fmt.Sprintf("Updated Device %d", i),
					Status: "online",
					Value:  float64(i) * 2.0,
				}
				storage.Store(id, device)
			} else {
				_, _ = storage.Load(id)
			}
			i++
		}
	})
}

// BenchmarkSyncMapSingleKeyContention tests all goroutines accessing same key
func BenchmarkSyncMapSingleKeyContention(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	hotKey := ids[0] // All goroutines will access this key

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = storage.Load(hotKey)
		}
	})
}

// BenchmarkRWMutexSingleKeyContention tests all goroutines accessing same key
func BenchmarkRWMutexSingleKeyContention(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(1000)

	// Populate storage
	for i, device := range devices {
		storage.Store(ids[i], device)
	}

	hotKey := ids[0] // All goroutines will access this key

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = storage.Load(hotKey)
		}
	})
}

// O(1) lookup validation benchmarks - test across different dataset sizes
// These validate that lookup time grows sub-linearly with dataset size

func BenchmarkSyncMap_SingleLookup_100(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(100)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}

func BenchmarkSyncMap_SingleLookup_1000(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(1000)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}

func BenchmarkSyncMap_SingleLookup_10000(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(10000)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}

func BenchmarkSyncMap_SingleLookup_100000(b *testing.B) {
	storage := NewSyncMapStorage()
	devices, ids := populateStorage(100000)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}

func BenchmarkRWMutex_SingleLookup_100(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(100)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}

func BenchmarkRWMutex_SingleLookup_1000(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(1000)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}

func BenchmarkRWMutex_SingleLookup_10000(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(10000)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}

func BenchmarkRWMutex_SingleLookup_100000(b *testing.B) {
	storage := NewRWMutexStorage()
	devices, ids := populateStorage(100000)
	for i, device := range devices {
		storage.Store(ids[i], device)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = storage.Load(ids[i%len(ids)])
	}
}
