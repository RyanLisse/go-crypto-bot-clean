package risk

import (
	"testing"
)

func TestRiskLimits_Validate(t *testing.T) {
	tests := []struct {
		name    string
		limits  RiskLimits
		wantErr bool
	}{
		{"valid limits", RiskLimits{1.0, 10000}, false},
		{"zero per-trade", RiskLimits{0, 10000}, true},
		{"negative per-trade", RiskLimits{-1, 10000}, true},
		{"zero total exposure", RiskLimits{1.0, 0}, true},
		{"negative total exposure", RiskLimits{1.0, -100}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.limits.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewExposureTracker(t *testing.T) {
	validLimits := RiskLimits{1.0, 10000}
	_, err := NewExposureTracker(validLimits, 10000)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = NewExposureTracker(RiskLimits{0, 10000}, 10000)
	if err == nil {
		t.Errorf("Expected error for zero per-trade risk")
	}

	_, err = NewExposureTracker(validLimits, 0)
	if err == nil {
		t.Errorf("Expected error for zero account balance")
	}
}

func TestExposureTracker_CanOpenTrade_NormalAndBreach(t *testing.T) {
	limits := RiskLimits{1.0, 10000} // 1% per trade, 10k total
	accountBalance := 10000.0
	tracker, err := NewExposureTracker(limits, accountBalance)
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	// Normal: trade risk = 50 (0.5%), total exposure after = 50
	ok, err := tracker.CanOpenTrade(50)
	if !ok || err != nil {
		t.Errorf("Expected trade within limits, got ok=%v, err=%v", ok, err)
	}

	tracker.AddOrUpdatePosition(Position{Symbol: "BTCUSDT", Size: 0.001, Value: 50})

	// Breach per-trade risk: trade risk = 200 (2%)
	ok, err = tracker.CanOpenTrade(200)
	if ok || err == nil {
		t.Errorf("Expected breach of per-trade risk, got ok=%v, err=%v", ok, err)
	}

	// Breach total exposure: add positions to reach near limit
	tracker.AddOrUpdatePosition(Position{Symbol: "ETHUSDT", Size: 0.01, Value: 9950})
	ok, err = tracker.CanOpenTrade(100)
	if ok || err == nil {
		t.Errorf("Expected breach of total exposure, got ok=%v, err=%v", ok, err)
	}
}

func TestExposureTracker_CanOpenTrade_EdgeCases(t *testing.T) {
	limits := RiskLimits{1.0, 10000}
	tracker, err := NewExposureTracker(limits, 10000)
	if err != nil {
		t.Fatalf("Failed to create tracker: %v", err)
	}

	// Zero position value
	ok, err := tracker.CanOpenTrade(0)
	if ok || err == nil {
		t.Errorf("Expected error for zero position value, got ok=%v, err=%v", ok, err)
	}

	// Negative position value
	ok, err = tracker.CanOpenTrade(-100)
	if ok || err == nil {
		t.Errorf("Expected error for negative position value, got ok=%v, err=%v", ok, err)
	}
}
