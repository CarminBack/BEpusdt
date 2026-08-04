package task

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBuildLogsRequestFiltersContracts(t *testing.T) {
	e := evm{Contracts: []string{"0xusdt"}}
	body, err := e.buildLogsRequest(evmBlock{From: 10, To: 19})
	if err != nil {
		t.Fatalf("buildLogsRequest() error = %v", err)
	}

	var request struct {
		Method string `json:"method"`
		Params []struct {
			FromBlock string   `json:"fromBlock"`
			ToBlock   string   `json:"toBlock"`
			Address   []string `json:"address"`
			Topics    []string `json:"topics"`
		} `json:"params"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		t.Fatalf("decode request: %v", err)
	}

	filter := request.Params[0]
	if request.Method != "eth_getLogs" || filter.FromBlock != "0xa" || filter.ToBlock != "0x13" {
		t.Fatalf("unexpected request: %s", body)
	}
	if len(filter.Address) != 1 || filter.Address[0] != "0xusdt" {
		t.Fatalf("unexpected address filter: %#v", filter.Address)
	}
	if len(filter.Topics) != 1 || filter.Topics[0] != evmTransferEvent {
		t.Fatalf("unexpected topics: %#v", filter.Topics)
	}
}

func TestBuildLogsRequestWithoutContractsOmitsFilter(t *testing.T) {
	e := evm{}
	body, err := e.buildLogsRequest(evmBlock{From: 1, To: 2})
	if err != nil {
		t.Fatalf("buildLogsRequest() error = %v", err)
	}

	var request struct {
		Params []map[string]interface{} `json:"params"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		t.Fatalf("decode request: %v", err)
	}
	if _, ok := request.Params[0]["address"]; ok {
		t.Fatalf("address filter should be omitted: %s", body)
	}
}

func TestSplitRpcEndpoints(t *testing.T) {
	got := splitRpcEndpoints(" https://rpc-a.example ,https://rpc-b.example\nhttps://rpc-a.example ")
	want := []string{"https://rpc-a.example", "https://rpc-b.example"}
	if len(got) != len(want) {
		t.Fatalf("splitRpcEndpoints() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("splitRpcEndpoints() = %#v, want %#v", got, want)
		}
	}
}

func TestRotateRpcEndpoint(t *testing.T) {
	e := evm{Network: "test"}
	endpoints := []string{"https://rpc-a.example", "https://rpc-b.example", "https://rpc-c.example"}
	if got := e.currentRpcEndpoint(endpoints); got != endpoints[0] {
		t.Fatalf("initial endpoint = %q, want %q", got, endpoints[0])
	}

	if next, ok := e.rotateRpcEndpointIn(endpoints, endpoints[0]); !ok || next != endpoints[1] {
		t.Fatalf("rotateRpcEndpointIn() = %q, %t; want %q, true", next, ok, endpoints[1])
	}
	if got := e.currentRpcEndpoint(endpoints); got != endpoints[1] {
		t.Fatalf("rotated endpoint = %q, want %q", got, endpoints[1])
	}
	if _, ok := e.rotateRpcEndpointIn(endpoints, endpoints[0]); ok {
		t.Fatal("stale failure should not rotate the current endpoint")
	}
}

func TestLookbackBlockCountSupportsSubSecondBscBlocks(t *testing.T) {
	now := time.Date(2026, time.August, 4, 16, 0, 0, 0, time.UTC)
	e := evm{AvgBlockTime: 750 * time.Millisecond}
	if got, want := e.lookbackBlockCount(now.Add(-20*time.Minute), now), int64(1630); got != want {
		t.Fatalf("lookbackBlockCount() = %d, want %d", got, want)
	}
}
