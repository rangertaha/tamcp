package talib

import "math"

// Candlestick pattern detectors. Output convention follows TA-Lib:
//   -100  bearish signal
//      0  no pattern
//   +100  bullish signal
// Some patterns are intrinsically directional (always +100 or always -100 when
// triggered); others can fire either way depending on the bars.

// Threshold helpers — TA-Lib has configurable CandleSettings for these. The
// values below are reasonable defaults and match the look-and-feel of the
// reference C implementation for typical OHLC data.

const (
	// trendLookback bars used to assess prior trend.
	trendLookback = 5
	// avgBodyPeriod bars used to compute the recent average body size.
	avgBodyPeriod = 14
	// dojiBodyRatio: a candle qualifies as a doji if body <= dojiBodyRatio * range.
	dojiBodyRatio = 0.1
	// shortBodyRatio: body considered "short" if <= ratio * avg body.
	shortBodyRatio = 0.5
	// longBodyRatio: body considered "long" if >= ratio * avg body.
	longBodyRatio = 1.5
)

func body(o, c float64) float64        { return math.Abs(c - o) }
func candleRange(h, l float64) float64 { return h - l }
func upperShadow(o, h, c float64) float64 {
	return h - math.Max(o, c)
}
func lowerShadow(o, l, c float64) float64 {
	return math.Min(o, c) - l
}
func bullish(o, c float64) bool { return c > o }
func bearish(o, c float64) bool { return c < o }

// avgBody returns the mean body size of bars [end-avgBodyPeriod, end-1].
// Returns 0 if not enough history.
func avgBody(open, close []float64, end int) float64 {
	start := end - avgBodyPeriod
	if start < 0 {
		return 0
	}
	var s float64
	for i := start; i < end; i++ {
		s += body(open[i], close[i])
	}
	return s / float64(avgBodyPeriod)
}

// downtrend returns true if close[i-1] is below close[i-trendLookback].
func downtrend(close []float64, i int) bool {
	if i < trendLookback {
		return false
	}
	return close[i-1] < close[i-trendLookback]
}

func uptrend(close []float64, i int) bool {
	if i < trendLookback {
		return false
	}
	return close[i-1] > close[i-trendLookback]
}

// CDLDOJIFn — Doji: open ≈ close (body very small relative to range).
func CDLDOJIFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		if rng > 0 && body(open[i], close[i]) <= dojiBodyRatio*rng {
			out[i] = 100
		}
	}
	return out
}

// CDLLONGLEGGEDDOJIFn — Doji with long upper and lower shadows.
func CDLLONGLEGGEDDOJIFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		if rng <= 0 || body(open[i], close[i]) > dojiBodyRatio*rng {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if us > rng*0.3 && ls > rng*0.3 {
			out[i] = 100
		}
	}
	return out
}

// CDLDRAGONFLYDOJIFn — Doji with long lower shadow, no upper shadow.
func CDLDRAGONFLYDOJIFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		if rng <= 0 || body(open[i], close[i]) > dojiBodyRatio*rng {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if us <= rng*0.1 && ls >= rng*0.6 {
			out[i] = 100
		}
	}
	return out
}

// CDLGRAVESTONEDOJIFn — Doji with long upper shadow, no lower shadow.
func CDLGRAVESTONEDOJIFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		if rng <= 0 || body(open[i], close[i]) > dojiBodyRatio*rng {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if ls <= rng*0.1 && us >= rng*0.6 {
			out[i] = -100
		}
	}
	return out
}

// CDLHAMMERFn — Small body, long lower shadow, little upper shadow, in a downtrend.
func CDLHAMMERFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if b <= rng*0.3 && ls >= 2*b && us <= b && downtrend(close, i) {
			out[i] = 100
		}
	}
	return out
}

// CDLHANGINGMANFn — Same geometry as Hammer but appears in an uptrend (bearish).
func CDLHANGINGMANFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if b <= rng*0.3 && ls >= 2*b && us <= b && uptrend(close, i) {
			out[i] = -100
		}
	}
	return out
}

// CDLINVERTEDHAMMERFn — Small body, long upper shadow, little lower shadow, in a downtrend (bullish).
func CDLINVERTEDHAMMERFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if b <= rng*0.3 && us >= 2*b && ls <= b && downtrend(close, i) {
			out[i] = 100
		}
	}
	return out
}

// CDLSHOOTINGSTARFn — Same geometry as Inverted Hammer but in an uptrend (bearish).
func CDLSHOOTINGSTARFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if b <= rng*0.3 && us >= 2*b && ls <= b && uptrend(close, i) {
			out[i] = -100
		}
	}
	return out
}

// CDLMARUBOZUFn — Long body with very small or no shadows. Sign matches body direction.
func CDLMARUBOZUFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 {
			continue
		}
		ab := avgBody(open, close, i)
		if ab > 0 && b >= longBodyRatio*ab && b >= 0.95*rng {
			if bullish(open[i], close[i]) {
				out[i] = 100
			} else if bearish(open[i], close[i]) {
				out[i] = -100
			}
		}
	}
	return out
}

// CDLLONGLINEFn — Body at least longBodyRatio × average. Sign matches body direction.
func CDLLONGLINEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		ab := avgBody(open, close, i)
		if ab <= 0 {
			continue
		}
		if body(open[i], close[i]) >= longBodyRatio*ab {
			if bullish(open[i], close[i]) {
				out[i] = 100
			} else if bearish(open[i], close[i]) {
				out[i] = -100
			}
		}
	}
	return out
}

// CDLSHORTLINEFn — Body at most shortBodyRatio × average. Sign matches body direction.
func CDLSHORTLINEFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		ab := avgBody(open, close, i)
		if ab <= 0 {
			continue
		}
		if body(open[i], close[i]) <= shortBodyRatio*ab {
			if bullish(open[i], close[i]) {
				out[i] = 100
			} else if bearish(open[i], close[i]) {
				out[i] = -100
			}
		}
	}
	return out
}

// CDLSPINNINGTOPFn — Small body with substantial shadows on both sides.
func CDLSPINNINGTOPFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 0; i < n; i++ {
		rng := candleRange(high[i], low[i])
		b := body(open[i], close[i])
		if rng <= 0 || b > rng*0.4 {
			continue
		}
		us, ls := upperShadow(open[i], high[i], close[i]), lowerShadow(open[i], low[i], close[i])
		if us > b && ls > b {
			if bullish(open[i], close[i]) {
				out[i] = 100
			} else if bearish(open[i], close[i]) {
				out[i] = -100
			}
		}
	}
	return out
}

// CDLENGULFINGFn — Two-bar pattern: bar 2 body fully engulfs bar 1's opposite-color body.
func CDLENGULFINGFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		// Bullish engulfing: prev bearish, curr bullish, curr engulfs prev body.
		if bearish(o1, c1) && bullish(o2, c2) && c2 >= o1 && o2 <= c1 {
			out[i] = 100
		} else if bullish(o1, c1) && bearish(o2, c2) && o2 >= c1 && c2 <= o1 {
			out[i] = -100
		}
	}
	return out
}

// CDLHARAMIFn — Two-bar pattern: bar 2 body fully inside bar 1's opposite-color body.
func CDLHARAMIFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		hi1, lo1 := math.Max(o1, c1), math.Min(o1, c1)
		hi2, lo2 := math.Max(o2, c2), math.Min(o2, c2)
		if hi2 > hi1 || lo2 < lo1 {
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

// CDLHARAMICROSSFn — Harami where the second bar is a doji.
func CDLHARAMICROSSFn(open, high, low, close []float64) []int {
	har := CDLHARAMIFn(open, high, low, close)
	doji := CDLDOJIFn(open, high, low, close)
	out := make([]int, len(har))
	for i := range out {
		if har[i] != 0 && doji[i] != 0 {
			out[i] = har[i]
		}
	}
	return out
}

// CDLMORNINGSTARFn — Three-bar bullish reversal: long bear, small body (gap down),
// long bull closing above midpoint of bar 1.
func CDLMORNINGSTARFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		c2 := close[i-1]
		o2 := open[i-1]
		o3, c3 := open[i], close[i]
		ab := avgBody(open, close, i-2)
		if ab <= 0 {
			continue
		}
		mid1 := (o1 + c1) / 2
		// Bar 1: long bear; bar 2: small body, gaps down from bar 1; bar 3: long bull, closes above bar 1 midpoint.
		if bearish(o1, c1) && body(o1, c1) >= longBodyRatio*ab &&
			body(o2, c2) <= shortBodyRatio*ab &&
			math.Max(o2, c2) < c1 &&
			bullish(o3, c3) && body(o3, c3) >= longBodyRatio*ab &&
			c3 > mid1 {
			out[i] = 100
		}
	}
	return out
}

// CDLEVENINGSTARFn — Three-bar bearish reversal mirror of Morning Star.
func CDLEVENINGSTARFn(open, high, low, close []float64) []int {
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
		mid1 := (o1 + c1) / 2
		if bullish(o1, c1) && body(o1, c1) >= longBodyRatio*ab &&
			body(o2, c2) <= shortBodyRatio*ab &&
			math.Min(o2, c2) > c1 &&
			bearish(o3, c3) && body(o3, c3) >= longBodyRatio*ab &&
			c3 < mid1 {
			out[i] = -100
		}
	}
	return out
}

// CDL3WHITESOLDIERSFn — Three consecutive bullish bars, each opening within
// the prior body and closing higher.
func CDL3WHITESOLDIERSFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 2; i < n; i++ {
		o1, c1 := open[i-2], close[i-2]
		o2, c2 := open[i-1], close[i-1]
		o3, c3 := open[i], close[i]
		if !bullish(o1, c1) || !bullish(o2, c2) || !bullish(o3, c3) {
			continue
		}
		if c2 <= c1 || c3 <= c2 {
			continue
		}
		if !(o2 > o1 && o2 < c1) || !(o3 > o2 && o3 < c2) {
			continue
		}
		out[i] = 100
	}
	return out
}

// CDL3BLACKCROWSFn — Three consecutive bearish bars, each opening within the
// prior body and closing lower.
func CDL3BLACKCROWSFn(open, high, low, close []float64) []int {
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
		if !(o2 < o1 && o2 > c1) || !(o3 < o2 && o3 > c2) {
			continue
		}
		out[i] = -100
	}
	return out
}

// CDLPIERCINGFn — Two-bar bullish reversal: long bear; bar 2 opens below bar 1 low,
// closes above the midpoint of bar 1's body.
func CDLPIERCINGFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		ab := avgBody(open, close, i-1)
		if ab <= 0 {
			continue
		}
		if bearish(o1, c1) && body(o1, c1) >= longBodyRatio*ab &&
			bullish(o2, c2) && o2 < low[i-1] &&
			c2 > (o1+c1)/2 && c2 < o1 {
			out[i] = 100
		}
	}
	return out
}

// CDLDARKCLOUDCOVERFn — Two-bar bearish reversal mirror of Piercing.
func CDLDARKCLOUDCOVERFn(open, high, low, close []float64) []int {
	n := len(open)
	out := make([]int, n)
	for i := 1; i < n; i++ {
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		ab := avgBody(open, close, i-1)
		if ab <= 0 {
			continue
		}
		if bullish(o1, c1) && body(o1, c1) >= longBodyRatio*ab &&
			bearish(o2, c2) && o2 > high[i-1] &&
			c2 < (o1+c1)/2 && c2 > o1 {
			out[i] = -100
		}
	}
	return out
}
