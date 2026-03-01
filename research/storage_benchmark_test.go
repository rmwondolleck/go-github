package research

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Device represents a smart home device in the system
type Device struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	State       string            `json:"state"`
	Attributes  map[string]string `json:"attributes"`
	LastUpdated time.Time         `json:"last_updated"`
}

// SyncMapStorage implements device storage using sync.Map
type SyncMapStorage struct {
	devices sync.Map
}

func NewSyncMapStorage() *SyncMapStorage {
	return &SyncMapStorage{}
}

func (s *SyncMapStorage) Store(id string, device *Device) {
	s.devices.Store(id, device)
}

func (s *SyncMapStorage) Load(id string) (*Device, bool) {
	val, ok := s.devices.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*Device), true
}

// RWMutexStorage implements device storage using RWMutex
type RWMutexStorage struct {
	mu      sync.RWMutex
	devices map[string]*Device
}

func NewRWMutexStorage() *RWMutexStorage {
	return &RWMutexStorage{
		devices: make(map[string]*Device),
	}
}

func (s *RWMutexStorage) Store(id string, device *Device) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices[id] = device
}

func (s *RWMutexStorage) Load(id string) (*Device, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, ok := s.devices[id]
	return device, ok
}

// Helper function to create a test device
func createTestDevice(id string) *Device {
	return &Device{
		ID:          id,
		Name:        fmt.Sprintf("Device %s", id),
		Type:        "smart_light",
		State:       "on",
		Attributes:  map[string]string{"brightness": "80", "color": "warm_white"},
		LastUpdated: time.Now(),
	}
}

// BenchmarkSyncMap_ConcurrentReads benchmarks concurrent reads with sync.Map (100 goroutines)
func BenchmarkSyncMap_ConcurrentReads(b *testing.B) {
	storage := NewSyncMapStorage()
	
	// Pre-populate with 1000 devices
	for i := 0; i < 1000; i++ {
		deviceID := fmt.Sprintf("device-%d", i)
		storage.Store(deviceID, createTestDevice(deviceID))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		deviceID := "device-500" // Read from middle of dataset
		for pb.Next() {
			_, _ = storage.Load(deviceID)
		}
	})
}

// BenchmarkRWMutex_ConcurrentReads benchmarks concurrent reads with RWMutex (100 goroutines)
func BenchmarkRWMutex_ConcurrentReads(b *testing.B) {
	storage := NewRWMutexStorage()
	
	// Pre-populate with 1000 devices
	for i := 0; i < 1000; i++ {
		deviceID := fmt.Sprintf("device-%d", i)
		storage.Store(deviceID, createTestDevice(deviceID))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		deviceID := "device-500" // Read from middle of dataset
		for pb.Next() {
			_, _ = storage.Load(deviceID)
		}
	})
}

// BenchmarkSyncMap_MixedWorkload benchmarks 90% reads, 10% writes with sync.Map
func BenchmarkSyncMap_MixedWorkload(b *testing.B) {
	storage := NewSyncMapStorage()
	
	// Pre-populate with 1000 devices
	for i := 0; i < 1000; i++ {
		deviceID := fmt.Sprintf("device-%d", i)
		storage.Store(deviceID, createTestDevice(deviceID))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			counter++
			if counter%10 == 0 {
				// 10% writes
				deviceID := fmt.Sprintf("device-%d", counter%1000)
				storage.Store(deviceID, createTestDevice(deviceID))
			} else {
				// 90% reads
				deviceID := fmt.Sprintf("device-%d", counter%1000)
				_, _ = storage.Load(deviceID)
			}
		}
	})
}

// BenchmarkRWMutex_MixedWorkload benchmarks 90% reads, 10% writes with RWMutex
func BenchmarkRWMutex_MixedWorkload(b *testing.B) {
	storage := NewRWMutexStorage()
	
	// Pre-populate with 1000 devices
	for i := 0; i < 1000; i++ {
		deviceID := fmt.Sprintf("device-%d", i)
		storage.Store(deviceID, createTestDevice(deviceID))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			counter++
			if counter%10 == 0 {
				// 10% writes
				deviceID := fmt.Sprintf("device-%d", counter%1000)
				storage.Store(deviceID, createTestDevice(deviceID))
			} else {
				// 90% reads
				deviceID := fmt.Sprintf("device-%d", counter%1000)
				_, _ = storage.Load(deviceID)
			}
		}
	})
}

// BenchmarkSyncMap_100Goroutines_ConcurrentReads explicitly tests with 100 goroutines
func BenchmarkSyncMap_100Goroutines_ConcurrentReads(b *testing.B) {
	storage := NewSyncMapStorage()
	
	// Pre-populate with 1000 devices
	for i := 0; i < 1000; i++ {
		deviceID := fmt.Sprintf("device-%d", i)
		storage.Store(deviceID, createTestDevice(deviceID))
	}

	b.ResetTimer()
	
	// Run N iterations
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(100)
		
		// Launch 100 goroutines
		for i := 0; i < 100; i++ {
			go func(goroutineID int) {
				defer wg.Done()
				// Each goroutine performs 100 reads
				for j := 0; j < 100; j++ {
					deviceID := fmt.Sprintf("device-%d", (goroutineID*10+j)%1000)
					_, _ = storage.Load(deviceID)
				}
			}(i)
		}
		
		wg.Wait()
	}
}

// BenchmarkRWMutex_100Goroutines_ConcurrentReads explicitly tests with 100 goroutines
func BenchmarkRWMutex_100Goroutines_ConcurrentReads(b *testing.B) {
	storage := NewRWMutexStorage()
	
	// Pre-populate with 1000 devices
	for i := 0; i < 1000; i++ {
		deviceID := fmt.Sprintf("device-%d", i)
		storage.Store(deviceID, createTestDevice(deviceID))
	}

	b.ResetTimer()
	
	// Run N iterations
	for n := 0; n < b.N; n++ {
		var wg sync.WaitGroup
		wg.Add(100)
		
		// Launch 100 goroutines
		for i := 0; i < 100; i++ {
			go func(goroutineID int) {
				defer wg.Done()
				// Each goroutine performs 100 reads
				for j := 0; j < 100; j++ {
					deviceID := fmt.Sprintf("device-%d", (goroutineID*10+j)%1000)
					_, _ = storage.Load(deviceID)
				}
			}(i)
		}
		
		wg.Wait()
	}
}

// BenchmarkSyncMap_SingleLookup benchmarks O(1) lookup performance for sync.Map
func BenchmarkSyncMap_SingleLookup(b *testing.B) {
	storage := NewSyncMapStorage()
	
	// Pre-populate with various dataset sizes to verify O(1) behavior
	sizes := []int{100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Clear and repopulate
			storage = NewSyncMapStorage()
			for i := 0; i < size; i++ {
				deviceID := fmt.Sprintf("device-%d", i)
				storage.Store(deviceID, createTestDevice(deviceID))
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				deviceID := fmt.Sprintf("device-%d", size/2) // Lookup from middle
				_, _ = storage.Load(deviceID)
			}
		})
	}
}

// BenchmarkRWMutex_SingleLookup benchmarks O(1) lookup performance for RWMutex
func BenchmarkRWMutex_SingleLookup(b *testing.B) {
	storage := NewRWMutexStorage()
	
	// Pre-populate with various dataset sizes to verify O(1) behavior
	sizes := []int{100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Clear and repopulate
			storage = NewRWMutexStorage()
			for i := 0; i < size; i++ {
				deviceID := fmt.Sprintf("device-%d", i)
				storage.Store(deviceID, createTestDevice(deviceID))
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				deviceID := fmt.Sprintf("device-%d", size/2) // Lookup from middle
				_, _ = storage.Load(deviceID)
			}
		})
	}
}
