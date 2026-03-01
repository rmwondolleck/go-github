package research

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Device represents a smart home device for benchmarking
type Device struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	State       string            `json:"state"`
	Attributes  map[string]string `json:"attributes"`
	LastUpdated time.Time         `json:"last_updated"`
}

// generateDevices creates a slice of n devices for benchmarking
func generateDevices(n int) []Device {
	devices := make([]Device, n)
	for i := 0; i < n; i++ {
		devices[i] = Device{
			ID:    "device_" + strconv.Itoa(i),
			Name:  "Device " + strconv.Itoa(i),
			Type:  "sensor",
			State: "on",
			Attributes: map[string]string{
				"temperature": "22.5",
				"humidity":    "45",
				"battery":     "85",
			},
			LastUpdated: time.Now(),
		}
	}
	return devices
}

// Benchmark stdlib encoding/json with 50 devices
func BenchmarkStdlib_50Devices(b *testing.B) {
	devices := generateDevices(50)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(devices)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark jsoniter with ConfigCompatibleWithStandardLibrary (50 devices)
func BenchmarkJsoniter_50Devices_Compatible(b *testing.B) {
	devices := generateDevices(50)
	var jsonCompat = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := jsonCompat.Marshal(devices)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark jsoniter with ConfigFastest (50 devices)
func BenchmarkJsoniter_50Devices_Fastest(b *testing.B) {
	devices := generateDevices(50)
	var jsonFastest = jsoniter.ConfigFastest
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := jsonFastest.Marshal(devices)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark jsoniter ConfigFastest with Stream API (50 devices)
func BenchmarkJsoniter_50Devices_FastestStream(b *testing.B) {
	devices := generateDevices(50)
	var jsonFastest = jsoniter.ConfigFastest
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stream := jsonFastest.BorrowStream(nil)
		stream.WriteVal(devices)
		if stream.Error != nil {
			b.Fatal(stream.Error)
		}
		_ = stream.Buffer()
		jsonFastest.ReturnStream(stream)
	}
}
