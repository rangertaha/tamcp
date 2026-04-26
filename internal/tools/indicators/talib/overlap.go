package talib

import "math"

// SMAFn — Simple Moving Average. Output[i<period-1]=0; Output[i>=period-1]=mean(real[i-period+1..i]).
func SMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	var sum float64
	for i := 0; i < period; i++ {
		sum += real[i]
	}
	out[period-1] = sum / float64(period)
	for i := period; i < n; i++ {
		sum += real[i] - real[i-period]
		out[i] = sum / float64(period)
	}
	return out
}

// EMAFn — Exponential Moving Average using TA-Lib's default smoothing k=2/(period+1).
// Seed value is the SMA of the first `period` samples.
func EMAFn(real []float64, period int) []float64 {
	return ema(real, period, 2.0/float64(period+1))
}

// ema is the EMA primitive with caller-controlled smoothing factor.
// Seed = SMA of first `period` samples; output is zero before the seed index.
func ema(real []float64, period int, k float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	var sum float64
	for i := 0; i < period; i++ {
		sum += real[i]
	}
	prev := sum / float64(period)
	out[period-1] = prev
	for i := period; i < n; i++ {
		prev = (real[i]-prev)*k + prev
		out[i] = prev
	}
	return out
}

// WMAFn — Weighted Moving Average with linearly increasing weights 1..period.
func WMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	denom := float64(period*(period+1)) / 2
	for i := period - 1; i < n; i++ {
		var num float64
		for j := 0; j < period; j++ {
			num += real[i-period+1+j] * float64(j+1)
		}
		out[i] = num / denom
	}
	return out
}

// DEMAFn — Double EMA: 2*EMA(real) - EMA(EMA(real)).
func DEMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < 2*period-1 {
		return out
	}
	e1 := EMAFn(real, period)
	// Build EMA over the populated tail of e1, then re-align into output.
	tail := e1[period-1:]
	e2 := EMAFn(tail, period)
	// e2[period-1] aligns to original index 2*(period-1).
	for i := 2 * (period - 1); i < n; i++ {
		out[i] = 2*e1[i] - e2[i-(period-1)]
	}
	return out
}

// TEMAFn — Triple EMA: 3*E1 - 3*E2 + E3 where E1=EMA, E2=EMA(E1), E3=EMA(E2).
func TEMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	lookback := 3 * (period - 1)
	if period < 1 || n <= lookback {
		return out
	}
	e1 := EMAFn(real, period)
	e2 := EMAFn(e1[period-1:], period)
	e3 := EMAFn(e2[period-1:], period)
	for i := lookback; i < n; i++ {
		out[i] = 3*e1[i] - 3*e2[i-(period-1)] + e3[i-2*(period-1)]
	}
	return out
}

// TRIMAFn — Triangular MA: SMA(SMA(real, ceil(period/2)+1), floor(period/2)+1)
// for even periods; for odd periods, both windows are (period+1)/2.
func TRIMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	var w1, w2 int
	if period%2 == 0 {
		w1 = period/2 + 1
		w2 = period / 2
	} else {
		w1 = (period + 1) / 2
		w2 = (period + 1) / 2
	}
	pass1 := SMAFn(real, w1)
	// SMA the populated tail of pass1.
	tail := pass1[w1-1:]
	pass2 := SMAFn(tail, w2)
	// pass2[w2-1] aligns to original index w1-1+w2-1 = period-1.
	for i := period - 1; i < n; i++ {
		out[i] = pass2[i-(w1-1)]
	}
	return out
}

// KAMAFn — Kaufman Adaptive Moving Average (period for ER; fast=2, slow=30).
func KAMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	const fast, slow = 2.0, 30.0
	fastSC := 2.0 / (fast + 1)
	slowSC := 2.0 / (slow + 1)
	// Seed: SMA(real, period) at index period-1, but TA-Lib seeds at index period.
	// Following TA-Lib semantics: out[period] is first KAMA value; seeded with real[period-1].
	prev := real[period-1]
	out[period] = computeKamaStep(real, period, period, prev, fastSC, slowSC)
	for i := period + 1; i < n; i++ {
		prev = out[i-1]
		out[i] = computeKamaStep(real, i, period, prev, fastSC, slowSC)
	}
	return out
}

func computeKamaStep(real []float64, i, period int, prev, fastSC, slowSC float64) float64 {
	change := math.Abs(real[i] - real[i-period])
	var volatility float64
	for j := i - period + 1; j <= i; j++ {
		volatility += math.Abs(real[j] - real[j-1])
	}
	er := 0.0
	if volatility > 0 {
		er = change / volatility
	}
	sc := er*(fastSC-slowSC) + slowSC
	sc = sc * sc
	return prev + sc*(real[i]-prev)
}

// MIDPOINTFn — (highest(real,period) + lowest(real,period)) / 2.
func MIDPOINTFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		hi, lo := real[i-period+1], real[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if real[j] > hi {
				hi = real[j]
			}
			if real[j] < lo {
				lo = real[j]
			}
		}
		out[i] = (hi + lo) / 2
	}
	return out
}

// MIDPRICEFn — (highest(high,period) + lowest(low,period)) / 2.
func MIDPRICEFn(high, low []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		hi, lo := high[i-period+1], low[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if high[j] > hi {
				hi = high[j]
			}
			if low[j] < lo {
				lo = low[j]
			}
		}
		out[i] = (hi + lo) / 2
	}
	return out
}

// BBANDSFn — Bollinger Bands: middle = MA(real, period, maType); upper/lower = middle ± dev*StdDev.
func BBANDSFn(real []float64, period int, devUp, devDn float64, t MaType) (upper, middle, lower []float64) {
	n := len(real)
	upper = make([]float64, n)
	lower = make([]float64, n)
	middle = MA(real, period, t)
	sd := STDDEVFn(real, period, 1.0)
	for i := 0; i < n; i++ {
		upper[i] = middle[i] + devUp*sd[i]
		lower[i] = middle[i] - devDn*sd[i]
	}
	return
}

// T3Fn — Tillson T3.
// y6 with cascaded EMA chain and T3 mixing coefficients.
func T3Fn(real []float64, period int, vFactor float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	lookback := 6 * (period - 1)
	if period < 1 || n <= lookback {
		return out
	}
	e1 := EMAFn(real, period)
	e2 := EMAFn(e1[period-1:], period)
	e3 := EMAFn(e2[period-1:], period)
	e4 := EMAFn(e3[period-1:], period)
	e5 := EMAFn(e4[period-1:], period)
	e6 := EMAFn(e5[period-1:], period)
	v := vFactor
	c1 := -v * v * v
	c2 := 3*v*v + 3*v*v*v
	c3 := -6*v*v - 3*v - 3*v*v*v
	c4 := 1 + 3*v + v*v*v + 3*v*v
	for i := lookback; i < n; i++ {
		out[i] = c1*e6[i-5*(period-1)] + c2*e5[i-4*(period-1)] + c3*e4[i-3*(period-1)] + c4*e3[i-2*(period-1)]
	}
	return out
}

// SARFn — Parabolic SAR (Wilder). Single-output series.
// acceleration: starting AF (default 0.02). maximum: cap on AF (default 0.2).
func SARFn(high, low []float64, acceleration, maximum float64) []float64 {
	n := len(high)
	out := make([]float64, n)
	if n < 2 {
		return out
	}
	if acceleration <= 0 {
		acceleration = 0.02
	}
	if maximum <= 0 {
		maximum = 0.2
	}

	// Initial trend: compare first +DM vs -DM analog using high/low.
	upMove := high[1] - high[0]
	downMove := low[0] - low[1]
	long := true
	if downMove > upMove && downMove > 0 {
		long = false
	}

	var sar, ep, af float64
	af = acceleration
	if long {
		sar = low[0]
		ep = high[0]
		if high[1] > ep {
			ep = high[1]
		}
	} else {
		sar = high[0]
		ep = low[0]
		if low[1] < ep {
			ep = low[1]
		}
	}
	out[0] = 0
	out[1] = sar
	for i := 2; i < n; i++ {
		sar = sar + af*(ep-sar)
		if long {
			// SAR cannot exceed prior two lows.
			if sar > low[i-1] {
				sar = low[i-1]
			}
			if sar > low[i-2] {
				sar = low[i-2]
			}
			if low[i] < sar {
				// Reverse to short.
				long = false
				sar = ep
				ep = low[i]
				af = acceleration
			} else {
				if high[i] > ep {
					ep = high[i]
					af += acceleration
					if af > maximum {
						af = maximum
					}
				}
			}
		} else {
			if sar < high[i-1] {
				sar = high[i-1]
			}
			if sar < high[i-2] {
				sar = high[i-2]
			}
			if high[i] > sar {
				long = true
				sar = ep
				ep = high[i]
				af = acceleration
			} else {
				if low[i] < ep {
					ep = low[i]
					af += acceleration
					if af > maximum {
						af = maximum
					}
				}
			}
		}
		out[i] = sar
	}
	return out
}
