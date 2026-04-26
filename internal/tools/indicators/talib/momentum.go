package talib

import "math"

// MOMFn — Momentum: real[i] - real[i-period].
func MOMFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	for i := period; i < n; i++ {
		out[i] = real[i] - real[i-period]
	}
	return out
}

// ROCFn — Rate of Change: ((real[i]/real[i-period]) - 1) * 100.
func ROCFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	for i := period; i < n; i++ {
		if real[i-period] != 0 {
			out[i] = (real[i]/real[i-period] - 1) * 100
		}
	}
	return out
}

// ROCPFn — Rate of Change Percentage: (real[i] - real[i-period]) / real[i-period].
func ROCPFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	for i := period; i < n; i++ {
		if real[i-period] != 0 {
			out[i] = (real[i] - real[i-period]) / real[i-period]
		}
	}
	return out
}

// ROCRFn — Rate of Change Ratio: real[i] / real[i-period].
func ROCRFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	for i := period; i < n; i++ {
		if real[i-period] != 0 {
			out[i] = real[i] / real[i-period]
		}
	}
	return out
}

// ROCR100Fn — Rate of Change Ratio (×100): (real[i] / real[i-period]) * 100.
func ROCR100Fn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	for i := period; i < n; i++ {
		if real[i-period] != 0 {
			out[i] = real[i] / real[i-period] * 100
		}
	}
	return out
}

// RSIFn — Wilder's Relative Strength Index.
func RSIFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		ch := real[i] - real[i-1]
		if ch > 0 {
			avgGain += ch
		} else {
			avgLoss -= ch
		}
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)
	out[period] = rsiValue(avgGain, avgLoss)
	pf := float64(period)
	for i := period + 1; i < n; i++ {
		ch := real[i] - real[i-1]
		gain, loss := 0.0, 0.0
		if ch > 0 {
			gain = ch
		} else {
			loss = -ch
		}
		avgGain = (avgGain*(pf-1) + gain) / pf
		avgLoss = (avgLoss*(pf-1) + loss) / pf
		out[i] = rsiValue(avgGain, avgLoss)
	}
	return out
}

func rsiValue(gain, loss float64) float64 {
	if loss == 0 {
		if gain == 0 {
			return 0
		}
		return 100
	}
	rs := gain / loss
	return 100 - 100/(1+rs)
}

// CMOFn — Chande Momentum Oscillator: 100 * (sumGain - sumLoss) / (sumGain + sumLoss).
func CMOFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	var gain, loss float64
	for i := 1; i <= period; i++ {
		ch := real[i] - real[i-1]
		if ch > 0 {
			gain += ch
		} else {
			loss -= ch
		}
	}
	if gain+loss > 0 {
		out[period] = 100 * (gain - loss) / (gain + loss)
	}
	pf := float64(period)
	// TA-Lib's CMO uses Wilder-smoothed gain/loss after the seed window.
	for i := period + 1; i < n; i++ {
		ch := real[i] - real[i-1]
		g, l := 0.0, 0.0
		if ch > 0 {
			g = ch
		} else {
			l = -ch
		}
		gain = (gain*(pf-1) + g) / pf
		loss = (loss*(pf-1) + l) / pf
		if gain+loss > 0 {
			out[i] = 100 * (gain - loss) / (gain + loss)
		}
	}
	return out
}

// BOPFn — Balance of Power: (close - open) / (high - low).
func BOPFn(open, high, low, close []float64) []float64 {
	n := len(open)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		hl := high[i] - low[i]
		if hl > 0 {
			out[i] = (close[i] - open[i]) / hl
		}
	}
	return out
}

// WILLRFn — Williams' %R: -100 * (highestHigh - close) / (highestHigh - lowestLow).
func WILLRFn(high, low, close []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		hh, ll := high[i-period+1], low[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if high[j] > hh {
				hh = high[j]
			}
			if low[j] < ll {
				ll = low[j]
			}
		}
		if hh != ll {
			out[i] = -100 * (hh - close[i]) / (hh - ll)
		}
	}
	return out
}

// CCIFn — Commodity Channel Index: (TP - SMA(TP, period)) / (0.015 * meanDev).
func CCIFn(high, low, close []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	tp := TYPPRICEFn(high, low, close)
	for i := period - 1; i < n; i++ {
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			sum += tp[j]
		}
		mean := sum / float64(period)
		var dev float64
		for j := i - period + 1; j <= i; j++ {
			dev += math.Abs(tp[j] - mean)
		}
		md := dev / float64(period)
		if md > 0 {
			out[i] = (tp[i] - mean) / (0.015 * md)
		}
	}
	return out
}

// MACDFn — MACD: EMA(real, fast) - EMA(real, slow); signal = EMA(macd, signalPeriod); hist = macd - signal.
func MACDFn(real []float64, fastPeriod, slowPeriod, signalPeriod int) (macd, signal, hist []float64) {
	n := len(real)
	macd = make([]float64, n)
	signal = make([]float64, n)
	hist = make([]float64, n)
	if fastPeriod < 1 || slowPeriod < 1 || signalPeriod < 1 {
		return
	}
	if fastPeriod > slowPeriod {
		fastPeriod, slowPeriod = slowPeriod, fastPeriod
	}
	fast := EMAFn(real, fastPeriod)
	slow := EMAFn(real, slowPeriod)
	// MACD valid from slowPeriod-1 onward.
	tmp := make([]float64, n-(slowPeriod-1))
	for i := slowPeriod - 1; i < n; i++ {
		tmp[i-(slowPeriod-1)] = fast[i] - slow[i]
	}
	sig := EMAFn(tmp, signalPeriod)
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

// MACDFIXFn — MACD with fixed 12/26 fast/slow.
func MACDFIXFn(real []float64, signalPeriod int) (macd, signal, hist []float64) {
	return MACDFn(real, 12, 26, signalPeriod)
}

// APOFn — Absolute Price Oscillator: MA(real, fast, type) - MA(real, slow, type).
func APOFn(real []float64, fastPeriod, slowPeriod int, t MaType) []float64 {
	if fastPeriod > slowPeriod {
		fastPeriod, slowPeriod = slowPeriod, fastPeriod
	}
	fast := MA(real, fastPeriod, t)
	slow := MA(real, slowPeriod, t)
	n := len(real)
	out := make([]float64, n)
	for i := slowPeriod - 1; i < n; i++ {
		out[i] = fast[i] - slow[i]
	}
	return out
}

// PPOFn — Percentage Price Oscillator: 100 * (fast - slow) / slow.
func PPOFn(real []float64, fastPeriod, slowPeriod int, t MaType) []float64 {
	if fastPeriod > slowPeriod {
		fastPeriod, slowPeriod = slowPeriod, fastPeriod
	}
	fast := MA(real, fastPeriod, t)
	slow := MA(real, slowPeriod, t)
	n := len(real)
	out := make([]float64, n)
	for i := slowPeriod - 1; i < n; i++ {
		if slow[i] != 0 {
			out[i] = 100 * (fast[i] - slow[i]) / slow[i]
		}
	}
	return out
}

// AROONFn — Aroon: down=100*(period-sinceLowest)/period; up=100*(period-sinceHighest)/period.
func AROONFn(high, low []float64, period int) (down, up []float64) {
	n := len(high)
	down = make([]float64, n)
	up = make([]float64, n)
	if period < 1 || n <= period {
		return
	}
	for i := period; i < n; i++ {
		hi, lo := high[i-period], low[i-period]
		hiIdx, loIdx := i-period, i-period
		for j := i - period + 1; j <= i; j++ {
			if high[j] >= hi {
				hi = high[j]
				hiIdx = j
			}
			if low[j] <= lo {
				lo = low[j]
				loIdx = j
			}
		}
		up[i] = 100 * float64(period-(i-hiIdx)) / float64(period)
		down[i] = 100 * float64(period-(i-loIdx)) / float64(period)
	}
	return
}

// AROONOSCFn — Aroon Oscillator: up - down.
func AROONOSCFn(high, low []float64, period int) []float64 {
	d, u := AROONFn(high, low, period)
	out := make([]float64, len(d))
	for i := range out {
		out[i] = u[i] - d[i]
	}
	return out
}

// dmTrSeed seeds Wilder smoothed +DM, -DM, and TR over the first `period` bars.
// Returns smoothed values and the per-bar TR series for follow-on smoothing.
func dmTrSmooth(high, low, close []float64, period int) (plusDM, minusDM, tr []float64) {
	n := len(high)
	plusDM = make([]float64, n)
	minusDM = make([]float64, n)
	tr = make([]float64, n)
	if period < 1 || n <= period {
		return
	}
	tr0 := trueRange(high, low, close)
	var sumPlus, sumMinus, sumTR float64
	for i := 1; i <= period; i++ {
		up := high[i] - high[i-1]
		dn := low[i-1] - low[i]
		if up > dn && up > 0 {
			sumPlus += up
		}
		if dn > up && dn > 0 {
			sumMinus += dn
		}
		sumTR += tr0[i]
	}
	plusDM[period] = sumPlus
	minusDM[period] = sumMinus
	tr[period] = sumTR
	pf := float64(period)
	for i := period + 1; i < n; i++ {
		up := high[i] - high[i-1]
		dn := low[i-1] - low[i]
		pd, md := 0.0, 0.0
		if up > dn && up > 0 {
			pd = up
		}
		if dn > up && dn > 0 {
			md = dn
		}
		plusDM[i] = plusDM[i-1] - plusDM[i-1]/pf + pd
		minusDM[i] = minusDM[i-1] - minusDM[i-1]/pf + md
		tr[i] = tr[i-1] - tr[i-1]/pf + tr0[i]
	}
	return
}

// PLUSDMFn — Plus Directional Movement (Wilder smoothed).
func PLUSDMFn(high, low []float64, period int) []float64 {
	// Reuse dmTrSmooth with a dummy close (TR not needed here).
	close := make([]float64, len(high))
	for i := range close {
		close[i] = (high[i] + low[i]) / 2
	}
	p, _, _ := dmTrSmooth(high, low, close, period)
	return p
}

// MINUSDMFn — Minus Directional Movement (Wilder smoothed).
func MINUSDMFn(high, low []float64, period int) []float64 {
	close := make([]float64, len(high))
	for i := range close {
		close[i] = (high[i] + low[i]) / 2
	}
	_, m, _ := dmTrSmooth(high, low, close, period)
	return m
}

// PLUSDIFn — Plus Directional Indicator: 100 * +DM / TR.
func PLUSDIFn(high, low, close []float64, period int) []float64 {
	pd, _, tr := dmTrSmooth(high, low, close, period)
	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		if tr[i] > 0 {
			out[i] = 100 * pd[i] / tr[i]
		}
	}
	return out
}

// MINUSDIFn — Minus Directional Indicator: 100 * -DM / TR.
func MINUSDIFn(high, low, close []float64, period int) []float64 {
	_, md, tr := dmTrSmooth(high, low, close, period)
	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		if tr[i] > 0 {
			out[i] = 100 * md[i] / tr[i]
		}
	}
	return out
}

// DXFn — Directional Movement Index: 100 * |+DI - -DI| / (+DI + -DI).
func DXFn(high, low, close []float64, period int) []float64 {
	pdi := PLUSDIFn(high, low, close, period)
	mdi := MINUSDIFn(high, low, close, period)
	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		s := pdi[i] + mdi[i]
		if s > 0 {
			out[i] = 100 * math.Abs(pdi[i]-mdi[i]) / s
		}
	}
	return out
}

// ADXFn — Average Directional Movement Index: Wilder-smoothed DX.
func ADXFn(high, low, close []float64, period int) []float64 {
	dx := DXFn(high, low, close, period)
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n < 2*period {
		return out
	}
	// Seed at index 2*period-1 = avg of dx[period..2*period-1].
	var sum float64
	for i := period; i < 2*period; i++ {
		sum += dx[i]
	}
	prev := sum / float64(period)
	out[2*period-1] = prev
	pf := float64(period)
	for i := 2 * period; i < n; i++ {
		prev = (prev*(pf-1) + dx[i]) / pf
		out[i] = prev
	}
	return out
}

// ADXRFn — ADX Rating: (ADX[i] + ADX[i-period]) / 2.
func ADXRFn(high, low, close []float64, period int) []float64 {
	adx := ADXFn(high, low, close, period)
	n := len(high)
	out := make([]float64, n)
	for i := 2*period - 1 + period - 1; i < n; i++ {
		out[i] = (adx[i] + adx[i-period+1]) / 2
	}
	return out
}

// MFIFn — Money Flow Index.
func MFIFn(high, low, close, volume []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	tp := TYPPRICEFn(high, low, close)
	mf := make([]float64, n)
	for i := 0; i < n; i++ {
		mf[i] = tp[i] * volume[i]
	}
	for i := period; i < n; i++ {
		var pos, neg float64
		for j := i - period + 1; j <= i; j++ {
			if tp[j] > tp[j-1] {
				pos += mf[j]
			} else if tp[j] < tp[j-1] {
				neg += mf[j]
			}
		}
		if neg == 0 {
			out[i] = 100
		} else {
			ratio := pos / neg
			out[i] = 100 - 100/(1+ratio)
		}
	}
	return out
}

// TRIXFn — 1-day ROC of triple-smoothed EMA.
func TRIXFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < 3*period {
		return out
	}
	e1 := EMAFn(real, period)
	e2 := EMAFn(e1[period-1:], period)
	e3 := EMAFn(e2[period-1:], period)
	// e3 is aligned starting at original index 3*(period-1).
	start := 3 * (period - 1)
	for i := start + 1; i < n; i++ {
		prev := e3[i-1-start]
		curr := e3[i-start]
		if prev != 0 {
			out[i] = 100 * (curr - prev) / prev
		}
	}
	return out
}

// STOCHFn — Stochastic Oscillator: slow %K and slow %D.
func STOCHFn(high, low, close []float64, fastKPeriod, slowKPeriod int, slowKMA MaType, slowDPeriod int, slowDMA MaType) (slowK, slowD []float64) {
	n := len(high)
	slowK = make([]float64, n)
	slowD = make([]float64, n)
	if fastKPeriod < 1 || slowKPeriod < 1 || slowDPeriod < 1 || n < fastKPeriod {
		return
	}
	fastK := make([]float64, n)
	for i := fastKPeriod - 1; i < n; i++ {
		hh, ll := high[i-fastKPeriod+1], low[i-fastKPeriod+1]
		for j := i - fastKPeriod + 2; j <= i; j++ {
			if high[j] > hh {
				hh = high[j]
			}
			if low[j] < ll {
				ll = low[j]
			}
		}
		if hh != ll {
			fastK[i] = 100 * (close[i] - ll) / (hh - ll)
		}
	}
	slowK = MA(fastK, slowKPeriod, slowKMA)
	slowD = MA(slowK, slowDPeriod, slowDMA)
	return
}

// STOCHFFn — Stochastic Fast: fast %K and fast %D.
func STOCHFFn(high, low, close []float64, fastKPeriod, fastDPeriod int, fastDMA MaType) (fastK, fastD []float64) {
	n := len(high)
	fastK = make([]float64, n)
	fastD = make([]float64, n)
	if fastKPeriod < 1 || fastDPeriod < 1 || n < fastKPeriod {
		return
	}
	for i := fastKPeriod - 1; i < n; i++ {
		hh, ll := high[i-fastKPeriod+1], low[i-fastKPeriod+1]
		for j := i - fastKPeriod + 2; j <= i; j++ {
			if high[j] > hh {
				hh = high[j]
			}
			if low[j] < ll {
				ll = low[j]
			}
		}
		if hh != ll {
			fastK[i] = 100 * (close[i] - ll) / (hh - ll)
		}
	}
	fastD = MA(fastK, fastDPeriod, fastDMA)
	return
}

// STOCHRSIFn — Stochastic RSI: STOCHF applied to RSI series.
func STOCHRSIFn(real []float64, period, fastKPeriod, fastDPeriod int, fastDMA MaType) (fastK, fastD []float64) {
	rsi := RSIFn(real, period)
	// Use RSI as both high/low/close for the stoch transform.
	return STOCHFFn(rsi, rsi, rsi, fastKPeriod, fastDPeriod, fastDMA)
}

// ULTOSCFn — Ultimate Oscillator across three periods.
func ULTOSCFn(high, low, close []float64, p1, p2, p3 int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if n < 2 {
		return out
	}
	bp := make([]float64, n)
	tr := make([]float64, n)
	for i := 1; i < n; i++ {
		tl := math.Min(low[i], close[i-1])
		th := math.Max(high[i], close[i-1])
		bp[i] = close[i] - tl
		tr[i] = th - tl
	}
	maxP := p1
	if p2 > maxP {
		maxP = p2
	}
	if p3 > maxP {
		maxP = p3
	}
	if n <= maxP {
		return out
	}
	for i := maxP; i < n; i++ {
		bp1, tr1 := windowSum(bp, i, p1), windowSum(tr, i, p1)
		bp2, tr2 := windowSum(bp, i, p2), windowSum(tr, i, p2)
		bp3, tr3 := windowSum(bp, i, p3), windowSum(tr, i, p3)
		if tr1 == 0 || tr2 == 0 || tr3 == 0 {
			continue
		}
		out[i] = 100 * (4*bp1/tr1 + 2*bp2/tr2 + bp3/tr3) / 7
	}
	return out
}

func windowSum(s []float64, end, period int) float64 {
	var sum float64
	for j := end - period + 1; j <= end; j++ {
		sum += s[j]
	}
	return sum
}
