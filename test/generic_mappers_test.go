package ezutil_test

import (
	"testing"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFromProtoTime(t *testing.T) {
	t.Run("nil timestamp", func(t *testing.T) {
		result := ezutil.FromProtoTime(nil)
		assert.True(t, result.IsZero())
	})

	t.Run("valid timestamp", func(t *testing.T) {
		now := time.Now()
		proto := timestamppb.New(now)
		result := ezutil.FromProtoTime(proto)
		assert.True(t, now.Equal(result))
	})
}

func TestDecimalToMoney(t *testing.T) {
	tests := []struct {
		name         string
		decimal      decimal.Decimal
		currencyCode string
		expectedUnits int64
		expectedNanos int32
	}{
		{
			name:          "zero value",
			decimal:       decimal.Zero,
			currencyCode:  "USD",
			expectedUnits: 0,
			expectedNanos: 0,
		},
		{
			name:          "whole number",
			decimal:       decimal.NewFromInt(100),
			currencyCode:  "USD",
			expectedUnits: 100,
			expectedNanos: 0,
		},
		{
			name:          "decimal with cents",
			decimal:       decimal.NewFromFloat(123.45),
			currencyCode:  "USD",
			expectedUnits: 123,
			expectedNanos: 450000000,
		},
		{
			name:          "negative value",
			decimal:       decimal.NewFromFloat(-50.25),
			currencyCode:  "EUR",
			expectedUnits: -50,
			expectedNanos: -250000000,
		},
		{
			name:          "value with over 1 billion nanos",
			decimal:       decimal.NewFromFloat(1.5),
			currencyCode:  "USD",
			expectedUnits: 1,
			expectedNanos: 500000000,
		},
		{
			name:          "negative with sign normalization",
			decimal:       decimal.NewFromFloat(-0.5),
			currencyCode:  "USD",
			expectedUnits: 0,
			expectedNanos: -500000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.DecimalToMoney(tt.decimal, tt.currencyCode)
			assert.Equal(t, tt.currencyCode, result.CurrencyCode)
			assert.Equal(t, tt.expectedUnits, result.Units)
			assert.Equal(t, tt.expectedNanos, result.Nanos)
		})
	}
}

func TestMoneyToDecimal(t *testing.T) {
	t.Run("nil money", func(t *testing.T) {
		result := ezutil.MoneyToDecimal(nil)
		assert.True(t, result.Equal(decimal.Zero))
	})

	tests := []struct {
		name     string
		money    *money.Money
		expected decimal.Decimal
	}{
		{
			name: "zero value",
			money: &money.Money{
				CurrencyCode: "USD",
				Units:        0,
				Nanos:        0,
			},
			expected: decimal.Zero,
		},
		{
			name: "whole number",
			money: &money.Money{
				CurrencyCode: "USD",
				Units:        100,
				Nanos:        0,
			},
			expected: decimal.NewFromInt(100),
		},
		{
			name: "with nanos",
			money: &money.Money{
				CurrencyCode: "USD",
				Units:        123,
				Nanos:        450000000,
			},
			expected: decimal.NewFromFloat(123.45),
		},
		{
			name: "negative value",
			money: &money.Money{
				CurrencyCode: "EUR",
				Units:        -50,
				Nanos:        -250000000,
			},
			expected: decimal.NewFromFloat(-50.25),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.MoneyToDecimal(tt.money)
			assert.True(t, tt.expected.Equal(result))
		})
	}
}

func TestDecimalToMoneyRounded(t *testing.T) {
	t.Run("high precision decimal", func(t *testing.T) {
		d, _ := decimal.NewFromString("123.123456789123456")
		result := ezutil.DecimalToMoneyRounded(d, "USD")
		
		assert.Equal(t, "USD", result.CurrencyCode)
		assert.Equal(t, int64(123), result.Units)
		assert.Equal(t, int32(123456789), result.Nanos)
	})

	t.Run("already rounded decimal", func(t *testing.T) {
		d := decimal.NewFromFloat(100.5)
		result := ezutil.DecimalToMoneyRounded(d, "EUR")
		
		assert.Equal(t, "EUR", result.CurrencyCode)
		assert.Equal(t, int64(100), result.Units)
		assert.Equal(t, int32(500000000), result.Nanos)
	})
}

func TestValidateMoney(t *testing.T) {
	t.Run("nil money", func(t *testing.T) {
		err := ezutil.ValidateMoney(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "money cannot be nil")
	})

	t.Run("valid money", func(t *testing.T) {
		m := &money.Money{
			CurrencyCode: "USD",
			Units:        100,
			Nanos:        500000000,
		}
		err := ezutil.ValidateMoney(m)
		assert.NoError(t, err)
	})

	t.Run("nanos out of range - too high", func(t *testing.T) {
		m := &money.Money{
			CurrencyCode: "USD",
			Units:        100,
			Nanos:        1000000000,
		}
		err := ezutil.ValidateMoney(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nanos out of range")
	})

	t.Run("nanos out of range - too low", func(t *testing.T) {
		m := &money.Money{
			CurrencyCode: "USD",
			Units:        100,
			Nanos:        -1000000000,
		}
		err := ezutil.ValidateMoney(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nanos out of range")
	})

	t.Run("mismatched signs - positive units, negative nanos", func(t *testing.T) {
		m := &money.Money{
			CurrencyCode: "USD",
			Units:        100,
			Nanos:        -500000000,
		}
		err := ezutil.ValidateMoney(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have same sign")
	})

	t.Run("mismatched signs - negative units, positive nanos", func(t *testing.T) {
		m := &money.Money{
			CurrencyCode: "USD",
			Units:        -100,
			Nanos:        500000000,
		}
		err := ezutil.ValidateMoney(m)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must have same sign")
	})

	t.Run("zero units with positive nanos", func(t *testing.T) {
		m := &money.Money{
			CurrencyCode: "USD",
			Units:        0,
			Nanos:        500000000,
		}
		err := ezutil.ValidateMoney(m)
		assert.NoError(t, err)
	})

	t.Run("negative units with zero nanos", func(t *testing.T) {
		m := &money.Money{
			CurrencyCode: "USD",
			Units:        -100,
			Nanos:        0,
		}
		err := ezutil.ValidateMoney(m)
		assert.NoError(t, err)
	})
}
