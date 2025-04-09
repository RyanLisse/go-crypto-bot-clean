package validation

import (
	"testing"
)

type TestStruct struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Age      int    `validate:"min=0,max=150"`
	Website  string `validate:"url"`
	Category string `validate:"oneof=A B C"`
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name: "valid struct",
			input: TestStruct{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				Website:  "https://example.com",
				Category: "A",
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			input: TestStruct{
				Email:    "john@example.com",
				Age:      30,
				Website:  "https://example.com",
				Category: "A",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: TestStruct{
				Name:     "John Doe",
				Email:    "invalid-email",
				Age:      30,
				Website:  "https://example.com",
				Category: "A",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePositiveFloat(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{"positive value", 10.5, false},
		{"zero value", 0, true},
		{"negative value", -5.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositiveFloat(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositiveFloat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePercentage(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{"valid percentage", 50.0, false},
		{"zero percentage", 0.0, false},
		{"hundred percentage", 100.0, false},
		{"negative percentage", -1.0, true},
		{"over hundred", 101.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePercentage(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePercentage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSymbol(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		wantErr bool
	}{
		{"valid symbol", "BTCUSDT", false},
		{"empty symbol", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSymbol() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOrderType(t *testing.T) {
	tests := []struct {
		name      string
		orderType string
		wantErr   bool
	}{
		{"market order", "market", false},
		{"limit order", "limit", false},
		{"invalid order type", "stop", true},
		{"empty order type", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderType(tt.orderType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrderType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTradeRequest(t *testing.T) {
	price := 100.0
	tests := []struct {
		name      string
		symbol    string
		amount    float64
		orderType string
		price     *float64
		wantErr   bool
	}{
		{
			name:      "valid market order",
			symbol:    "BTCUSDT",
			amount:    1.0,
			orderType: "market",
			price:     nil,
			wantErr:   false,
		},
		{
			name:      "valid limit order",
			symbol:    "BTCUSDT",
			amount:    1.0,
			orderType: "limit",
			price:     &price,
			wantErr:   false,
		},
		{
			name:      "limit order without price",
			symbol:    "BTCUSDT",
			amount:    1.0,
			orderType: "limit",
			price:     nil,
			wantErr:   true,
		},
		{
			name:      "invalid amount",
			symbol:    "BTCUSDT",
			amount:    0.0,
			orderType: "market",
			price:     nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTradeRequest(tt.symbol, tt.amount, tt.orderType, tt.price)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTradeRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRiskParameters(t *testing.T) {
	tests := []struct {
		name           string
		maxDrawdown    float64
		riskPerTrade   float64
		maxExposure    float64
		dailyLossLimit float64
		minBalance     float64
		wantErr        bool
	}{
		{
			name:           "valid parameters",
			maxDrawdown:    20.0,
			riskPerTrade:   2.0,
			maxExposure:    50.0,
			dailyLossLimit: 10.0,
			minBalance:     1000.0,
			wantErr:        false,
		},
		{
			name:           "invalid maxDrawdown",
			maxDrawdown:    150.0,
			riskPerTrade:   2.0,
			maxExposure:    50.0,
			dailyLossLimit: 10.0,
			minBalance:     1000.0,
			wantErr:        true,
		},
		{
			name:           "invalid minBalance",
			maxDrawdown:    20.0,
			riskPerTrade:   2.0,
			maxExposure:    50.0,
			dailyLossLimit: 10.0,
			minBalance:     0.0,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRiskParameters(tt.maxDrawdown, tt.riskPerTrade, tt.maxExposure, tt.dailyLossLimit, tt.minBalance)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRiskParameters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
