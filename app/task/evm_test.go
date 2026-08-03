package task

import (
	"encoding/json"
	"testing"
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
