package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		want     bool
	}{
		{
			name:     "supported currency USD",
			currency: "USD",
			want:     true,
		},
		{
			name:     "supported currency EUR",
			currency: "EUR",
			want:     true,
		},
		{
			name:     "supported currency CAD",
			currency: "CAD",
			want:     true,
		},
		{
			name:     "unsupported currency",
			currency: "unsupported-currency",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSupportedCurrency(tt.currency)
			require.Equal(t, tt.want, got)
		})
	}
}
