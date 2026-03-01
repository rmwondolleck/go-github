package research

import (
	"sync"
	"testing"
)

// DeviceStore using RWMutex
type RWMutexStore struct {
	mu      sync.RWMutex
	devices map[string]*Device
}

func NewRWMutexStore() *RWMutexStore {
	return &RWMutexStore{
		devices: make(map[string]*Device),
	}
}

func (s *RWMutexStore) Get(id string) (*Device, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, ok := s.devices[id]
	return device, ok
}

func (s *RWMutexStore) Set(id string, device *Device) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices[id] = device
}

// DeviceStore using sync.Map
type SyncMapStore struct {
	devices sync.Map
}

func NewSyncMapStore() *SyncMapStore {
	return &SyncMapStore{}
}

func (s *SyncMapStore) Get(id string) (*Device, bool) {
	val, ok := s.devices.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*Device), true
}

func (s *SyncMapStore) Set(id string, device *Device) {
	s.devices.Store(id, device)
}

// Setup function to populate stores with test data
func populateStore(rwStore *RWMutexStore, smStore *SyncMapStore, count int) {
	devices := generateDevices(count)
	for i, device := range devices {
		id := device.ID
		if rwStore != nil {
			rwStore.Set(id, &devices[i])
		}
		if smStore != nil {
			smStore.Set(id, &devices[i])
		}
	}
}

// Benchmark concurrent reads with RWMutex (100 goroutines, read-heavy workload)
func BenchmarkRWMutex_ConcurrentReads(b *testing.B) {
	store := NewRWMutexStore()
	populateStore(store, nil, 50)
	
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = store.Get("device_0")
		}
	})
}

// Benchmark concurrent reads with sync.Map (100 goroutines, read-heavy workload)
func BenchmarkSyncMap_ConcurrentReads(b *testing.B) {
	store := NewSyncMapStore()
	populateStore(nil, store, 50)
	
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = store.Get("device_0")
		}
	})
}

// Benchmark mixed read/write workload with RWMutex (90% reads, 10% writes)
func BenchmarkRWMutex_MixedWorkload(b *testing.B) {
	store := NewRWMutexStore()
	populateStore(store, nil, 50)
	device := &Device{
		ID:    "device_new",
		Name:  "New Device",
		Type:  "sensor",
		State: "on",
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				// 10% writes
				store.Set("device_0", device)
			} else {
				// 90% reads
				_, _ = store.Get("device_0")
			}
			i++
		}
	})
}

// Benchmark mixed read/write workload with sync.Map (90% reads, 10% writes)
func BenchmarkSyncMap_MixedWorkload(b *testing.B) {
	store := NewSyncMapStore()
	populateStore(nil, store, 50)
	device := &Device{
		ID:    "device_new",
		Name:  "New Device",
		Type:  "sensor",
		State: "on",
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 == 0 {
				// 10% writes
				store.Set("device_0", device)
			} else {
				// 90% reads
				_, _ = store.Get("device_0")
			}
			i++
		}
	})
}
