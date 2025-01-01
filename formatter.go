package hscale

// version: 0.0.1

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

type InvalidMinValueError struct {
	Value string
	Err   error
}

type DivisionError struct {
	Number    string
	ScaleName string
	Err       error
}

type InvalidNumberError struct {
	Value string
	Err   error
}

type FloorError struct {
	Value string
	Err   error
}
type Scale struct {
	Value string
	Name  string
}

func (e InvalidNumberError) Error() string {
	return fmt.Sprintf("invalid number %q: %v", e.Value, e.Err)
}

func (e InvalidMinValueError) Error() string {
	return fmt.Sprintf("invalid min value %q: %v", e.Value, e.Err)
}

func (e DivisionError) Error() string {
	return fmt.Sprintf("division error for number=%q scale=%q: %v",
		e.Number, e.ScaleName, e.Err)
}

func (e FloorError) Error() string {
	return fmt.Sprintf("floor error for %q: %v", e.Value, e.Err)
}

type RoundingError struct {
	Value string
	Err   error
}

func (e RoundingError) Error() string {
	return fmt.Sprintf("rounding error for %q: %v", e.Value, e.Err)
}

func Formatter(
	num string,
	min string,
	scales []Scale,
	fallback func(string) string,
) (string, error) {

	if countTrailingZeros(num) <= 2 {
		return fallback(num), nil
	}

	number := new(apd.Decimal)
	if _, _, err := number.SetString(num); err != nil {
		return "", InvalidNumberError{Value: num, Err: err}
	}

	minValue := new(apd.Decimal)
	if _, _, err := minValue.SetString(min); err != nil {
		return "", InvalidMinValueError{Value: min, Err: err}
	}

	if number.Cmp(minValue) < 0 {
		return fallback(num), nil
	}

	ctx := apd.BaseContext.WithPrecision(17)

	for _, sc := range scales {
		scaleVal := new(apd.Decimal)
		if _, _, err := scaleVal.SetString(sc.Value); err != nil {
			continue
		}

		if number.Cmp(scaleVal) >= 0 {
			ratio := new(apd.Decimal)

			if _, err := ctx.Quo(ratio, number, scaleVal); err != nil {
				return "", DivisionError{
					Number:    num,
					ScaleName: sc.Name,
					Err:       err,
				}
			}

			if sc.Name == "thousand" {
				tempFloor := new(apd.Decimal)
				if _, err := ctx.Floor(tempFloor, ratio); err != nil {
					return "", FloorError{Value: ratio.String(), Err: err}
				}
				if ratio.Cmp(tempFloor) != 0 {
					return fallback(num), nil
				}
			}

			ratioRounded := new(apd.Decimal)
			temp := new(apd.Decimal)

			ctx.Mul(temp, ratio, apd.New(1000, 0))
			if _, err := ctx.RoundToIntegralExact(temp, temp); err != nil {
				return "", RoundingError{Value: ratio.String(), Err: err}
			}
			ctx.Quo(ratioRounded, temp, apd.New(1000, 0))

			reconstructed := new(apd.Decimal)
			ctx.Mul(reconstructed, ratioRounded, scaleVal)
			if reconstructed.Cmp(number) != 0 {
				return fallback(num), nil
			}

			result := formatWithUpTo3Decimal(ratio)
			return fmt.Sprintf("%s %s", result, sc.Name), nil
		}
	}

	return fallback(num), nil
}

func formatWithUpTo3Decimal(dec *apd.Decimal) string {
	ctx := apd.BaseContext.WithPrecision(17)

	temp := new(apd.Decimal)
	ctx.Mul(temp, dec, apd.New(1000, 0))

	ctx.RoundToIntegralExact(temp, temp)

	scaled := new(apd.Decimal)
	ctx.Quo(scaled, temp, apd.New(1000, 0))

	return stripTrailingZeros(scaled.String())
}

func countTrailingZeros(numStr string) int {
	count := 0
	for i := len(numStr) - 1; i >= 0; i-- {
		if numStr[i] == '0' {
			count++
		} else {
			break
		}
	}
	return count
}

func stripTrailingZeros(s string) string {
	if strings.Contains(s, ".") {
		for strings.HasSuffix(s, "0") {
			s = s[:len(s)-1]
		}
		if strings.HasSuffix(s, ".") {
			s = s[:len(s)-1]
		}
	}
	return s
}
