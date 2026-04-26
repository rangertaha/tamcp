package talib

import "math"

// Final batch of multi-bar candlestick patterns. Same -100/0/+100 convention.

// CDLADVANCEBLOCKFn — Three bullish bars losing momentum (each body smaller, upper shadows growing). Bearish.
func CDLADVANCEBLOCKFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		if !bullish(o1, c1) || !bullish(o2, c2) || !bullish(o3, c3) {
			continue
		}
		// Each opens within prior body and closes higher.
		if !(o2 > o1 && o2 < c1 && c2 > c1) || !(o3 > o2 && o3 < c2 && c3 > c2) {
			continue
		}
		b1, b2, b3 := body(o1, c1), body(o2, c2), body(o3, c3)
		if !(b2 < b1 && b3 < b2) {
			continue
		}
		ab := avgBody(open, close, i-2)
		if ab <= 0 || b1 < longBodyRatio*ab {
			continue
		}
		// Upper shadows getting longer signals exhaustion.
		us2 := upperShadow(o2, high[i-1], c2)
		us3 := upperShadow(o3, high[i], c3)
		if us3 > b3*0.3 || us2 > b2*0.3 {
			out[i] = -100
		}
	}
	return out
}

// CDLBREAKAWAYFn — Five-bar pattern: long body, gap, three same-direction bars,
// then opposite-direction bar closing back into the gap.
func CDLBREAKAWAYFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 4; i < n; i++ {
		o1, c1 := open[i-4], close[i-4]
		o2, c2 := open[i-3], close[i-3]
		o3, c3 := open[i-2], close[i-2]
		o4, c4 := open[i-1], close[i-1]
		o5, c5 := open[i], close[i]
		ab := avgBody(open, close, i-4)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		// Bullish breakaway: long bear, gap-down bear, two more bears, bull closing into gap.
		if bearish(o1, c1) && bearish(o2, c2) && high[i-3] < low[i-4] &&
			bearish(o3, c3) && bearish(o4, c4) && c4 < c3 && c3 < c2 &&
			bullish(o5, c5) && c5 > c2 && c5 < o2 {
			out[i] = 100
		}
		if bullish(o1, c1) && bullish(o2, c2) && low[i-3] > high[i-4] &&
			bullish(o3, c3) && bullish(o4, c4) && c4 > c3 && c3 > c2 &&
			bearish(o5, c5) && c5 < c2 && c5 > o2 {
			out[i] = -100
		}
	}
	return out
}

// CDLCONCEALBABYSWALLFn — Four-bar bullish reversal: two black marubozus, then small black with upper shadow above the prior high, then large black engulfing prior shadow.
func CDLCONCEALBABYSWALLFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	maru := CDLMARUBOZUFn(open, high, low, close)
	for i := 3; i < n; i++ {
		// Bars 1 and 2: black marubozus.
		if maru[i-3] != -100 || maru[i-2] != -100 {
			continue
		}
		o3, c3 := open[i-1], close[i-1]
		o4, c4 := open[i], close[i]
		if !bearish(o3, c3) || !bearish(o4, c4) {
			continue
		}
		// Bar 3: small black with upper shadow above bar 2 high.
		if upperShadow(o3, high[i-1], c3) <= 0 || high[i-1] <= high[i-2] {
			continue
		}
		// Bar 4: long black engulfing bar 3's shadow.
		if !(o4 > high[i-1] && c4 < low[i-1]) {
			continue
		}
		out[i] = 100
	}
	return out
}

// CDLHIKKAKEFn — 5-bar inside-bar false-break reversal.
func CDLHIKKAKEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 4; i < n; i++ {
		// Bar 2 is inside bar 1 (lower high, higher low).
		h1, l1 := high[i-3], low[i-3]
		h2, l2 := high[i-2], low[i-2]
		if !(h2 < h1 && l2 > l1) {
			continue
		}
		h3 := high[i-1]
		l3 := low[i-1]
		// Bullish hikkake: bar 3 breaks below bar 2 low; bar 4/5 reclaims above bar 2 high.
		if l3 < l2 && (close[i-1] >= h2 || close[i] >= h2) {
			out[i] = 100
			continue
		}
		// Bearish hikkake: bar 3 breaks above bar 2 high; bar 4/5 fails below bar 2 low.
		if h3 > h2 && (close[i-1] <= l2 || close[i] <= l2) {
			out[i] = -100
		}
	}
	return out
}

// CDLHIKKAKEMODFn — Hikkake with prior-trend confirmation: requires that the
// inside bar (bar 2) be preceded by a similarly-colored bar in the same trend.
func CDLHIKKAKEMODFn(open, high, low, close []float64) []int {
	hk := CDLHIKKAKEFn(open, high, low, close)
	out := make([]int, len(hk))
	for i := 5; i < len(hk); i++ {
		if hk[i] == 0 {
			continue
		}
		// Confirm with bar i-5 trend direction.
		prevTrend := close[i-5]
		if hk[i] == 100 && prevTrend > close[i-3] {
			out[i] = 100
		} else if hk[i] == -100 && prevTrend < close[i-3] {
			out[i] = -100
		}
	}
	return out
}

// CDLLADDERBOTTOMFn — Five-bar bullish reversal: three bears with consecutively lower opens
// and closes, fourth bear with upper shadow, fifth bull opening above bar 4's body.
func CDLLADDERBOTTOMFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 4; i < n; i++ {
		o1, c1 := open[i-4], close[i-4]
		o2, c2 := open[i-3], close[i-3]
		o3, c3 := open[i-2], close[i-2]
		o4, c4 := open[i-1], close[i-1]
		o5, c5 := open[i], close[i]
		if !bearish(o1, c1) || !bearish(o2, c2) || !bearish(o3, c3) || !bearish(o4, c4) {
			continue
		}
		if !(o2 < o1 && o3 < o2) || !(c2 < c1 && c3 < c2 && c4 < c3) {
			continue
		}
		// Bar 4 has upper shadow.
		if upperShadow(o4, high[i-1], c4) <= 0 {
			continue
		}
		// Bar 5: bullish, opening above bar 4 body, closing higher.
		if bullish(o5, c5) && o5 > math.Max(o4, c4) && c5 > o5 {
			out[i] = 100
		}
	}
	return out
}

// CDLMATHOLDFn — Five-bar bullish continuation: long bull, then three small bears
// staying within the first bull's range, then long bull breaking out.
func CDLMATHOLDFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 4; i < n; i++ {
		o1, c1 := open[i-4], close[i-4]
		o2, c2 := open[i-3], close[i-3]
		o3, c3 := open[i-2], close[i-2]
		o4, c4 := open[i-1], close[i-1]
		o5, c5 := open[i], close[i]
		ab := avgBody(open, close, i-4)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		if !bullish(o1, c1) || !bullish(o5, c5) || body(o5, c5) < longBodyRatio*ab {
			continue
		}
		if !(bearish(o2, c2) && bearish(o3, c3) && bearish(o4, c4)) {
			continue
		}
		// All three small bodies stay within bar 1's range.
		if !(low[i-3] > l1Range(o1, c1) && low[i-2] > l1Range(o1, c1) && low[i-1] > l1Range(o1, c1)) {
			continue
		}
		if c5 > c1 {
			out[i] = 100
		}
	}
	return out
}

func l1Range(o, c float64) float64 {
	if o < c {
		return o
	}
	return c
}

// CDLRISEFALL3METHODSFn — Five-bar continuation: long body, three small opposite bodies inside,
// then a long same-direction body breaking out.
func CDLRISEFALL3METHODSFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 4; i < n; i++ {
		o1, c1 := open[i-4], close[i-4]
		o2, c2 := open[i-3], close[i-3]
		o3, c3 := open[i-2], close[i-2]
		o4, c4 := open[i-1], close[i-1]
		o5, c5 := open[i], close[i]
		ab := avgBody(open, close, i-4)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab || body(o5, c5) < longBodyRatio*ab {
			continue
		}
		// Bullish rise: bull, three bears inside bar 1 range, bull breaking higher.
		if bullish(o1, c1) && bullish(o5, c5) && c5 > c1 &&
			bearish(o2, c2) && bearish(o3, c3) && bearish(o4, c4) &&
			high[i-3] < high[i-4] && low[i-3] > low[i-4] &&
			high[i-2] < high[i-4] && low[i-2] > low[i-4] &&
			high[i-1] < high[i-4] && low[i-1] > low[i-4] &&
			c2 > c3 && c3 > c4 {
			out[i] = 100
		}
		// Bearish fall: mirror.
		if bearish(o1, c1) && bearish(o5, c5) && c5 < c1 &&
			bullish(o2, c2) && bullish(o3, c3) && bullish(o4, c4) &&
			high[i-3] < high[i-4] && low[i-3] > low[i-4] &&
			high[i-2] < high[i-4] && low[i-2] > low[i-4] &&
			high[i-1] < high[i-4] && low[i-1] > low[i-4] &&
			c2 < c3 && c3 < c4 {
			out[i] = -100
		}
	}
	return out
}

// CDL3STARSINSOUTHFn — Rare bullish reversal: three black candles with progressively smaller
// bodies and ranges, in a downtrend.
func CDL3STARSINSOUTHFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		if !bearish(o1, c1) || !bearish(o2, c2) || !bearish(o3, c3) {
			continue
		}
		ab := avgBody(open, close, i-2)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		// Bar 1 has long lower shadow.
		if lowerShadow(o1, low[i-2], c1) <= body(o1, c1)*0.3 {
			continue
		}
		// Bar 2 is smaller and stays inside bar 1's body extent.
		if !(body(o2, c2) < body(o1, c1) && low[i-1] > low[i-2] && high[i-1] < o1) {
			continue
		}
		// Bar 3 is short, no shadows, range inside bar 2.
		if body(o3, c3) > shortBodyRatio*ab {
			continue
		}
		if upperShadow(o3, high[i], c3) > 0.05*body(o3, c3) ||
			lowerShadow(o3, low[i], c3) > 0.05*body(o3, c3) {
			continue
		}
		if high[i] >= high[i-1] || low[i] <= low[i-1] {
			continue
		}
		out[i] = 100
	}
	return out
}

// CDLXSIDEGAP3METHODSFn — Five-bar continuation with a gap: 3 same-direction bars, gap,
// then opposite-color bar filling the gap, then same-direction continuation.
func CDLXSIDEGAP3METHODSFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 4; i < n; i++ {
		o1, c1 := open[i-4], close[i-4]
		o2, c2 := open[i-3], close[i-3]
		o3, c3 := open[i-2], close[i-2]
		o4, c4 := open[i-1], close[i-1]
		o5, c5 := open[i], close[i]
		// Bullish: three bulls with gap up between bar 2 and 3, bear filling gap, bull continuation.
		if bullish(o1, c1) && bullish(o2, c2) && bullish(o3, c3) &&
			low[i-2] > high[i-3] &&
			bearish(o4, c4) && o4 > c3 && c4 < o3 && c4 > c2 &&
			bullish(o5, c5) && c5 > high[i-2] {
			out[i] = 100
		}
		if bearish(o1, c1) && bearish(o2, c2) && bearish(o3, c3) &&
			high[i-2] < low[i-3] &&
			bullish(o4, c4) && o4 < c3 && c4 > o3 && c4 < c2 &&
			bearish(o5, c5) && c5 < low[i-2] {
			out[i] = -100
		}
	}
	return out
}
