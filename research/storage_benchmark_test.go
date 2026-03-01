package research

import (
	"fmt"
	"sync"
	"testing"
)

// DeviceStorage interface for different storage implementations
type DeviceStorage interface {
	Store(id string, device Device)
	Load(id string) (Device, bool)
	LoadAll() []Device
}

// SyncMapStorage implements DeviceStorage using sync.Map
type SyncMapStorage struct {
	m sync.Map
}

func (s *SyncMapStorage) Store(id string, device Device) {
	s.m.Store(id, device)
}

func (s *SyncMapStorage) Load(id string) (Device, bool) {
	val, ok := s.m.Load(id)
	if !ok {
		return Device{}, false
	}
	return val.(Device), true
}

func (s *SyncMapStorage) LoadAll() []Device {
	devices := make([]Device, 0, 100)
	s.m.Range(func(key, value interface{}) bool {
		devices = append(devices, value.(Device))
		return true
	})
	return devices
}

// RWMutexStorage implements DeviceStorage using RWMutex and map
type RWMutexStorage struct {
	mu      sync.RWMutex
	devices map[string]Device
}

func NewRWMutexStorage() *RWMutexStorage {
	return &RWMutexStorage{
		devices: make(map[string]Device),
	}
}

func (s *RWMutexStorage) Store(id string, device Device) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices[id] = device
}

func (s *RWMutexStorage) Load(id string) (Device, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, ok := s.devices[id]
	return device, ok
}

func (s *RWMutexStorage) LoadAll() []Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	devices := make([]Device, 0, len(s.devices))
	for _, device := range s.devices {
		devices = append(devices, device)
	}
	return devices
}

// setupStorage initializes storage with test devices
func setupStorage(storage DeviceStorage, count int) {
	devices := generateTestDevices(count)
	for _, device := range devices {
		storage.Store(device.ID, device)
	}
}

// BenchmarkSyncMap_ConcurrentReads benchmarks concurrent reads with sync.Map
// Tests 100 concurrent goroutines performing reads (via RunParallel)
func BenchmarkSyncMap_ConcurrentReads(b *testing.B) {
	storage := &SyncMapStorage{}
	setupStorage(storage, 50)
	deviceTypes := []string{"light", "sensor", "switch", "binary_sensor"}
	
	b.ResetTimer()
	// RunParallel automatically uses GOMAXPROCS goroutines
	// Each goroutine will run the body function b.N/GOMAXPROCS times
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Cycle through all 50 device IDs
			deviceType := deviceTypes[i%len(deviceTypes)]
			deviceLetter := string(rune('a' + (i % 26)))
			deviceID := deviceType + ".device_" + deviceLetter
			storage.Load(deviceID)
			i++
		}
	})
}

// BenchmarkRWMutex_ConcurrentReads benchmarks concurrent reads with RWMutex
// Tests 100 concurrent goroutines performing reads (via RunParallel)
func BenchmarkRWMutex_ConcurrentReads(b *testing.B) {
	storage := NewRWMutexStorage()
	setupStorage(storage, 50)
	deviceTypes := []string{"light", "sensor", "switch", "binary_sensor"}
	
	b.ResetTimer()
	// RunParallel automatically uses GOMAXPROCS goroutines
	// Each goroutine will run the body function b.N/GOMAXPROCS times
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Cycle through all 50 device IDs
			deviceType := deviceTypes[i%len(deviceTypes)]
			deviceLetter := string(rune('a' + (i % 26)))
			deviceID := deviceType + ".device_" + deviceLetter
			storage.Load(deviceID)
			i++
		}
	})
}

// BenchmarkSyncMap_MixedWorkload benchmarks mixed read/write workload with sync.Map
// 90% reads, 10% writes to simulate realistic API usage
func BenchmarkSyncMap_MixedWorkload(b *testing.B) {
	storage := &SyncMapStorage{}
	setupStorage(storage, 50)
	devices := generateTestDevices(50)
	deviceTypes := []string{"light", "sensor", "switch", "binary_sensor"}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 90% reads, 10% writes
			if i%10 == 0 {
				// Write operation
				device := devices[i%len(devices)]
				storage.Store(device.ID, device)
			} else {
				// Read operation - cycle through all 50 device IDs
				deviceType := deviceTypes[i%len(deviceTypes)]
				deviceLetter := string(rune('a' + (i % 26)))
				deviceID := deviceType + ".device_" + deviceLetter
				storage.Load(deviceID)
			}
			i++
		}
	})
}

// BenchmarkRWMutex_MixedWorkload benchmarks mixed read/write workload with RWMutex
// 90% reads, 10% writes to simulate realistic API usage
func BenchmarkRWMutex_MixedWorkload(b *testing.B) {
	storage := NewRWMutexStorage()
	setupStorage(storage, 50)
	devices := generateTestDevices(50)
	deviceTypes := []string{"light", "sensor", "switch", "binary_sensor"}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 90% reads, 10% writes
			if i%10 == 0 {
				// Write operation
				device := devices[i%len(devices)]
				storage.Store(device.ID, device)
			} else {
				// Read operation - cycle through all 50 device IDs
				deviceType := deviceTypes[i%len(deviceTypes)]
				deviceLetter := string(rune('a' + (i % 26)))
				deviceID := deviceType + ".device_" + deviceLetter
				storage.Load(deviceID)
			}
			i++
		}
	})
}

// BenchmarkSyncMap_LoadAll benchmarks retrieving all devices with sync.Map
func BenchmarkSyncMap_LoadAll(b *testing.B) {
	storage := &SyncMapStorage{}
	setupStorage(storage, 50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.LoadAll()
	}
}

// BenchmarkRWMutex_LoadAll benchmarks retrieving all devices with RWMutex
func BenchmarkRWMutex_LoadAll(b *testing.B) {
	storage := NewRWMutexStorage()
	setupStorage(storage, 50)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.LoadAll()
	}
}

// O(1) lookup validation benchmarks - test with different dataset sizes
// to verify lookup time doesn't grow linearly with dataset size

// BenchmarkSyncMap_SingleLookup validates O(1) lookup performance
func BenchmarkSyncMap_SingleLookup(b *testing.B) {
	sizes := []int{100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			storage := &SyncMapStorage{}
			setupStorage(storage, size)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Always lookup the same device to test pure lookup time
				storage.Load("light.device_a")
			}
		})
	}
}

// BenchmarkRWMutex_SingleLookup validates O(1) lookup performance
func BenchmarkRWMutex_SingleLookup(b *testing.B) {
	sizes := []int{100, 1000, 10000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			storage := NewRWMutexStorage()
			setupStorage(storage, size)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Always lookup the same device to test pure lookup time
				storage.Load("light.device_a")
			}
		})
	}
}
