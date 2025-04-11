package risk

import (
	"math"
	"testing"
)

func TestCalculateSLTP_PercentageBased(t *testing.T) {
	params := RiskParams{
		RiskPercentage: 1.0,
		UseATR:         false,
		RRRatio:        2.0,
	}
	entryPrice := 100.0
	sltp, err := CalculateSLTP(entryPrice, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedSL := 99.0
	expectedTP := 102.0
	if math.Abs(sltp.StopLoss-expectedSL) > 1e-8 {
		t.Errorf("expected SL %.2f, got %.8f", expectedSL, sltp.StopLoss)
	}
	if math.Abs(sltp.TakeProfit-expectedTP) > 1e-8 {
		t.Errorf("expected TP %.2f, got %.8f", expectedTP, sltp.TakeProfit)
	}
}

func TestCalculateSLTP_ATRBased(t *testing.T) {
	params := RiskParams{
		UseATR:         true,
		ATR:            2.5,
		RRRatio:        3.0,
		RiskPercentage: 1.0, // ignored when UseATR=true
	}
	entryPrice := 50.0
	sltp, err := CalculateSLTP(entryPrice, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedSL := 47.5
	expectedTP := 57.5
	if math.Abs(sltp.StopLoss-expectedSL) > 1e-8 {
		t.Errorf("expected SL %.2f, got %.8f", expectedSL, sltp.StopLoss)
	}
	if math.Abs(sltp.TakeProfit-expectedTP) > 1e-8 {
		t.Errorf("expected TP %.2f, got %.8f", expectedTP, sltp.TakeProfit)
	}
}

func TestCalculateSLTP_InvalidParams(t *testing.T) {
	tests := []struct {
		name   string
		params RiskParams
		price  float64
	}{
		{"ZeroEntryPrice", RiskParams{RiskPercentage: 1, RRRatio: 2}, 0},
		{"NegativeEntryPrice", RiskParams{RiskPercentage: 1, RRRatio: 2}, -10},
		{"ZeroRiskPercent", RiskParams{RiskPercentage: 0, RRRatio: 2}, 100},
		{"NegativeRiskPercent", RiskParams{RiskPercentage: -1, RRRatio: 2}, 100},
		{"ExtremeRiskPercent", RiskParams{RiskPercentage: 100, RRRatio: 2}, 100},
		{"ZeroRRRatio", RiskParams{RiskPercentage: 1, RRRatio: 0}, 100},
		{"NegativeRRRatio", RiskParams{RiskPercentage: 1, RRRatio: -2}, 100},
		{"UseATRTrueZeroATR", RiskParams{UseATR: true, ATR: 0, RRRatio: 2, RiskPercentage: 1}, 100},
		{"UseATRTrueNegativeATR", RiskParams{UseATR: true, ATR: -5, RRRatio: 2, RiskPercentage: 1}, 100},
	}
	for _, tc := range tests {
		_, err := CalculateSLTP(tc.price, tc.params)
		if err == nil {
			t.Errorf("%s: expected error, got nil", tc.name)
		}
	}
}

func TestIsStopLossHit(t *testing.T) {
	sltp := SLTPLevel{StopLoss: 95.0, TakeProfit: 110.0}
	if !IsStopLossHit(94.99, sltp) {
		t.Error("expected stop-loss hit at 94.99")
	}
	if !IsStopLossHit(95.0, sltp) {
		t.Error("expected stop-loss hit at 95.0")
	}
	if IsStopLossHit(95.01, sltp) {
		t.Error("did not expect stop-loss hit at 95.01")
	}
}

func TestIsTakeProfitHit(t *testing.T) {
	sltp := SLTPLevel{StopLoss: 95.0, TakeProfit: 110.0}
	if !IsTakeProfitHit(110.0, sltp) {
		t.Error("expected take-profit hit at 110.0")
	}
	if !IsTakeProfitHit(111.0, sltp) {
		t.Error("expected take-profit hit at 111.0")
	}
	if IsTakeProfitHit(109.99, sltp) {
		t.Error("did not expect take-profit hit at 109.99")
	}
}

func TestUpdateSLTP(t *testing.T) {
	params := RiskParams{
		RiskPercentage: 2.0,
		UseATR:         false,
		RRRatio:        1.5,
	}
	entryPrice := 200.0
	sltp, err := UpdateSLTP(entryPrice, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedSL := 196.0
	expectedTP := 206.0
	if math.Abs(sltp.StopLoss-expectedSL) > 1e-8 {
		t.Errorf("expected SL %.2f, got %.8f", expectedSL, sltp.StopLoss)
	}
	if math.Abs(sltp.TakeProfit-expectedTP) > 1e-8 {
		t.Errorf("expected TP %.2f, got %.8f", expectedTP, sltp.TakeProfit)
	}
}

func TestCancelSLTP(t *testing.T) {
	sltp := CancelSLTP()
	if !math.IsNaN(sltp.StopLoss) || !math.IsNaN(sltp.TakeProfit) {
		t.Error("expected NaN values after canceling SLTP")
	}
}
