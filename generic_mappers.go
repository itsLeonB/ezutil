package ezutil

import (
	"time"

	"github.com/itsLeonB/ungerr"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromProtoTime(t *timestamppb.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.AsTime()
}

func DecimalToMoney(d decimal.Decimal, currencyCode string) *money.Money {
	units := d.Truncate(0).IntPart()
	fractional := d.Sub(decimal.New(units, 0))
	nanos := fractional.Mul(decimal.New(1000000000, 0)).Round(0).IntPart()

	// Normalize: fold whole billions into units
	units += nanos / 1_000_000_000
	nanos = nanos % 1_000_000_000

	// Adjust signs so units and nanos have same overall sign
	if nanos < 0 && units > 0 {
		units--
		nanos += 1_000_000_000
	} else if nanos > 0 && units < 0 {
		units++
		nanos -= 1_000_000_000
	}

	// Validate nanos bounds and int32 range
	if nanos < -999_999_999 || nanos > 999_999_999 || nanos < int64(int32(-2147483648)) || nanos > int64(int32(2147483647)) {
		// Clamp to valid range
		if nanos > 999_999_999 {
			nanos = 999_999_999
		} else if nanos < -999_999_999 {
			nanos = -999_999_999
		}
	}

	return &money.Money{
		CurrencyCode: currencyCode,
		Units:        units,
		Nanos:        int32(nanos),
	}
}

// MoneyToDecimal converts google.type.Money to decimal.Decimal
func MoneyToDecimal(m *money.Money) decimal.Decimal {
	if m == nil {
		return decimal.Zero
	}

	// Convert units to decimal
	units := decimal.New(m.Units, 0)

	// Convert nanos to decimal (nanos / 10^9)
	nanos := decimal.New(int64(m.Nanos), -9)

	return units.Add(nanos)
}

// DecimalToMoneyRounded converts decimal to Money with proper rounding
// This handles cases where decimal has more precision than nanos can represent
func DecimalToMoneyRounded(d decimal.Decimal, currencyCode string) *money.Money {
	// Round to 9 decimal places (nano precision)
	rounded := d.Round(9)
	return DecimalToMoney(rounded, currencyCode)
}

// ValidateMoney checks if Money values are valid
func ValidateMoney(m *money.Money) error {
	if m == nil {
		return ungerr.Unknown("money cannot be nil")
	}

	// Nanos must be in range [-999,999,999, 999,999,999]
	if m.Nanos < -999999999 || m.Nanos > 999999999 {
		return ungerr.Unknownf("nanos out of range: %d", m.Nanos)
	}

	// Units and nanos must have the same sign (or one of them is zero)
	if (m.Units > 0 && m.Nanos < 0) || (m.Units < 0 && m.Nanos > 0) {
		return ungerr.Unknownf("units (%d) and nanos (%d) must have same sign", m.Units, m.Nanos)
	}

	return nil
}
