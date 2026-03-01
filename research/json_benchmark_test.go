package research

import (
	"encoding/json"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// Device represents a HomeAssistant device with all its properties
type Device struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastUpdated time.Time              `json:"last_updated"`
}

// generateTestDevices creates a slice of N test devices with realistic data
func generateTestDevices(count int) []Device {
	devices := make([]Device, count)
	deviceTypes := []string{"light", "sensor", "switch", "binary_sensor"}
	states := []string{"on", "off", "available", "unavailable"}

	for i := 0; i < count; i++ {
		deviceType := deviceTypes[i%len(deviceTypes)]
		devices[i] = Device{
			ID:    deviceType + "." + "device_" + string(rune('a'+i%26)),
			Name:  "Test Device " + string(rune('A'+i%26)),
			Type:  deviceType,
			State: states[i%len(states)],
			Attributes: map[string]interface{}{
				"friendly_name": "Test Device " + string(rune('A'+i%26)),
				"brightness":    128,
				"temperature":   23.5,
				"humidity":      65,
				"battery":       95,
			},
			LastUpdated: time.Now(),
		}
	}

	return devices
}

// BenchmarkStdlib_50Devices benchmarks stdlib encoding/json for encoding 50 devices
func BenchmarkStdlib_50Devices(b *testing.B) {
	devices := generateTestDevices(50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(devices)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJsoniter_50Devices_Fastest benchmarks jsoniter with ConfigFastest for encoding 50 devices
func BenchmarkJsoniter_50Devices_Fastest(b *testing.B) {
	devices := generateTestDevices(50)
	jsonAPI := jsoniter.ConfigFastest

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jsonAPI.Marshal(devices)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJsoniter_50Devices_Compatible benchmarks jsoniter with ConfigCompatibleWithStandardLibrary
func BenchmarkJsoniter_50Devices_Compatible(b *testing.B) {
	devices := generateTestDevices(50)
	jsonAPI := jsoniter.ConfigCompatibleWithStandardLibrary

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jsonAPI.Marshal(devices)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStdlib_50Devices_Stream benchmarks stdlib encoding/json using Encoder (stream API)
func BenchmarkStdlib_50Devices_Stream(b *testing.B) {
	devices := generateTestDevices(50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate writing to a buffer (like an HTTP response)
		buf := make([]byte, 0, 8192)
		encoder := json.NewEncoder(&bytesBuffer{buf: buf})
		if err := encoder.Encode(devices); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJsoniter_50Devices_Stream benchmarks jsoniter using Stream API
func BenchmarkJsoniter_50Devices_Stream(b *testing.B) {
	devices := generateTestDevices(50)
	jsonAPI := jsoniter.ConfigFastest

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, 8192)
		stream := jsonAPI.BorrowStream(&bytesBuffer{buf: buf})
		stream.WriteVal(devices)
		if stream.Error != nil {
			b.Fatal(stream.Error)
		}
		jsonAPI.ReturnStream(stream)
	}
}

// bytesBuffer is a simple writer that appends to a byte slice
type bytesBuffer struct {
	buf []byte
}

func (b *bytesBuffer) Write(p []byte) (n int, err error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}
