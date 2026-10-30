package scanner

import (
	"net"
	"runtime"
	"testing"
)

func TestScanPreservesEachRespondingPort(t *testing.T) {
	// Let the scan enqueue both probes before either goroutine runs.
	previous := runtime.GOMAXPROCS(1)
	t.Cleanup(func() { runtime.GOMAXPROCS(previous) })

	ports := make([]int, 0, 2)
	for i := 0; i < 2; i++ {
		listener, err := net.Listen("tcp4", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { listener.Close() })
		ports = append(ports, listener.Addr().(*net.TCPAddr).Port)
	}

	scanner := NewScanner(WithPorts(ports), WithParallelRunners(2))
	if err := scanner.AddIP("127.0.0.1"); err != nil {
		t.Fatal(err)
	}
	results := scanner.Scan()
	if len(results) != len(ports) {
		t.Fatalf("got %d responses, want %d", len(results), len(ports))
	}
	remaining := map[int]bool{ports[0]: true, ports[1]: true}
	for _, result := range results {
		if !result.IP.Equal(net.ParseIP("127.0.0.1")) || result.Protocol != "tcp" {
			t.Fatalf("unexpected response: %+v", result)
		}
		if !remaining[result.Port] {
			t.Fatalf("unexpected or duplicate responding port %d; remaining %v", result.Port, remaining)
		}
		delete(remaining, result.Port)
	}
}
