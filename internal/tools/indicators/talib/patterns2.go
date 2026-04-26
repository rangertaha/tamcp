package talib

import "math"

// Additional candlestick patterns. Same -100/0/+100 output convention as patterns.go.

// CDLBELTHOLDFn — Long body opening at one extreme of the range.
// Bullish: long bull with no lower shadow (open at low).
// Bearish: long bear with no upper shadow (open at high).
func CDLBELTHOLDFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		ab := avgBody(open, close, i)
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if ab <= 0 || rng <= 0 || b < longBodyRatio*ab {
			continue
		}
		if bullish(open[i], close[i]) && lowerShadow(open[i], low[i], close[i]) <= 0.1*rng {
			out[i] = 100
		} else if bearish(open[i], close[i]) && upperShadow(open[i], high[i], close[i]) <= 0.1*rng {
			out[i] = -100
		}
	}
	return out
}

// CDLCLOSINGMARUBOZUFn — Long body closing at one extreme.
// Bullish: long bull closing at high (no upper shadow).
// Bearish: long bear closing at low (no lower shadow).
func CDLCLOSINGMARUBOZUFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		ab := avgBody(open, close, i)
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if ab <= 0 || rng <= 0 || b < longBodyRatio*ab {
			continue
		}
		if bullish(open[i], close[i]) && upperShadow(open[i], high[i], close[i]) <= 0.1*rng {
			out[i] = 100
		} else if bearish(open[i], close[i]) && lowerShadow(open[i], low[i], close[i]) <= 0.1*rng {
			out[i] = -100
		}
	}
	return out
}

// CDLHIGHWAVEFn — Small body with very long upper and lower shadows.
func CDLHIGHWAVEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		ab := avgBody(open, close, i)
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if ab <= 0 || rng <= 0 || b > shortBodyRatio*ab {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if us > 2*b && ls > 2*b && rng > 2*ab {
			if bullish(open[i], close[i]) {
				out[i] = 100
			} else if bearish(open[i], close[i]) {
				out[i] = -100
			}
		}
	}
	return out
}

// CDLRICKSHAWMANFn — Long-legged doji with body near the middle of the range.
func CDLRICKSHAWMANFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 || b > dojiBodyRatio*rng {
			continue
		}
		mid := (high[i] + low[i]) / 2
		bodyMid := (open[i] + close[i]) / 2
		if math.Abs(bodyMid-mid) <= 0.1*rng {
			out[i] = 100
		}
	}
	return out
}

// CDLTAKURIFn — Dragonfly doji with an unusually long lower shadow.
func CDLTAKURIFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 || b > dojiBodyRatio*rng {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if us <= 0.05*rng && ls >= 0.7*rng {
			out[i] = 100
		}
	}
	return out
}

// CDLDOJISTARFn — Long body followed by a doji that gaps in the trend direction.
func CDLDOJISTARFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		ab := avgBody(open, close, i-1)
		if ab <= 0 || body(open[i-1], close[i-1]) < longBodyRatio*ab {
			continue
		}
		rng := candleRange(high[i], low[i])
		if rng <= 0 || body(open[i], close[i]) > dojiBodyRatio*rng {
			continue
		}
		// Bullish setup: prior bear, doji gaps below.
		if bearish(open[i-1], close[i-1]) && math.Max(open[i], close[i]) < close[i-1] {
			out[i] = 100
		} else if bullish(open[i-1], close[i-1]) && math.Min(open[i], close[i]) > close[i-1] {
			out[i] = -100
		}
	}
	return out
}

// CDLHOMINGPIGEONFn — Bullish: two consecutive bears, second body inside first.
func CDLHOMINGPIGEONFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		if !bearish(o1, c1) || !bearish(o2, c2) {
			continue
		}
		if o2 < o1 && c2 > c1 {
			out[i] = 100
		}
	}
	return out
}

// CDLINNECKFn — Bearish continuation: long bear, then small bull closing slightly into prior body.
func CDLINNECKFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		ab := avgBody(open, close, i-1)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		if bearish(o1, c1) && bullish(o2, c2) && o2 < c1 && c2 >= c1 && c2 < c1+0.05*body(o1, c1) {
			out[i] = -100
		}
	}
	return out
}

// CDLONNECKFn — Bearish continuation similar to In-Neck but second close exactly at prior low.
func CDLONNECKFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		ab := avgBody(open, close, i-1)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		if bearish(o1, c1) && bullish(o2, c2) && o2 < c1 && math.Abs(c2-low[i-1]) <= 0.02*c1 {
			out[i] = -100
		}
	}
	return out
}

// CDLTHRUSTINGFn — Bearish continuation: long bear, then bull closing into lower half of prior body.
func CDLTHRUSTINGFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		ab := avgBody(open, close, i-1)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		mid := (o1 + c1) / 2
		if bearish(o1, c1) && bullish(o2, c2) && o2 < c1 && c2 > c1 && c2 < mid {
			out[i] = -100
		}
	}
	return out
}

// CDLMATCHINGLOWFn — Bullish: two bears with the same closing price.
func CDLMATCHINGLOWFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		if bearish(o1, c1) && bearish(o2, c2) && math.Abs(c1-c2) <= 0.005*c1 {
			out[i] = 100
		}
	}
	return out
}

// CDLSEPARATINGLINESFn — Continuation: two bars opening at the same price in opposite directions.
func CDLSEPARATINGLINESFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		if math.Abs(o1-o2) > 0.005*o1 {
			continue
		}
		if bearish(o1, c1) && bullish(o2, c2) {
			out[i] = 100
		} else if bullish(o1, c1) && bearish(o2, c2) {
			out[i] = -100
		}
	}
	return out
}

// CDLCOUNTERATTACKFn — Reversal: opposite-color bars closing at the same level.
func CDLCOUNTERATTACKFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		if math.Abs(c1-c2) > 0.005*c1 {
			continue
		}
		ab := avgBody(open, close, i-1)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab || body(o2, c2) < longBodyRatio*ab {
			continue
		}
		if bearish(o1, c1) && bullish(o2, c2) {
			out[i] = 100
		} else if bullish(o1, c1) && bearish(o2, c2) {
			out[i] = -100
		}
	}
	return out
}

// CDLKICKINGFn — Two opposite marubozus separated by a gap.
func CDLKICKINGFn(open, high, low, close []float64) []int {
	n := len(open)
	maru := CDLMARUBOZUFn(open, high, low, close)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		if maru[i-1] == 0 || maru[i] == 0 {
			continue
		}
		// Bullish: bear marubozu then bull marubozu gapping up.
		if maru[i-1] == -100 && maru[i] == 100 && low[i] > high[i-1] {
			out[i] = 100
		} else if maru[i-1] == 100 && maru[i] == -100 && high[i] < low[i-1] {
			out[i] = -100
		}
	}
	return out
}

// CDLKICKINGBYLENGTHFn — Same as Kicking but the longer marubozu sets the direction.
func CDLKICKINGBYLENGTHFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	maru := CDLMARUBOZUFn(open, high, low, close)
	for i := 1; i < n; i++ {
		if maru[i-1] == 0 || maru[i] == 0 {
			continue
		}
		if !((maru[i-1] == -100 && maru[i] == 100 && low[i] > high[i-1]) ||
			(maru[i-1] == 100 && maru[i] == -100 && high[i] < low[i-1])) {
			continue
		}
		if body(open[i], close[i]) >= body(open[i-1], close[i-1]) {
			if bullish(open[i], close[i]) {
				out[i] = 100
			} else {
				out[i] = -100
			}
		} else {
			if bullish(open[i-1], close[i-1]) {
				out[i] = 100
			} else {
				out[i] = -100
			}
		}
	}
	return out
}

// CDL2CROWSFn — Bearish reversal: bull, gap-up bear, bear closing inside first body.
func CDL2CROWSFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		ab := avgBody(open, close, i-2)
		if ab <= 0 || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		if !bullish(o1, c1) || !bearish(o2, c2) || !bearish(o3, c3) {
			continue
		}
		if math.Min(o2, c2) > c1 && o3 > c2 && c3 < c1 && c3 > o1 {
			out[i] = -100
		}
	}
	return out
}

// CDL3INSIDEFn — Three Inside Up/Down: harami + same-direction confirmation bar.
func CDL3INSIDEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	har := CDLHARAMIFn(open, high, low, close)
	for i := 2; i < n; i++ {
		if har[i-1] == 0 {
			continue
		}
		if har[i-1] == 100 && bullish(open[i], close[i]) && close[i] > close[i-1] {
			out[i] = 100
		} else if har[i-1] == -100 && bearish(open[i], close[i]) && close[i] < close[i-1] {
			out[i] = -100
		}
	}
	return out
}

// CDL3OUTSIDEFn — Three Outside Up/Down: engulfing + same-direction confirmation bar.
func CDL3OUTSIDEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	eng := CDLENGULFINGFn(open, high, low, close)
	for i := 2; i < n; i++ {
		if eng[i-1] == 0 {
			continue
		}
		if eng[i-1] == 100 && bullish(open[i], close[i]) && close[i] > close[i-1] {
			out[i] = 100
		} else if eng[i-1] == -100 && bearish(open[i], close[i]) && close[i] < close[i-1] {
			out[i] = -100
		}
	}
	return out
}

// CDLEVENINGDOJISTARFn — Evening Star where the middle bar is a doji.
func CDLEVENINGDOJISTARFn(open, high, low, close []float64) []int {
	es := CDLEVENINGSTARFn(open, high, low, close)
	doji := CDLDOJIFn(open, high, low, close)
	out := make([]int, len(es))
	for i := range out {
		if i >= 1 && es[i] != 0 && doji[i-1] != 0 {
			out[i] = es[i]
		}
	}
	return out
}

// CDLMORNINGDOJISTARFn — Morning Star where the middle bar is a doji.
func CDLMORNINGDOJISTARFn(open, high, low, close []float64) []int {
	ms := CDLMORNINGSTARFn(open, high, low, close)
	doji := CDLDOJIFn(open, high, low, close)
	out := make([]int, len(ms))
	for i := range out {
		if i >= 1 && ms[i] != 0 && doji[i-1] != 0 {
			out[i] = ms[i]
		}
	}
	return out
}

// CDLABANDONEDBABYFn — Three-bar reversal with a gapped doji isolated by both neighbors.
func CDLABANDONEDBABYFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	doji := CDLDOJIFn(open, high, low, close)
	for i := 2; i < n; i++ {
		if doji[i-1] == 0 {
			continue
		}
		o1, c1 := open[i-2], close[i-2]
		o3, c3 := open[i], close[i]
		// Bullish: long bear, gap-down doji isolated below, then gap-up bull above doji high.
		if bearish(o1, c1) && high[i-1] < low[i-2] && bullish(o3, c3) && low[i] > high[i-1] {
			out[i] = 100
		} else if bullish(o1, c1) && low[i-1] > high[i-2] && bearish(o3, c3) && high[i] < low[i-1] {
			out[i] = -100
		}
	}
	return out
}

// CDLIDENTICAL3CROWSFn — Three black crows where each opens at the prior close.
func CDLIDENTICAL3CROWSFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		if !bearish(o1, c1) || !bearish(o2, c2) || !bearish(o3, c3) {
			continue
		}
		if c2 >= c1 || c3 >= c2 {
			continue
		}
		if math.Abs(o2-c1) > 0.005*c1 || math.Abs(o3-c2) > 0.005*c2 {
			continue
		}
		out[i] = -100
	}
	return out
}

// CDLSTALLEDPATTERNFn — Three white soldiers losing momentum (small last body, gap up).
func CDLSTALLEDPATTERNFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		if !bullish(o1, c1) || !bullish(o2, c2) || !bullish(o3, c3) {
			continue
		}
		if !(c2 > c1 && c3 > c2 && o2 > o1 && o3 > o2) {
			continue
		}
		ab := avgBody(open, close, i-2)
		if ab <= 0 {
			continue
		}
		if body(o1, c1) >= longBodyRatio*ab && body(o2, c2) >= longBodyRatio*ab && body(o3, c3) <= shortBodyRatio*ab {
			out[i] = -100
		}
	}
	return out
}

// CDLSTICKSANDWICHFn — Bullish: bear, bull, bear with the two bears closing at the same price.
func CDLSTICKSANDWICHFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		if !bearish(o1, c1) || !bullish(o2, c2) || !bearish(o3, c3) {
			continue
		}
		if math.Abs(c1-c3) > 0.005*c1 {
			continue
		}
		if low[i-1] > c1 {
			out[i] = 100
		}
	}
	return out
}

// CDLTRISTARFn — Three consecutive dojis. Direction inferred from gap pattern.
func CDLTRISTARFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	doji := CDLDOJIFn(open, high, low, close)
	for i := 2; i < n; i++ {
		if doji[i-2] == 0 || doji[i-1] == 0 || doji[i] == 0 {
			continue
		}
		// Bullish: middle doji gaps below the others.
		mid := (open[i-1] + close[i-1]) / 2
		left := (open[i-2] + close[i-2]) / 2
		right := (open[i] + close[i]) / 2
		if mid < left && mid < right {
			out[i] = 100
		} else if mid > left && mid > right {
			out[i] = -100
		}
	}
	return out
}

// CDLUNIQUE3RIVERFn — Bullish reversal: long bear, hammer-like bear inside, small bull below.
func CDLUNIQUE3RIVERFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		ab := avgBody(open, close, i-2)
		if ab <= 0 {
			continue
		}
		if !bearish(o1, c1) || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		if !bearish(o2, c2) || low[i-1] >= low[i-2] {
			continue
		}
		// Bar 2 has long lower shadow and is inside bar 1.
		us2, ls2 := upperShadow(o2, high[i-1], c2), lowerShadow(o2, low[i-1], c2)
		if !(ls2 > 2*body(o2, c2) && us2 <= body(o2, c2)) {
			continue
		}
		if bullish(o3, c3) && body(o3, c3) <= shortBodyRatio*ab && c3 < c2 {
			out[i] = 100
		}
	}
	return out
}

// CDLUPSIDEGAP2CROWSFn — Bearish: bull, gap-up bear, bear engulfing the prior bear.
func CDLUPSIDEGAP2CROWSFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		ab := avgBody(open, close, i-2)
		if ab <= 0 {
			continue
		}
		if !bullish(o1, c1) || body(o1, c1) < longBodyRatio*ab {
			continue
		}
		if !bearish(o2, c2) || math.Min(o2, c2) <= c1 {
			continue
		}
		if !bearish(o3, c3) || o3 <= o2 || c3 >= c2 || c3 <= c1 {
			continue
		}
		out[i] = -100
	}
	return out
}

// CDLTASUKIGAPFn — 3-bar continuation: gap-bar-bar where bar 3 is opposite color and stays in the gap.
func CDLTASUKIGAPFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		// Bullish Tasuki: bar 1 bull, bar 2 bull gapping up, bar 3 bear opening inside bar 2 body, closing inside the gap.
		if bullish(o1, c1) && bullish(o2, c2) && low[i-1] > high[i-2] {
			if bearish(o3, c3) && o3 > o2 && o3 < c2 && c3 < o2 && c3 > c1 {
				out[i] = 100
			}
		}
		if bearish(o1, c1) && bearish(o2, c2) && high[i-1] < low[i-2] {
			if bullish(o3, c3) && o3 < o2 && o3 > c2 && c3 > o2 && c3 < c1 {
				out[i] = -100
			}
		}
	}
	return out
}

// CDLGAPSIDESIDEWHITEFn — Two same-color candles after a gap: bullish/bearish continuation.
func CDLGAPSIDESIDEWHITEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		// Bullish: prior bull, then two bulls of similar size gapping up.
		if bullish(o1, c1) && bullish(o2, c2) && bullish(o3, c3) &&
			low[i-1] > high[i-2] && math.Abs(o3-o2) <= 0.01*o2 &&
			math.Abs(body(o2, c2)-body(o3, c3)) <= 0.3*body(o2, c2) {
			out[i] = 100
		}
		if bearish(o1, c1) && bearish(o2, c2) && bullish(o3, c3) &&
			high[i-1] < low[i-2] && math.Abs(o3-o2) <= 0.01*o2 &&
			math.Abs(body(o2, c2)-body(o3, c3)) <= 0.3*body(o2, c2) {
			// In a downtrend, two side-by-side whites still indicate continuation lower in TA-Lib semantics.
			out[i] = -100
		}
	}
	return out
}

// CDL3LINESTRIKEFn — Four-bar reversal: three same-direction bars then opposite engulfing all.
func CDL3LINESTRIKEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 3; i < n; i++ {
		o1, c1 := open[i-3], close[i-3]
		o2, c2 := open[i-2], close[i-2]
		o3, c3 := open[i-1], close[i-1]
		o4, c4 := open[i], close[i]
		// Bullish strike: 3 bears closing lower, then bull opening below c3 and closing above o1.
		if bearish(o1, c1) && bearish(o2, c2) && bearish(o3, c3) &&
			c2 < c1 && c3 < c2 &&
			bullish(o4, c4) && o4 < c3 && c4 > o1 {
			out[i] = 100
		}
		if bullish(o1, c1) && bullish(o2, c2) && bullish(o3, c3) &&
			c2 > c1 && c3 > c2 &&
			bearish(o4, c4) && o4 > c3 && c4 < o1 {
			out[i] = -100
		}
	}
	return out
}
