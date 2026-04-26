// Package ta implements TA-Lib's technical analysis indicators in pure Go.
//
// Conventions:
//   - Output slices have the same length as input.
//   - Indices before an indicator's warmup period are zero.
//   - Period parameters use TA-Lib's default semantics (lookback = period - 1
//     for simple moving averages; period for indicators like ATR/ADX that
//     consume one extra bar for diff computation).
package talib

import (
	"fmt"
	"math"
	"strings"
)

// MaType selects the moving-average kernel used by indicators that accept it.
type MaType int

const (
	SMA MaType = iota
	EMA
	WMA
	DEMA
	TEMA
	TRIMA
	KAMA
	MAMA
	T3MA
)

// MaTypeFromString maps a TA-Lib MA name to a MaType.
// Empty string defaults to SMA.
func MaTypeFromString(s string) (MaType, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "", "SMA":
		return SMA, nil
	case "EMA":
		return EMA, nil
	case "WMA":
		return WMA, nil
	case "DEMA":
		return DEMA, nil
	case "TEMA":
		return TEMA, nil
	case "TRIMA":
		return TRIMA, nil
	case "KAMA":
		return KAMA, nil
	case "MAMA":
		return MAMA, nil
	case "T3", "T3MA":
		return T3MA, nil
	}
	return 0, fmt.Errorf("unknown MA type %q", s)
}

// MA dispatches to the requested moving-average implementation.
// MAMA is not supported via this entry point (it is a two-output indicator).
func MA(real []float64, period int, t MaType) []float64 {
	if period <= 1 {
		out := make([]float64, len(real))
		copy(out, real)
		return out
	}
	switch t {
	case SMA:
		return SMAFn(real, period)
	case EMA:
		return EMAFn(real, period)
	case WMA:
		return WMAFn(real, period)
	case DEMA:
		return DEMAFn(real, period)
	case TEMA:
		return TEMAFn(real, period)
	case TRIMA:
		return TRIMAFn(real, period)
	case KAMA:
		return KAMAFn(real, period)
	case T3MA:
		return T3Fn(real, period, 0.7)
	}
	return SMAFn(real, period)
}

// roundDown returns max(0, n).
func clampZero(n int) int {
	if n < 0 {
		return 0
	}
	return n
}

// trueRange returns the per-bar True Range series:
// max(high-low, |high-prevClose|, |low-prevClose|), with TR[0] = high[0]-low[0].
func trueRange(high, low, close []float64) []float64 {
	n := len(high)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	out[0] = high[0] - low[0]
	for i := 1; i < n; i++ {
		hl := high[i] - low[i]
		hc := math.Abs(high[i] - close[i-1])
		lc := math.Abs(low[i] - close[i-1])
		tr := hl
		if hc > tr {
			tr = hc
		}
		if lc > tr {
			tr = lc
		}
		out[i] = tr
	}
	return out
}
