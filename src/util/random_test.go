package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min, max := int64(1), int64(10)
	result := RandomInt(min, max)

	require.GreaterOrEqual(t, result, min)
	require.LessOrEqual(t, result, max)

}

func TestRandomOwner(t *testing.T) {
	result := RandomOwner()

	require.Len(t, result, 6)
}

func TestRandomMoney(t *testing.T) {
	result := RandomMoney()

	require.GreaterOrEqual(t, result, int64(0))
	require.LessOrEqual(t, result, int64(1000))
}

func TestRandomCurrency(t *testing.T) {
	result := RandomCurrency()

	require.True(t, strings.EqualFold(result, EUR) || strings.EqualFold(result, USD) || strings.EqualFold(result, CAD))
}

func TestRandomEmail(t *testing.T) {
	result := RandomEmail()

	require.Contains(t, result, "@")
	require.Contains(t, result, ".")
}
