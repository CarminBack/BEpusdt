package epay

import (
	"net/http/httptest"
	"testing"

	"github.com/v03413/bepusdt/app/model"
)

func TestNormalizeTradeType(t *testing.T) {
	if got := normalizeTradeType(usdtBep20Alias); got != string(model.UsdtBep20) {
		t.Fatalf("normalizeTradeType() = %q, want %q", got, model.UsdtBep20)
	}
	if got := normalizeTradeType(string(model.UsdtTrc20)); got != string(model.UsdtTrc20) {
		t.Fatalf("normalizeTradeType() changed existing type to %q", got)
	}
}

func TestRequestScheme(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/submit.php", nil)
	if got := requestScheme(req); got != "http" {
		t.Fatalf("requestScheme() = %q, want http", got)
	}

	req.Header.Set("X-Forwarded-Proto", "https")
	if got := requestScheme(req); got != "https" {
		t.Fatalf("requestScheme() = %q, want https", got)
	}
}

func TestVerifyUSDTBEP20AliasDefaultsToUSD(t *testing.T) {
	params, err := (Epay{}).verify(map[string]string{
		"pid":          Pid,
		"type":         usdtBep20Alias,
		"out_trade_no": "order-1",
		"notify_url":   "https://example.com/notify",
		"return_url":   "https://example.com/return",
		"name":         "Recharge",
		"money":        "10.00",
		"sign":         "placeholder",
	})
	if err != nil {
		t.Fatalf("verify() error = %v", err)
	}
	if params.Fiat != model.USD {
		t.Fatalf("verify() fiat = %q, want %q", params.Fiat, model.USD)
	}
}

func TestVerifyUSDTBEP20AliasKeepsExplicitFiat(t *testing.T) {
	params, err := (Epay{}).verify(map[string]string{
		"pid":          Pid,
		"type":         usdtBep20Alias,
		"out_trade_no": "order-1",
		"notify_url":   "https://example.com/notify",
		"return_url":   "https://example.com/return",
		"name":         "Recharge",
		"money":        "10.00",
		"sign":         "placeholder",
		"fiat":         string(model.CNY),
	})
	if err != nil {
		t.Fatalf("verify() error = %v", err)
	}
	if params.Fiat != model.CNY {
		t.Fatalf("verify() fiat = %q, want %q", params.Fiat, model.CNY)
	}
}
