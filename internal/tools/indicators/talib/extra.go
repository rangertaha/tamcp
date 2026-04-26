package talib

import "math"

// MAMAFn — MESA Adaptive Moving Average. Returns mama and fama (Following Adaptive MA).
//
// Reuses the Hilbert transform pipeline (htCompute) to estimate the dominant
// cycle phase, then derives an adaptive smoothing factor alpha that varies
// inversely with the rate of phase change. fama follows mama with half its
// alpha, providing a smoother trailing signal.
//
// fastLimit and slowLimit clamp alpha. TA-Lib defaults: fastLimit=0.5, slowLimit=0.05.
func MAMAFn(real []float64, fastLimit, slowLimit float64) (mama, fama []float64) {
	n := len(real)
	mama = make([]float64, n)
	fama = make([]float64, n)
	if n <= htWarmup {
		return
	}
	if fastLimit <= 0 {
		fastLimit = 0.5
	}
	if slowLimit <= 0 {
		slowLimit = 0.05
	}
	st := htCompute(real)
	prevPhase := 0.0
	for i := htWarmup; i < n; i++ {
		var phase float64
		if st.i1[i] != 0 {
			phase = math.Atan(st.q1[i]/st.i1[i]) * 180 / math.Pi
		}
		deltaPhase := prevPhase - phase
		if deltaPhase < 1 {
			deltaPhase = 1
		}
		alpha := fastLimit / deltaPhase
		if alpha < slowLimit {
			alpha = slowLimit
		}
		if alpha > fastLimit {
			alpha = fastLimit
		}
		mama[i] = alpha*real[i] + (1-alpha)*mama[i-1]
		fama[i] = 0.5*alpha*mama[i] + (1-0.5*alpha)*fama[i-1]
		prevPhase = phase
	}
	return
}

// MAVPFn — Moving Average with Variable Period.
// periods[i] specifies the lookback for output[i], clamped to [minPeriod, maxPeriod].
// To avoid recomputing identical-period SMAs, results are cached by period value.
func MAVPFn(real, periods []float64, minPeriod, maxPeriod int, t MaType) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n == 0 || len(periods) != n {
		return out
	}
	if minPeriod < 2 {
		minPeriod = 2
	}
	if maxPeriod < minPeriod {
		maxPeriod = minPeriod
	}
	cache := make(map[int][]float64)
	for i := 0; i < n; i++ {
		p := int(math.Round(periods[i]))
		if p < minPeriod {
			p = minPeriod
		}
		if p > maxPeriod {
			p = maxPeriod
		}
		series, ok := cache[p]
		if !ok {
			series = MA(real, p, t)
			cache[p] = series
		}
		out[i] = series[i]
	}
	return out
}

// SAREXTFn — Extended Parabolic SAR with separate long/short AF parameters.
// Defaults (when zero): startValue=0 (auto-detect), offsetOnReverse=0,
// AF init/step/max long = 0.02/0.02/0.2, same for short.
//
// startValue: positive forces initial long state at that SAR; negative forces short.
// offsetOnReverse: percent (0..1) added to SAR when reversing.
func SAREXTFn(high, low []float64,
	startValue, offsetOnReverse,
	afInitLong, afLong, afMaxLong,
	afInitShort, afShort, afMaxShort float64) []float64 {

	n := len(high)
	out := make([]float64, n)
	if n < 2 {
		return out
	}
	if afInitLong <= 0 {
		afInitLong = 0.02
	}
	if afLong <= 0 {
		afLong = 0.02
	}
	if afMaxLong <= 0 {
		afMaxLong = 0.2
	}
	if afInitShort <= 0 {
		afInitShort = 0.02
	}
	if afShort <= 0 {
		afShort = 0.02
	}
	if afMaxShort <= 0 {
		afMaxShort = 0.2
	}

	long := true
	switch {
	case startValue > 0:
		long = true
	case startValue < 0:
		long = false
	default:
		// Detect initial trend from first +DM vs -DM.
		upMove := high[1] - high[0]
		downMove := low[0] - low[1]
		if downMove > upMove && downMove > 0 {
			long = false
		}
	}

	var sar, ep, af float64
	if long {
		sar = low[0]
		if startValue > 0 {
			sar = startValue
		}
		ep = high[0]
		if high[1] > ep {
			ep = high[1]
		}
		af = afInitLong
	} else {
		sar = high[0]
		if startValue < 0 {
			sar = -startValue
		}
		ep = low[0]
		if low[1] < ep {
			ep = low[1]
		}
		af = afInitShort
	}
	out[1] = sar

	for i := 2; i < n; i++ {
		sar = sar + af*(ep-sar)
		if long {
			if sar > low[i-1] {
				sar = low[i-1]
			}
			if sar > low[i-2] {
				sar = low[i-2]
			}
			if low[i] < sar {
				long = false
				sar = ep + offsetOnReverse*ep
				ep = low[i]
				af = afInitShort
			} else {
				if high[i] > ep {
					ep = high[i]
					af += afLong
					if af > afMaxLong {
						af = afMaxLong
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
				sar = ep - offsetOnReverse*ep
				ep = high[i]
				af = afInitLong
			} else {
				if low[i] < ep {
					ep = low[i]
					af += afShort
					if af > afMaxShort {
						af = afMaxShort
					}
				}
			}
		}
		out[i] = sar
	}
	return out
}

// MACDEXTFn — MACD with separate MA-type selection for fast, slow, and signal lines.
func MACDEXTFn(real []float64, fastPeriod int, fastMA MaType, slowPeriod int, slowMA MaType, signalPeriod int, signalMA MaType) (macd, signal, hist []float64) {
	n := len(real)
	macd = make([]float64, n)
	signal = make([]float64, n)
	hist = make([]float64, n)
	if fastPeriod < 1 || slowPeriod < 1 || signalPeriod < 1 {
		return
	}
	if fastPeriod > slowPeriod {
		fastPeriod, slowPeriod = slowPeriod, fastPeriod
		fastMA, slowMA = slowMA, fastMA
	}
	fast := MA(real, fastPeriod, fastMA)
	slow := MA(real, slowPeriod, slowMA)
	tmp := make([]float64, n-(slowPeriod-1))
	for i := slowPeriod - 1; i < n; i++ {
		tmp[i-(slowPeriod-1)] = fast[i] - slow[i]
	}
	sig := MA(tmp, signalPeriod, signalMA)
	for i := slowPeriod - 1; i < n; i++ {
		macd[i] = tmp[i-(slowPeriod-1)]
	}
	signalStart := slowPeriod - 1 + signalPeriod - 1
	for i := signalStart; i < n; i++ {
		signal[i] = sig[i-(slowPeriod-1)]
		hist[i] = macd[i] - signal[i]
	}
	return
}
