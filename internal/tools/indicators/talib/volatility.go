package talib

// TRANGEFn — True Range per bar.
func TRANGEFn(high, low, close []float64) []float64 {
	return trueRange(high, low, close)
}

// ATRFn — Average True Range using Wilder smoothing (alpha = 1/period).
// Seed = SMA of TR over the first `period` bars at index `period`.
func ATRFn(high, low, close []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	tr := trueRange(high, low, close)
	var sum float64
	for i := 1; i <= period; i++ {
		sum += tr[i]
	}
	prev := sum / float64(period)
	out[period] = prev
	pf := float64(period)
	for i := period + 1; i < n; i++ {
		prev = (prev*(pf-1) + tr[i]) / pf
		out[i] = prev
	}
	return out
}

// NATRFn — Normalized ATR: 100 * ATR / close.
func NATRFn(high, low, close []float64, period int) []float64 {
	atr := ATRFn(high, low, close, period)
	out := make([]float64, len(atr))
	for i, v := range atr {
		if close[i] != 0 {
			out[i] = 100 * v / close[i]
		}
	}
	return out
}
