package risk

import (
	"math"
	"testing"
)

func TestPositionSize(t *testing.T) {
	type args struct {
		accountBalance   float64
		riskPercentage   float64
		marketVolatility float64
	}
	tests := []struct {
		name     string
		args     args
		wantSize float64
		wantErr  bool
	}{
		{
			name:     "Typical case",
			args:     args{accountBalance: 10000, riskPercentage: 0.01, marketVolatility: 50},
			wantSize: 2,
			wantErr:  false,
		},
		{
			name:     "Zero account balance",
			args:     args{accountBalance: 0, riskPercentage: 0.01, marketVolatility: 50},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name:     "Negative account balance",
			args:     args{accountBalance: -1000, riskPercentage: 0.01, marketVolatility: 50},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name:     "Zero risk percentage",
			args:     args{accountBalance: 10000, riskPercentage: 0, marketVolatility: 50},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name:     "Negative risk percentage",
			args:     args{accountBalance: 10000, riskPercentage: -0.01, marketVolatility: 50},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name:     "Zero market volatility",
			args:     args{accountBalance: 10000, riskPercentage: 0.01, marketVolatility: 0},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name:     "Negative market volatility",
			args:     args{accountBalance: 10000, riskPercentage: 0.01, marketVolatility: -50},
			wantSize: 0,
			wantErr:  true,
		},
		{
			name:     "Very large values",
			args:     args{accountBalance: 1e12, riskPercentage: 0.5, marketVolatility: 1e6},
			wantSize: 5e5,
			wantErr:  false,
		},
		{
			name:     "Very small values",
			args:     args{accountBalance: 1e-6, riskPercentage: 1e-3, marketVolatility: 1e-6},
			wantSize: 1e-9,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSize, err := PositionSize(tt.args.accountBalance, tt.args.riskPercentage, tt.args.marketVolatility)
			if (err != nil) != tt.wantErr {
				t.Errorf("PositionSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(gotSize-tt.wantSize) > 1e-9 {
				t.Errorf("PositionSize() = %v, want %v", gotSize, tt.wantSize)
			}
		})
	}
}
