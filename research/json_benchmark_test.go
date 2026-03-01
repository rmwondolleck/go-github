package research

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

// Device represents a HomeAssistant device for benchmarking
type Device struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastUpdated string                 `json:"last_updated"`
}

// DeviceListResponse represents the response with 50 devices
type DeviceListResponse struct {
	Devices   []Device `json:"devices"`
	Total     int      `json:"total"`
	RequestID string   `json:"request_id"`
}

// generateMockDevices creates a slice of 50 mock devices
func generateMockDevices() []Device {
	devices := make([]Device, 50)
	for i := 0; i < 50; i++ {
		devices[i] = Device{
			ID:    "device." + string(rune('a'+i%26)) + string(rune('0'+(i/26))),
			Name:  "Test Device " + string(rune('0'+i/10)) + string(rune('0'+i%10)),
			Type:  []string{"light", "sensor", "switch", "binary_sensor"}[i%4],
			State: []string{"on", "off", "active", "inactive"}[i%4],
			Attributes: map[string]interface{}{
				"brightness":  i * 5,
				"temperature": 20.0 + float64(i),
				"color_temp":  370,
				"room":        "Room " + string(rune('A'+i%10)),
			},
			LastUpdated: "2026-02-28T10:30:00Z",
		}
	}
	return devices
}

// BenchmarkStdlibJSON_Single tests encoding a single device with stdlib
func BenchmarkStdlibJSON_Single(b *testing.B) {
	device := Device{
		ID:    "light.living_room",
		Name:  "Living Room Light",
		Type:  "light",
		State: "on",
		Attributes: map[string]interface{}{
			"brightness": 200,
			"color_temp": 370,
		},
		LastUpdated: "2026-02-28T10:30:00Z",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(device)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJsoniter_Single tests encoding a single device with jsoniter
func BenchmarkJsoniter_Single(b *testing.B) {
	device := Device{
		ID:    "light.living_room",
		Name:  "Living Room Light",
		Type:  "light",
		State: "on",
		Attributes: map[string]interface{}{
			"brightness": 200,
			"color_temp": 370,
		},
		LastUpdated: "2026-02-28T10:30:00Z",
	}

	var json = jsoniter.ConfigFastest
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(device)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStdlibJSON_50Devices tests encoding 50 devices with stdlib
func BenchmarkStdlibJSON_50Devices(b *testing.B) {
	devices := generateMockDevices()
	response := DeviceListResponse{
		Devices:   devices,
		Total:     50,
		RequestID: "req_benchmark_123",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(response)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJsoniter_50Devices tests encoding 50 devices with jsoniter
func BenchmarkJsoniter_50Devices(b *testing.B) {
	devices := generateMockDevices()
	response := DeviceListResponse{
		Devices:   devices,
		Total:     50,
		RequestID: "req_benchmark_123",
	}

	var json = jsoniter.ConfigFastest
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(response)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStdlibJSON_50Devices_NoAlloc tests encoding with pre-allocated buffer
func BenchmarkStdlibJSON_50Devices_NoAlloc(b *testing.B) {
	devices := generateMockDevices()
	response := DeviceListResponse{
		Devices:   devices,
		Total:     50,
		RequestID: "req_benchmark_123",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(response)
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

// BenchmarkJsoniter_50Devices_NoAlloc tests encoding with jsoniter pre-allocated buffer
func BenchmarkJsoniter_50Devices_NoAlloc(b *testing.B) {
	devices := generateMockDevices()
	response := DeviceListResponse{
		Devices:   devices,
		Total:     50,
		RequestID: "req_benchmark_123",
	}

	var json = jsoniter.ConfigFastest
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(response)
		if err != nil {
			b.Fatal(err)
		}
		_ = data
	}
}

// BenchmarkJsoniter_50Devices_Compatible tests encoding with ConfigCompatibleWithStandardLibrary
func BenchmarkJsoniter_50Devices_Compatible(b *testing.B) {
	devices := generateMockDevices()
	response := DeviceListResponse{
		Devices:   devices,
		Total:     50,
		RequestID: "req_benchmark_123",
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(response)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJsoniter_50Devices_Stream tests encoding with Stream API
func BenchmarkJsoniter_50Devices_Stream(b *testing.B) {
	devices := generateMockDevices()
	response := DeviceListResponse{
		Devices:   devices,
		Total:     50,
		RequestID: "req_benchmark_123",
	}

	var json = jsoniter.ConfigFastest
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stream := json.BorrowStream(nil)
		stream.WriteVal(response)
		if stream.Error != nil {
			b.Fatal(stream.Error)
		}
		_ = stream.Buffer()
		json.ReturnStream(stream)
	}
}

