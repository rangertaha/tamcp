package talib

// Community indicators sourced from Pandas TA, sdcoffey/techan,
// cinar/indicator, Ta-Lib-Rust, and Yatala — i.e. things outside
// the original TA-Lib C library that are widely used in practice.

import "math"

// SUPERTRENDFn — ATR-based trend follower.
// Returns the active band as `trend` (lower band when uptrend, upper band
// when downtrend) and `direction` ∈ {+1,-1} (0 before warm-up).
func SUPERTRENDFn(high, low, close []float64, period int, multiplier float64) (trend, direction []float64) {
	n := len(close)
	trend = make([]float64, n)
	direction = make([]float64, n)
	if period < 1 || n <= period {
		return
	}
	atr := ATRFn(high, low, close, period)
	upper := make([]float64, n)
	lower := make([]float64, n)
	dir := 1
	for i := period; i < n; i++ {
		hl2 := (high[i] + low[i]) / 2
		bUp := hl2 + multiplier*atr[i]
		bDn := hl2 - multiplier*atr[i]
		if i == period {
			upper[i] = bUp
			lower[i] = bDn
		} else {
			if bUp < upper[i-1] || close[i-1] > upper[i-1] {
				upper[i] = bUp
			} else {
				upper[i] = upper[i-1]
			}
			if bDn > lower[i-1] || close[i-1] < lower[i-1] {
				lower[i] = bDn
			} else {
				lower[i] = lower[i-1]
			}
			if dir > 0 && close[i] < lower[i-1] {
				dir = -1
			} else if dir < 0 && close[i] > upper[i-1] {
				dir = 1
			}
		}
		direction[i] = float64(dir)
		if dir > 0 {
			trend[i] = lower[i]
		} else {
			trend[i] = upper[i]
		}
	}
	return
}

// VWAPFn — cumulative Volume Weighted Average Price using HLC3 typical price.
// No session reset; pre-slice inputs per session if needed.
func VWAPFn(high, low, close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	var pv, vv float64
	for i := 0; i < n; i++ {
		tp := (high[i] + low[i] + close[i]) / 3
		pv += tp * volume[i]
		vv += volume[i]
		if vv != 0 {
			out[i] = pv / vv
		}
	}
	return out
}

// VWMAFn — Volume Weighted Moving Average:
// sum(real*vol, period) / sum(vol, period).
func VWMAFn(real, volume []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	var pv, vv float64
	for i := 0; i < period; i++ {
		pv += real[i] * volume[i]
		vv += volume[i]
	}
	if vv != 0 {
		out[period-1] = pv / vv
	}
	for i := period; i < n; i++ {
		pv += real[i]*volume[i] - real[i-period]*volume[i-period]
		vv += volume[i] - volume[i-period]
		if vv != 0 {
			out[i] = pv / vv
		}
	}
	return out
}

// HMAFn — Hull Moving Average: WMA(2*WMA(real, p/2) - WMA(real, p), √p).
func HMAFn(real []float64, period int) []float64 {
	n := len(real)
	if period < 2 || n < period {
		return make([]float64, n)
	}
	half := period / 2
	if half < 1 {
		half = 1
	}
	sq := int(math.Round(math.Sqrt(float64(period))))
	if sq < 1 {
		sq = 1
	}
	w1 := WMAFn(real, half)
	w2 := WMAFn(real, period)
	diff := make([]float64, n)
	for i := 0; i < n; i++ {
		diff[i] = 2*w1[i] - w2[i]
	}
	return WMAFn(diff, sq)
}

// ZLEMAFn — Zero-Lag EMA. Pre-de-lags input by ⌊(p-1)/2⌋ samples then EMAs.
func ZLEMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	lag := (period - 1) / 2
	src := make([]float64, n)
	for i := 0; i < n; i++ {
		if i < lag {
			src[i] = real[i]
		} else {
			src[i] = 2*real[i] - real[i-lag]
		}
	}
	return EMAFn(src, period)
}

// SMMAFn — Wilder smoothing / Smoothed MA / RMA:
// (prev*(p-1)+curr)/p, seeded with the SMA of the first p bars.
func SMMAFn(real []float64, period int) []float64 {
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
	pf := float64(period)
	for i := period; i < n; i++ {
		prev = (prev*(pf-1) + real[i]) / pf
		out[i] = prev
	}
	return out
}

// KDJFn — Stochastic-derived KDJ. Returns K, D, J series.
//
//	RSV  = (close - LL_p) / (HH_p - LL_p) * 100
//	K[i] = 2/3 * K[i-1] + 1/3 * RSV[i]   (init 50)
//	D[i] = 2/3 * D[i-1] + 1/3 * K[i]      (init 50)
//	J    = 3*K - 2*D
func KDJFn(high, low, close []float64, period int) (k, d, j []float64) {
	n := len(close)
	k = make([]float64, n)
	d = make([]float64, n)
	j = make([]float64, n)
	if period < 1 || n < period {
		return
	}
	prevK, prevD := 50.0, 50.0
	for i := period - 1; i < n; i++ {
		hh, ll := high[i-period+1], low[i-period+1]
		for jj := i - period + 2; jj <= i; jj++ {
			if high[jj] > hh {
				hh = high[jj]
			}
			if low[jj] < ll {
				ll = low[jj]
			}
		}
		rsv := 0.0
		if hh != ll {
			rsv = (close[i] - ll) / (hh - ll) * 100
		}
		curK := (2.0/3.0)*prevK + (1.0/3.0)*rsv
		curD := (2.0/3.0)*prevD + (1.0/3.0)*curK
		k[i] = curK
		d[i] = curD
		j[i] = 3*curK - 2*curD
		prevK, prevD = curK, curD
	}
	return
}

// AOFn — Bill Williams' Awesome Oscillator: SMA((H+L)/2,5) − SMA((H+L)/2,34).
func AOFn(high, low []float64) []float64 {
	n := len(high)
	median := make([]float64, n)
	for i := 0; i < n; i++ {
		median[i] = (high[i] + low[i]) / 2
	}
	s5 := SMAFn(median, 5)
	s34 := SMAFn(median, 34)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = s5[i] - s34[i]
	}
	return out
}

// CMFFn — Chaikin Money Flow over `period` bars.
//
//	mfm = ((c-l) - (h-c)) / (h-l)
//	mfv = mfm * volume
//	cmf = SUM(mfv, period) / SUM(volume, period)
func CMFFn(high, low, close, volume []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	mfv := make([]float64, n)
	for i := 0; i < n; i++ {
		rng := high[i] - low[i]
		if rng == 0 {
			continue
		}
		mfm := ((close[i] - low[i]) - (high[i] - close[i])) / rng
		mfv[i] = mfm * volume[i]
	}
	var sumMfv, sumVol float64
	for i := 0; i < period; i++ {
		sumMfv += mfv[i]
		sumVol += volume[i]
	}
	if sumVol != 0 {
		out[period-1] = sumMfv / sumVol
	}
	for i := period; i < n; i++ {
		sumMfv += mfv[i] - mfv[i-period]
		sumVol += volume[i] - volume[i-period]
		if sumVol != 0 {
			out[i] = sumMfv / sumVol
		}
	}
	return out
}

// DONCHIANFn — Donchian Channels: highest(high,p), lowest(low,p), midpoint.
func DONCHIANFn(high, low []float64, period int) (upper, middle, lower []float64) {
	n := len(high)
	upper = make([]float64, n)
	lower = make([]float64, n)
	middle = make([]float64, n)
	if period < 1 || n < period {
		return
	}
	for i := period - 1; i < n; i++ {
		hh, ll := high[i-period+1], low[i-period+1]
		for jj := i - period + 2; jj <= i; jj++ {
			if high[jj] > hh {
				hh = high[jj]
			}
			if low[jj] < ll {
				ll = low[jj]
			}
		}
		upper[i] = hh
		lower[i] = ll
		middle[i] = (hh + ll) / 2
	}
	return
}

// FISHERFn — Ehlers Fisher Transform of price. Returns `fisher` and a
// one-bar lagged `signal` series.
//
//	x = 0.66 * (2*((hl2-LL)/(HH-LL)) - 1) + 0.67*x_prev,  clamp |x|<0.999
//	F = 0.5 * ln((1+x)/(1-x)) + 0.5 * F_prev
func FISHERFn(high, low []float64, period int) (fisher, signal []float64) {
	n := len(high)
	fisher = make([]float64, n)
	signal = make([]float64, n)
	if period < 1 || n < period {
		return
	}
	x, fprev := 0.0, 0.0
	for i := period - 1; i < n; i++ {
		hh, ll := high[i-period+1], low[i-period+1]
		for jj := i - period + 2; jj <= i; jj++ {
			if high[jj] > hh {
				hh = high[jj]
			}
			if low[jj] < ll {
				ll = low[jj]
			}
		}
		hl2 := (high[i] + low[i]) / 2
		var raw float64
		if hh != ll {
			raw = 2*((hl2-ll)/(hh-ll)) - 1
		}
		x = 0.66*raw + 0.67*x
		if x > 0.999 {
			x = 0.999
		} else if x < -0.999 {
			x = -0.999
		}
		f := 0.5*math.Log((1+x)/(1-x)) + 0.5*fprev
		if i > period-1 {
			signal[i] = fisher[i-1]
		}
		fisher[i] = f
		fprev = f
	}
	return
}

// TSIFn — True Strength Index. Returns the TSI series and an EMA-smoothed signal.
//
//	m   = close[t] - close[t-1]
//	TSI = 100 * EMA(EMA(m,  r), s) / EMA(EMA(|m|, r), s)
//	sig = EMA(TSI, signal)
func TSIFn(real []float64, r, s, signal int) (tsi, sig []float64) {
	n := len(real)
	tsi = make([]float64, n)
	sig = make([]float64, n)
	if n < 2 || r < 1 || s < 1 {
		return
	}
	mom := make([]float64, n)
	abs := make([]float64, n)
	for i := 1; i < n; i++ {
		d := real[i] - real[i-1]
		mom[i] = d
		if d < 0 {
			abs[i] = -d
		} else {
			abs[i] = d
		}
	}
	num := EMAFn(EMAFn(mom, r), s)
	den := EMAFn(EMAFn(abs, r), s)
	for i := 0; i < n; i++ {
		if den[i] != 0 {
			tsi[i] = 100 * num[i] / den[i]
		}
	}
	if signal > 0 {
		sig = EMAFn(tsi, signal)
	}
	return
}

// KSTFn — Know Sure Thing oscillator.
//
//	KST = 1*SMA(ROC(c,r1),n1) + 2*SMA(ROC(c,r2),n2)
//	    + 3*SMA(ROC(c,r3),n3) + 4*SMA(ROC(c,r4),n4)
//	sig = SMA(KST, sigPeriod)
func KSTFn(real []float64, r1, r2, r3, r4, n1, n2, n3, n4, sigPeriod int) (kst, sig []float64) {
	n := len(real)
	kst = make([]float64, n)
	sig = make([]float64, n)
	rc1 := SMAFn(ROCFn(real, r1), n1)
	rc2 := SMAFn(ROCFn(real, r2), n2)
	rc3 := SMAFn(ROCFn(real, r3), n3)
	rc4 := SMAFn(ROCFn(real, r4), n4)
	for i := 0; i < n; i++ {
		kst[i] = rc1[i] + 2*rc2[i] + 3*rc3[i] + 4*rc4[i]
	}
	if sigPeriod > 0 {
		sig = SMAFn(kst, sigPeriod)
	}
	return
}

// COPPOCKFn — Coppock Curve: WMA(ROC(c,long) + ROC(c,short), wmaPeriod).
func COPPOCKFn(real []float64, longPeriod, shortPeriod, wmaPeriod int) []float64 {
	n := len(real)
	r1 := ROCFn(real, longPeriod)
	r2 := ROCFn(real, shortPeriod)
	sum := make([]float64, n)
	for i := 0; i < n; i++ {
		sum[i] = r1[i] + r2[i]
	}
	return WMAFn(sum, wmaPeriod)
}

// CHOPFn — Choppiness Index over `period` bars.
//
//	CHOP = 100 * log10( SUM(TR, p) / (HHV(high,p) - LLV(low,p)) ) / log10(p)
func CHOPFn(high, low, close []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 2 || n <= period {
		return out
	}
	tr := trueRange(high, low, close)
	denomLog := math.Log10(float64(period))
	for i := period - 1; i < n; i++ {
		var sumTR float64
		hh, ll := high[i-period+1], low[i-period+1]
		for j := i - period + 1; j <= i; j++ {
			sumTR += tr[j]
			if high[j] > hh {
				hh = high[j]
			}
			if low[j] < ll {
				ll = low[j]
			}
		}
		rng := hh - ll
		if rng == 0 || sumTR == 0 {
			continue
		}
		out[i] = 100 * math.Log10(sumTR/rng) / denomLog
	}
	return out
}

// MASSIFn — Mass Index. EMA9 of (high-low), then EMA9 of that, ratio summed
// over `period` bars (Pandas TA default = 25).
func MASSIFn(high, low []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	hl := make([]float64, n)
	for i := 0; i < n; i++ {
		hl[i] = high[i] - low[i]
	}
	e1 := EMAFn(hl, 9)
	e2 := EMAFn(e1, 9)
	ratio := make([]float64, n)
	for i := 0; i < n; i++ {
		if e2[i] != 0 {
			ratio[i] = e1[i] / e2[i]
		}
	}
	var sum float64
	for i := 0; i < period; i++ {
		sum += ratio[i]
	}
	out[period-1] = sum
	for i := period; i < n; i++ {
		sum += ratio[i] - ratio[i-period]
		out[i] = sum
	}
	return out
}

// KVOFn — Klinger Volume Oscillator (Pandas TA flavor).
// Returns the KVO series and an EMA-smoothed signal.
//
//	hlc3   = (h+l+c)/3
//	trend  = sign(hlc3[t] - hlc3[t-1])     (+1 / -1)
//	dm     = high - low
//	cm[t]  = cm[t-1] + dm[t]      if trend[t] == trend[t-1]
//	         dm[t-1]   + dm[t]    otherwise
//	vf     = volume * |2*(dm/cm) - 1| * trend * 100
//	kvo    = EMA(vf, fast) - EMA(vf, slow)
//	sig    = EMA(kvo, signal)
func KVOFn(high, low, close, volume []float64, fast, slow, signal int) (kvo, sig []float64) {
	n := len(close)
	kvo = make([]float64, n)
	sig = make([]float64, n)
	if n < 2 {
		return
	}
	trend := make([]float64, n)
	dm := make([]float64, n)
	cm := make([]float64, n)
	vf := make([]float64, n)
	for i := 0; i < n; i++ {
		dm[i] = high[i] - low[i]
	}
	for i := 1; i < n; i++ {
		hlcCur := (high[i] + low[i] + close[i]) / 3
		hlcPrev := (high[i-1] + low[i-1] + close[i-1]) / 3
		switch {
		case hlcCur > hlcPrev:
			trend[i] = 1
		case hlcCur < hlcPrev:
			trend[i] = -1
		default:
			trend[i] = trend[i-1]
		}
		if trend[i] == trend[i-1] {
			cm[i] = cm[i-1] + dm[i]
		} else {
			cm[i] = dm[i-1] + dm[i]
		}
		if cm[i] != 0 {
			vf[i] = volume[i] * math.Abs(2*(dm[i]/cm[i])-1) * trend[i] * 100
		}
	}
	fastE := EMAFn(vf, fast)
	slowE := EMAFn(vf, slow)
	for i := 0; i < n; i++ {
		kvo[i] = fastE[i] - slowE[i]
	}
	if signal > 0 {
		sig = EMAFn(kvo, signal)
	}
	return
}

// EFIFn — Elder's Force Index: EMA((c - c[-1]) * volume, period).
func EFIFn(close, volume []float64, period int) []float64 {
	n := len(close)
	if n < 2 {
		return make([]float64, n)
	}
	fi := make([]float64, n)
	for i := 1; i < n; i++ {
		fi[i] = (close[i] - close[i-1]) * volume[i]
	}
	return EMAFn(fi, period)
}

// ERIFn — Elder Ray Index. Returns bull power and bear power.
//
//	bull = high - EMA(close, period)
//	bear = low  - EMA(close, period)
func ERIFn(high, low, close []float64, period int) (bull, bear []float64) {
	n := len(close)
	bull = make([]float64, n)
	bear = make([]float64, n)
	if period < 1 || n < period {
		return
	}
	e := EMAFn(close, period)
	for i := 0; i < n; i++ {
		bull[i] = high[i] - e[i]
		bear[i] = low[i] - e[i]
	}
	return
}

// UIFn — Ulcer Index over `period` bars.
//
//	maxC = highest(close, p)
//	dd   = 100 * (close - maxC) / maxC      (a non-positive series)
//	UI   = sqrt( SUM(dd^2, p) / p )
func UIFn(close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	dd2 := make([]float64, n)
	for i := period - 1; i < n; i++ {
		mx := close[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if close[j] > mx {
				mx = close[j]
			}
		}
		if mx != 0 {
			d := 100 * (close[i] - mx) / mx
			dd2[i] = d * d
		}
	}
	for i := period - 1; i < n; i++ {
		var sum float64
		start := i - period + 1
		if start < 0 {
			start = 0
		}
		for j := start; j <= i; j++ {
			sum += dd2[j]
		}
		out[i] = math.Sqrt(sum / float64(period))
	}
	return out
}

// ALMAFn — Arnaud Legoux Moving Average.
//
//	m  = offset * (window - 1)
//	s  = window / sigma
//	w  = exp(-(i - m)^2 / (2 s^2))
//	ALMA = sum(w * price) / sum(w)
func ALMAFn(real []float64, window int, sigma, offset float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if window < 1 || sigma <= 0 || n < window {
		return out
	}
	m := offset * float64(window-1)
	s := float64(window) / sigma
	w := make([]float64, window)
	var wsum float64
	for i := 0; i < window; i++ {
		w[i] = math.Exp(-math.Pow(float64(i)-m, 2) / (2 * s * s))
		wsum += w[i]
	}
	if wsum == 0 {
		return out
	}
	for i := window - 1; i < n; i++ {
		var v float64
		for k := 0; k < window; k++ {
			v += real[i-window+1+k] * w[k]
		}
		out[i] = v / wsum
	}
	return out
}

// RVIFn — Relative Vigor Index using a symmetric [1,2,2,1]/6 weighting.
// Returns RVI and a 4-bar SWMA signal.
//
//	num = SWMA(close - open, 4)
//	den = SWMA(high - low,  4)
//	RVI = SUM(num, p) / SUM(den, p)
func RVIFn(open, high, low, close []float64, period int) (rvi, sig []float64) {
	n := len(close)
	rvi = make([]float64, n)
	sig = make([]float64, n)
	if period < 1 || n < period+3 {
		return
	}
	co := make([]float64, n)
	hl := make([]float64, n)
	for i := 0; i < n; i++ {
		co[i] = close[i] - open[i]
		hl[i] = high[i] - low[i]
	}
	swma := func(s []float64) []float64 {
		o := make([]float64, n)
		for i := 3; i < n; i++ {
			o[i] = (s[i-3] + 2*s[i-2] + 2*s[i-1] + s[i]) / 6
		}
		return o
	}
	num := swma(co)
	den := swma(hl)
	for i := period + 2; i < n; i++ {
		var sn, sd float64
		for k := 0; k < period; k++ {
			sn += num[i-k]
			sd += den[i-k]
		}
		if sd != 0 {
			rvi[i] = sn / sd
		}
	}
	sig = swma(rvi)
	return
}

// EOMFn — Ease of Movement (Pandas TA / cinar default scale).
//
//	distance  = (h+l)/2 - (h_prev+l_prev)/2
//	box_ratio = (volume / divisor) / (h - l)
//	emv       = distance / box_ratio
//	EOM       = SMA(emv, period)
func EOMFn(high, low, volume []float64, period int, divisor float64) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n < period+1 || divisor == 0 {
		return out
	}
	emv := make([]float64, n)
	for i := 1; i < n; i++ {
		dist := (high[i]+low[i])/2 - (high[i-1]+low[i-1])/2
		rng := high[i] - low[i]
		if rng == 0 {
			continue
		}
		box := (volume[i] / divisor) / rng
		if box == 0 {
			continue
		}
		emv[i] = dist / box
	}
	return SMAFn(emv, period)
}

// SQUEEZEFn — LazyBear's Squeeze Momentum.
//
//	bb_upper = SMA(close, bbLen) + bbMult * stddev(close, bbLen)
//	bb_lower = SMA(close, bbLen) - bbMult * stddev(close, bbLen)
//	kc_upper = SMA(close, kcLen) + kcMult * SMA(TR, kcLen)
//	kc_lower = SMA(close, kcLen) - kcMult * SMA(TR, kcLen)
//	on   = 1 when bb_lower > kc_lower && bb_upper < kc_upper, else 0
//	src  = close - avg(avg(highest(high,momLen), lowest(low,momLen)), SMA(close,momLen))
//	mom  = LINEARREG(src, momLen)
func SQUEEZEFn(high, low, close []float64, bbLen int, bbMult float64, kcLen int, kcMult float64, momLen int) (squeeze, momentum []float64) {
	n := len(close)
	squeeze = make([]float64, n)
	momentum = make([]float64, n)
	if n < bbLen || n < kcLen || n < momLen {
		return
	}
	smaC := SMAFn(close, bbLen)
	std := STDDEVFn(close, bbLen, 1.0)
	smaCkc := SMAFn(close, kcLen)
	tr := trueRange(high, low, close)
	rngMA := SMAFn(tr, kcLen)
	for i := 0; i < n; i++ {
		bbUp := smaC[i] + bbMult*std[i]
		bbLo := smaC[i] - bbMult*std[i]
		kcUp := smaCkc[i] + kcMult*rngMA[i]
		kcLo := smaCkc[i] - kcMult*rngMA[i]
		if bbLo > kcLo && bbUp < kcUp {
			squeeze[i] = 1
		}
	}
	src := make([]float64, n)
	smaSrc := SMAFn(close, momLen)
	for i := momLen - 1; i < n; i++ {
		hh, ll := high[i-momLen+1], low[i-momLen+1]
		for j := i - momLen + 2; j <= i; j++ {
			if high[j] > hh {
				hh = high[j]
			}
			if low[j] < ll {
				ll = low[j]
			}
		}
		mid := (hh + ll) / 2
		src[i] = close[i] - (mid+smaSrc[i])/2
	}
	momentum = LINEARREGFn(src, momLen)
	return
}

// PVIFn — Positive Volume Index. Starts at 1000; updates on bars where
// volume rises vs the previous bar; otherwise carries forward.
func PVIFn(close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	out[0] = 1000
	for i := 1; i < n; i++ {
		if volume[i] > volume[i-1] && close[i-1] != 0 {
			out[i] = out[i-1] + (close[i]-close[i-1])/close[i-1]*out[i-1]
		} else {
			out[i] = out[i-1]
		}
	}
	return out
}

// NVIFn — Negative Volume Index. Mirror of PVI but updates on bars where
// volume falls vs the previous bar.
func NVIFn(close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	out[0] = 1000
	for i := 1; i < n; i++ {
		if volume[i] < volume[i-1] && close[i-1] != 0 {
			out[i] = out[i-1] + (close[i]-close[i-1])/close[i-1]*out[i-1]
		} else {
			out[i] = out[i-1]
		}
	}
	return out
}

// STCFn — Schaff Trend Cycle of MACD(fast, slow), then double-stochastic
// smoothed with factor `factor` over `cycle` bars. Pandas TA defaults
// fast=23, slow=50, cycle=10, factor=0.5.
func STCFn(real []float64, fast, slow, cycle int, factor float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	macd := make([]float64, n)
	ef := EMAFn(real, fast)
	es := EMAFn(real, slow)
	for i := 0; i < n; i++ {
		macd[i] = ef[i] - es[i]
	}
	stoch := func(src []float64) []float64 {
		out := make([]float64, n)
		k := make([]float64, n)
		for i := cycle - 1; i < n; i++ {
			ll, hh := src[i-cycle+1], src[i-cycle+1]
			for j := i - cycle + 2; j <= i; j++ {
				if src[j] < ll {
					ll = src[j]
				}
				if src[j] > hh {
					hh = src[j]
				}
			}
			if hh != ll {
				k[i] = (src[i] - ll) / (hh - ll) * 100
			}
		}
		out[cycle-1] = k[cycle-1]
		for i := cycle; i < n; i++ {
			out[i] = out[i-1] + factor*(k[i]-out[i-1])
		}
		return out
	}
	d1 := stoch(macd)
	d2 := stoch(d1)
	return d2
}

// PSLFn — Psychological Line: 100 * (count of bars where close > prev close
// over the last `period` bars) / period.
func PSLFn(close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	up := make([]float64, n)
	for i := 1; i < n; i++ {
		if close[i] > close[i-1] {
			up[i] = 1
		}
	}
	var s float64
	for i := 1; i <= period; i++ {
		s += up[i]
	}
	out[period] = 100 * s / float64(period)
	for i := period + 1; i < n; i++ {
		s += up[i] - up[i-period]
		out[i] = 100 * s / float64(period)
	}
	return out
}

// BIASFn — Bias: 100 * (close - SMA(close,p)) / SMA(close,p).
func BIASFn(close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	sma := SMAFn(close, period)
	for i := 0; i < n; i++ {
		if sma[i] != 0 {
			out[i] = 100 * (close[i] - sma[i]) / sma[i]
		}
	}
	return out
}

// CTIFn — Correlation Trend Indicator: rolling Pearson correlation between
// the price series and a linear time index over `period` bars.
func CTIFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < period {
		return out
	}
	t := make([]float64, n)
	for i := 0; i < n; i++ {
		t[i] = float64(i)
	}
	return CORRELFn(real, t, period)
}

// ERFn — Kaufman Efficiency Ratio: |close - close[-p]| / SUM(|Δclose|, p).
func ERFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	abs := func(x float64) float64 {
		if x < 0 {
			return -x
		}
		return x
	}
	for i := period; i < n; i++ {
		var vol float64
		for j := i - period + 1; j <= i; j++ {
			vol += abs(real[j] - real[j-1])
		}
		if vol != 0 {
			out[i] = abs(real[i]-real[i-period]) / vol
		}
	}
	return out
}

// FWMAFn — Fibonacci Weighted MA: weights are the Fibonacci numbers from
// F_2..F_(period+1), applied newest-weight-largest.
func FWMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	w := make([]float64, period)
	a, b := 1.0, 1.0
	for i := 0; i < period; i++ {
		w[i] = a
		a, b = b, a+b
	}
	var ws float64
	for _, v := range w {
		ws += v
	}
	if ws == 0 {
		return out
	}
	for i := period - 1; i < n; i++ {
		var v float64
		for k := 0; k < period; k++ {
			v += real[i-period+1+k] * w[k]
		}
		out[i] = v / ws
	}
	return out
}

// SWMAFn — Symmetric Weighted MA. The classic 4-bar [1,2,2,1]/6 filter.
func SWMAFn(real []float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	for i := 3; i < n; i++ {
		out[i] = (real[i-3] + 2*real[i-2] + 2*real[i-1] + real[i]) / 6
	}
	return out
}

// INERTIAFn — Pandas TA's Inertia: linear regression of RVI over `period`.
// Defaults: rvi_period=14, regression period=20.
func INERTIAFn(open, high, low, close []float64, rviPeriod, regPeriod int) []float64 {
	rvi, _ := RVIFn(open, high, low, close, rviPeriod)
	return LINEARREGFn(rvi, regPeriod)
}

// DPOFn — Detrended Price Oscillator: close[i - shift] - SMA(close, p)[i],
// where shift = p/2 + 1. Pandas TA `centered=False` flavor.
func DPOFn(close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	shift := period/2 + 1
	sma := SMAFn(close, period)
	for i := 0; i < n; i++ {
		idx := i - shift
		if idx < 0 {
			continue
		}
		out[i] = close[idx] - sma[i]
	}
	return out
}

// ABERRATIONFn — ATR-banded SMA of HLC3.
//
//	mid   = SMA(HLC3, len)
//	upper = mid + atrMult * ATR(h,l,c, atrLen)
//	lower = mid - atrMult * ATR(h,l,c, atrLen)
func ABERRATIONFn(high, low, close []float64, length, atrLen int, atrMult float64) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	if length < 1 || atrLen < 1 || n < length || n <= atrLen {
		return
	}
	hlc3 := make([]float64, n)
	for i := 0; i < n; i++ {
		hlc3[i] = (high[i] + low[i] + close[i]) / 3
	}
	mid := SMAFn(hlc3, length)
	atr := ATRFn(high, low, close, atrLen)
	for i := 0; i < n; i++ {
		middle[i] = mid[i]
		upper[i] = mid[i] + atrMult*atr[i]
		lower[i] = mid[i] - atrMult*atr[i]
	}
	return
}

// PGOFn — Pretty Good Oscillator: (close - SMA(close,p)) / EMA(TR, p).
func PGOFn(high, low, close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	sma := SMAFn(close, period)
	tr := trueRange(high, low, close)
	emaTR := EMAFn(tr, period)
	for i := 0; i < n; i++ {
		if emaTR[i] != 0 {
			out[i] = (close[i] - sma[i]) / emaTR[i]
		}
	}
	return out
}

// RMIFn — Relative Momentum Index. Compares close to close[t-momentum]
// rather than close[t-1] like RSI; Wilder-smoothed up/down moves.
func RMIFn(close []float64, period, momentum int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || momentum < 1 || n <= momentum {
		return out
	}
	up := make([]float64, n)
	dn := make([]float64, n)
	for i := momentum; i < n; i++ {
		d := close[i] - close[i-momentum]
		if d > 0 {
			up[i] = d
		} else {
			dn[i] = -d
		}
	}
	emaUp := SMMAFn(up, period)
	emaDn := SMMAFn(dn, period)
	for i := 0; i < n; i++ {
		den := emaUp[i] + emaDn[i]
		if den != 0 {
			out[i] = 100 * emaUp[i] / den
		}
	}
	return out
}

// MCGINLEYFn — McGinley Dynamic: MD[i] = MD[i-1] + (c - MD[i-1])/(k*p*(c/MD[i-1])^4).
// Seeded with the first close.
func MCGINLEYFn(close []float64, period int, k float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	if period < 1 {
		period = 10
	}
	if k <= 0 {
		k = 0.6
	}
	out[0] = close[0]
	for i := 1; i < n; i++ {
		prev := out[i-1]
		if prev == 0 {
			out[i] = close[i]
			continue
		}
		ratio := close[i] / prev
		denom := k * float64(period) * ratio * ratio * ratio * ratio
		if denom == 0 {
			out[i] = prev
			continue
		}
		out[i] = prev + (close[i]-prev)/denom
	}
	return out
}

// INCREASINGFn — 1 if close[i] > close[i-length], else 0.
func INCREASINGFn(real []float64, length int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if length < 1 || n <= length {
		return out
	}
	for i := length; i < n; i++ {
		if real[i] > real[i-length] {
			out[i] = 1
		}
	}
	return out
}

// DECREASINGFn — 1 if close[i] < close[i-length], else 0.
func DECREASINGFn(real []float64, length int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if length < 1 || n <= length {
		return out
	}
	for i := length; i < n; i++ {
		if real[i] < real[i-length] {
			out[i] = 1
		}
	}
	return out
}

// WADFn — Williams Accumulation/Distribution.
//
//	TRH = max(high, close[-1]),  TRL = min(low, close[-1])
//	if c >  c[-1]: ad = c - TRL
//	if c <  c[-1]: ad = c - TRH
//	if c == c[-1]: ad = 0
//	WAD = cumulative sum of ad
func WADFn(high, low, close []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	if n < 2 {
		return out
	}
	for i := 1; i < n; i++ {
		trh, trl := high[i], low[i]
		if close[i-1] > trh {
			trh = close[i-1]
		}
		if close[i-1] < trl {
			trl = close[i-1]
		}
		var ad float64
		switch {
		case close[i] > close[i-1]:
			ad = close[i] - trl
		case close[i] < close[i-1]:
			ad = close[i] - trh
		}
		out[i] = out[i-1] + ad
	}
	return out
}

// BRARFn — Pandas TA's BRAR. Returns AR and BR over `period` bars.
//
//	AR = 100 * SUM(high - open, p) / SUM(open - low, p)
//	BR = 100 * SUM(max(0, h - c[-1]), p) / SUM(max(0, c[-1] - l), p)
func BRARFn(open, high, low, close []float64, period int) (ar, br []float64) {
	n := len(close)
	ar = make([]float64, n)
	br = make([]float64, n)
	if period < 1 || n <= period {
		return
	}
	hoMinusO := make([]float64, n)
	oMinusL := make([]float64, n)
	hcPrev := make([]float64, n)
	cPrevL := make([]float64, n)
	for i := 0; i < n; i++ {
		hoMinusO[i] = high[i] - open[i]
		oMinusL[i] = open[i] - low[i]
		if i > 0 {
			d1 := high[i] - close[i-1]
			if d1 > 0 {
				hcPrev[i] = d1
			}
			d2 := close[i-1] - low[i]
			if d2 > 0 {
				cPrevL[i] = d2
			}
		}
	}
	var sumHO, sumOL, sumHC, sumCL float64
	for i := 0; i < period; i++ {
		sumHO += hoMinusO[i]
		sumOL += oMinusL[i]
		sumHC += hcPrev[i]
		sumCL += cPrevL[i]
	}
	if sumOL != 0 {
		ar[period-1] = 100 * sumHO / sumOL
	}
	if sumCL != 0 {
		br[period-1] = 100 * sumHC / sumCL
	}
	for i := period; i < n; i++ {
		sumHO += hoMinusO[i] - hoMinusO[i-period]
		sumOL += oMinusL[i] - oMinusL[i-period]
		sumHC += hcPrev[i] - hcPrev[i-period]
		sumCL += cPrevL[i] - cPrevL[i-period]
		if sumOL != 0 {
			ar[i] = 100 * sumHO / sumOL
		}
		if sumCL != 0 {
			br[i] = 100 * sumHC / sumCL
		}
	}
	return
}

// PWMAFn — Pascal Weighted MA. Weights are the binomial coefficients
// C(p-1, k); newest sample carries the central weight.
func PWMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	w := make([]float64, period)
	w[0] = 1
	for i := 1; i < period; i++ {
		w[i] = w[i-1] * float64(period-i) / float64(i)
	}
	var ws float64
	for _, v := range w {
		ws += v
	}
	if ws == 0 {
		return out
	}
	for i := period - 1; i < n; i++ {
		var v float64
		for k := 0; k < period; k++ {
			v += real[i-period+1+k] * w[k]
		}
		out[i] = v / ws
	}
	return out
}

// HILOFn — Gann Hi-Lo activator.
//
//	hi_ma = SMA(high, highLen)
//	lo_ma = SMA(low,  lowLen)
//	trend flips up when close > hi_ma, down when close < lo_ma
//	hilo  = lo_ma in uptrend, hi_ma in downtrend
//	dir   = +1 / -1
func HILOFn(high, low, close []float64, highLen, lowLen int) (hilo, direction []float64) {
	n := len(close)
	hilo = make([]float64, n)
	direction = make([]float64, n)
	if highLen < 1 || lowLen < 1 || n < highLen || n < lowLen {
		return
	}
	hiMA := SMAFn(high, highLen)
	loMA := SMAFn(low, lowLen)
	dir := 1
	for i := 0; i < n; i++ {
		if close[i] > hiMA[i] {
			dir = 1
		} else if close[i] < loMA[i] {
			dir = -1
		}
		direction[i] = float64(dir)
		if dir > 0 {
			hilo[i] = loMA[i]
		} else {
			hilo[i] = hiMA[i]
		}
	}
	return
}

// PDISTFn — Price Distance:
// 2*(high - low) - |close - open| + |open - close[-1]|.
func PDISTFn(open, high, low, close []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	for i := 0; i < n; i++ {
		v := 2 * (high[i] - low[i])
		v -= math.Abs(close[i] - open[i])
		if i > 0 {
			v += math.Abs(open[i] - close[i-1])
		}
		out[i] = v
	}
	return out
}

// VHFFn — Vertical Horizontal Filter:
// (HHV(close,p) - LLV(close,p)) / SUM(|Δclose|, p).
func VHFFn(close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	for i := period; i < n; i++ {
		hh, ll := close[i-period+1], close[i-period+1]
		var sumAbs float64
		for j := i - period + 1; j <= i; j++ {
			if close[j] > hh {
				hh = close[j]
			}
			if close[j] < ll {
				ll = close[j]
			}
			sumAbs += math.Abs(close[j] - close[j-1])
		}
		if sumAbs != 0 {
			out[i] = (hh - ll) / sumAbs
		}
	}
	return out
}

// KCFn — Keltner Channels: EMA(close,p) ± mult * ATR(h,l,c,p).
func KCFn(high, low, close []float64, period int, mult float64) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	if period < 1 || n < period {
		return
	}
	mid := EMAFn(close, period)
	atr := ATRFn(high, low, close, period)
	for i := 0; i < n; i++ {
		middle[i] = mid[i]
		upper[i] = mid[i] + mult*atr[i]
		lower[i] = mid[i] - mult*atr[i]
	}
	return
}

// BBPFn — Bollinger %B: (close - lower) / (upper - lower).
func BBPFn(close []float64, period int, devUp, devDn float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	upper, _, lower := BBANDSFn(close, period, devUp, devDn, SMA)
	for i := 0; i < n; i++ {
		w := upper[i] - lower[i]
		if w != 0 {
			out[i] = (close[i] - lower[i]) / w
		}
	}
	return out
}

// BBWFn — Bollinger Bandwidth (%): 100 * (upper - lower) / middle.
func BBWFn(close []float64, period int, devUp, devDn float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	upper, mid, lower := BBANDSFn(close, period, devUp, devDn, SMA)
	for i := 0; i < n; i++ {
		if mid[i] != 0 {
			out[i] = 100 * (upper[i] - lower[i]) / mid[i]
		}
	}
	return out
}

// CRSIFn — Connors RSI = (RSI(close, rsiPeriod) + RSI(streak, streakPeriod) +
// percent_rank(ROC(close,1), rankPeriod)) / 3.
func CRSIFn(close []float64, rsiPeriod, streakPeriod, rankPeriod int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if n < 2 {
		return out
	}
	rsi := RSIFn(close, rsiPeriod)
	streak := make([]float64, n)
	for i := 1; i < n; i++ {
		switch {
		case close[i] > close[i-1]:
			if streak[i-1] >= 0 {
				streak[i] = streak[i-1] + 1
			} else {
				streak[i] = 1
			}
		case close[i] < close[i-1]:
			if streak[i-1] <= 0 {
				streak[i] = streak[i-1] - 1
			} else {
				streak[i] = -1
			}
		}
	}
	rsiStreak := RSIFn(streak, streakPeriod)
	roc := make([]float64, n)
	for i := 1; i < n; i++ {
		if close[i-1] != 0 {
			roc[i] = (close[i] - close[i-1]) / close[i-1] * 100
		}
	}
	rank := make([]float64, n)
	for i := rankPeriod; i < n; i++ {
		var below int
		ref := roc[i]
		for j := i - rankPeriod + 1; j < i; j++ {
			if roc[j] < ref {
				below++
			}
		}
		rank[i] = float64(below) / float64(rankPeriod-1) * 100
	}
	for i := 0; i < n; i++ {
		out[i] = (rsi[i] + rsiStreak[i] + rank[i]) / 3
	}
	return out
}

// IFISHERFn — Inverse Fisher Transform of `signal`. Useful applied to
// RSI/etc to compress values into [-1, +1].
//
//	v   = amplitude * signal
//	out = (e^(2v) - 1) / (e^(2v) + 1)
func IFISHERFn(signal []float64, amplitude float64) []float64 {
	n := len(signal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		v := amplitude * signal[i]
		e := math.Exp(2 * v)
		out[i] = (e - 1) / (e + 1)
	}
	return out
}

// VSTOPFn — Volatility Stop (ATR trailing stop). Returns the trailing stop
// price and direction (+1/-1).
func VSTOPFn(high, low, close []float64, period int, mult float64) (stop, direction []float64) {
	n := len(close)
	stop = make([]float64, n)
	direction = make([]float64, n)
	if period < 1 || n <= period {
		return
	}
	atr := ATRFn(high, low, close, period)
	dir := 1
	stop[period] = close[period] - mult*atr[period]
	direction[period] = 1
	for i := period + 1; i < n; i++ {
		off := mult * atr[i]
		if dir > 0 {
			s := close[i] - off
			if s < stop[i-1] {
				s = stop[i-1]
			}
			if close[i] < stop[i-1] {
				dir = -1
				stop[i] = close[i] + off
			} else {
				stop[i] = s
			}
		} else {
			s := close[i] + off
			if s > stop[i-1] {
				s = stop[i-1]
			}
			if close[i] > stop[i-1] {
				dir = 1
				stop[i] = close[i] - off
			} else {
				stop[i] = s
			}
		}
		direction[i] = float64(dir)
	}
	return
}

// ENVELOPEFn — Moving-average envelope: SMA(close, p) * (1 ± pct/100).
func ENVELOPEFn(close []float64, period int, pct float64) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	mid := SMAFn(close, period)
	f := pct / 100.0
	for i := 0; i < n; i++ {
		middle[i] = mid[i]
		upper[i] = mid[i] * (1 + f)
		lower[i] = mid[i] * (1 - f)
	}
	return
}

// HLC3Fn — Typical Price (h+l+c)/3.
func HLC3Fn(high, low, close []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (high[i] + low[i] + close[i]) / 3
	}
	return out
}

// OHLC4Fn — (open + high + low + close) / 4.
func OHLC4Fn(open, high, low, close []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (open[i] + high[i] + low[i] + close[i]) / 4
	}
	return out
}

// MADFn — Mean Absolute Deviation over `period` bars around the rolling mean.
func MADFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	mean := SMAFn(real, period)
	for i := period - 1; i < n; i++ {
		var s float64
		for j := i - period + 1; j <= i; j++ {
			s += math.Abs(real[j] - mean[i])
		}
		out[i] = s / float64(period)
	}
	return out
}

// ENTROPYFn — Shannon entropy (natural log) of the sliding window normalized
// to a probability distribution. Negative / zero values contribute nothing.
func ENTROPYFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			if real[j] > 0 {
				sum += real[j]
			}
		}
		if sum == 0 {
			continue
		}
		var h float64
		for j := i - period + 1; j <= i; j++ {
			if real[j] > 0 {
				p := real[j] / sum
				h -= p * math.Log(p)
			}
		}
		out[i] = h
	}
	return out
}

// MEDIANFn — Rolling median over `period` bars.
func MEDIANFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	buf := make([]float64, period)
	for i := period - 1; i < n; i++ {
		copy(buf, real[i-period+1:i+1])
		// insertion sort — small windows.
		for k := 1; k < period; k++ {
			x := buf[k]
			j := k - 1
			for j >= 0 && buf[j] > x {
				buf[j+1] = buf[j]
				j--
			}
			buf[j+1] = x
		}
		if period%2 == 1 {
			out[i] = buf[period/2]
		} else {
			out[i] = (buf[period/2-1] + buf[period/2]) / 2
		}
	}
	return out
}

// FRAMAFn — Ehlers Fractal Adaptive Moving Average over an even `period`.
//
//	N1 = (max(H, n/2_first)  - min(L, n/2_first))  / (n/2)
//	N2 = (max(H, n/2_second) - min(L, n/2_second)) / (n/2)
//	N3 = (max(H, n)          - min(L, n))          / n
//	D     = (log(N1 + N2) - log(N3)) / log(2)
//	alpha = clamp(exp(-4.6 * (D - 1)), 0.01, 1)
//	out[i] = alpha * close[i] + (1 - alpha) * out[i-1]
func FRAMAFn(high, low, close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 2 || period%2 != 0 || n < period {
		return out
	}
	half := period / 2
	out[period-1] = close[period-1]
	for i := period; i < n; i++ {
		h1, l1 := high[i-period+1], low[i-period+1]
		for j := i - period + 1; j < i-period+1+half; j++ {
			if high[j] > h1 {
				h1 = high[j]
			}
			if low[j] < l1 {
				l1 = low[j]
			}
		}
		h2, l2 := high[i-period+1+half], low[i-period+1+half]
		for j := i - period + 1 + half; j <= i; j++ {
			if high[j] > h2 {
				h2 = high[j]
			}
			if low[j] < l2 {
				l2 = low[j]
			}
		}
		h3, l3 := high[i-period+1], low[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if high[j] > h3 {
				h3 = high[j]
			}
			if low[j] < l3 {
				l3 = low[j]
			}
		}
		N1 := (h1 - l1) / float64(half)
		N2 := (h2 - l2) / float64(half)
		N3 := (h3 - l3) / float64(period)
		var alpha float64 = 0.01
		if N1+N2 > 0 && N3 > 0 {
			D := (math.Log(N1+N2) - math.Log(N3)) / math.Log(2)
			alpha = math.Exp(-4.6 * (D - 1))
		}
		if alpha > 1 {
			alpha = 1
		} else if alpha < 0.01 {
			alpha = 0.01
		}
		out[i] = alpha*close[i] + (1-alpha)*out[i-1]
	}
	return out
}

// ZSCOREFn — rolling Z-score: (real - SMA(real,p)) / STDDEV(real,p).
func ZSCOREFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	mean := SMAFn(real, period)
	std := STDDEVFn(real, period, 1.0)
	for i := 0; i < n; i++ {
		if std[i] != 0 {
			out[i] = (real[i] - mean[i]) / std[i]
		}
	}
	return out
}

// DECAYFn — Linear decay: out[i] = max(real[i], out[i-1] - 1/length, 0).
// Pandas TA's `decay(mode='linear')` flavor.
func DECAYFn(real []float64, length int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	if length < 1 {
		length = 5
	}
	step := 1.0 / float64(length)
	out[0] = real[0]
	for i := 1; i < n; i++ {
		v := out[i-1] - step
		if real[i] > v {
			v = real[i]
		}
		if v < 0 {
			v = 0
		}
		out[i] = v
	}
	return out
}

// TIIFn — Trend Intensity Index over (smaPeriod, sumPeriod).
//
//	sma = SMA(close, smaPeriod)
//	positive = max(close - sma, 0); negative = max(sma - close, 0)
//	TII = 100 * SUM(positive, sumPeriod) / (SUM(positive, sumPeriod) + SUM(negative, sumPeriod))
func TIIFn(close []float64, smaPeriod, sumPeriod int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if smaPeriod < 1 || sumPeriod < 1 || n < smaPeriod+sumPeriod {
		return out
	}
	sma := SMAFn(close, smaPeriod)
	pos := make([]float64, n)
	neg := make([]float64, n)
	for i := 0; i < n; i++ {
		d := close[i] - sma[i]
		if d > 0 {
			pos[i] = d
		} else {
			neg[i] = -d
		}
	}
	var sP, sN float64
	for i := 0; i < sumPeriod; i++ {
		sP += pos[i]
		sN += neg[i]
	}
	for i := sumPeriod - 1; i < n; i++ {
		if i >= sumPeriod {
			sP += pos[i] - pos[i-sumPeriod]
			sN += neg[i] - neg[i-sumPeriod]
		}
		denom := sP + sN
		if denom != 0 {
			out[i] = 100 * sP / denom
		}
	}
	return out
}

// SMIFn — Stochastic Momentum Index. Returns smi and an EMA-smoothed signal.
//
//	mid   = (HHV(high, len) + LLV(low, len)) / 2
//	rng   = HHV(high, len) - LLV(low, len)
//	m1    = EMA(EMA(close - mid, smooth), smooth)
//	m2    = EMA(EMA(rng/2,        smooth), smooth)
//	SMI   = 100 * m1 / m2
//	sig   = EMA(SMI, signal)
func SMIFn(high, low, close []float64, length, smooth, signal int) (smi, sig []float64) {
	n := len(close)
	smi = make([]float64, n)
	sig = make([]float64, n)
	if length < 1 || n < length {
		return
	}
	diff := make([]float64, n)
	rng := make([]float64, n)
	for i := length - 1; i < n; i++ {
		hh, ll := high[i-length+1], low[i-length+1]
		for j := i - length + 2; j <= i; j++ {
			if high[j] > hh {
				hh = high[j]
			}
			if low[j] < ll {
				ll = low[j]
			}
		}
		diff[i] = close[i] - (hh+ll)/2
		rng[i] = (hh - ll) / 2
	}
	m1 := EMAFn(EMAFn(diff, smooth), smooth)
	m2 := EMAFn(EMAFn(rng, smooth), smooth)
	for i := 0; i < n; i++ {
		if m2[i] != 0 {
			smi[i] = 100 * m1[i] / m2[i]
		}
	}
	if signal > 0 {
		sig = EMAFn(smi, signal)
	}
	return
}

// KURTOSISFn — rolling sample excess kurtosis: m4/m2^2 - 3.
func KURTOSISFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < period {
		return out
	}
	mean := SMAFn(real, period)
	for i := period - 1; i < n; i++ {
		var m2, m4 float64
		for j := i - period + 1; j <= i; j++ {
			d := real[j] - mean[i]
			m2 += d * d
			m4 += d * d * d * d
		}
		m2 /= float64(period)
		m4 /= float64(period)
		if m2 != 0 {
			out[i] = m4/(m2*m2) - 3
		}
	}
	return out
}

// SKEWFn — rolling sample skewness: m3 / m2^(3/2).
func SKEWFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < period {
		return out
	}
	mean := SMAFn(real, period)
	for i := period - 1; i < n; i++ {
		var m2, m3 float64
		for j := i - period + 1; j <= i; j++ {
			d := real[j] - mean[i]
			m2 += d * d
			m3 += d * d * d
		}
		m2 /= float64(period)
		m3 /= float64(period)
		if m2 > 0 {
			out[i] = m3 / math.Pow(m2, 1.5)
		}
	}
	return out
}

// QUANTILEFn — rolling quantile q ∈ [0,1] with linear interpolation.
func QUANTILEFn(real []float64, period int, q float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period || q < 0 || q > 1 {
		return out
	}
	buf := make([]float64, period)
	for i := period - 1; i < n; i++ {
		copy(buf, real[i-period+1:i+1])
		for k := 1; k < period; k++ {
			x := buf[k]
			j := k - 1
			for j >= 0 && buf[j] > x {
				buf[j+1] = buf[j]
				j--
			}
			buf[j+1] = x
		}
		pos := q * float64(period-1)
		lo := int(math.Floor(pos))
		hi := int(math.Ceil(pos))
		if hi >= period {
			hi = period - 1
		}
		if lo == hi {
			out[i] = buf[lo]
		} else {
			frac := pos - float64(lo)
			out[i] = buf[lo]*(1-frac) + buf[hi]*frac
		}
	}
	return out
}

// ACCBANDSFn — Acceleration Bands.
//
//	hl_ratio = (h - l) / (h + l)
//	upper    = SMA(h * (1 + 4*hl_ratio), period)
//	lower    = SMA(l * (1 - 4*hl_ratio), period)
//	middle   = SMA(close, period)
func ACCBANDSFn(high, low, close []float64, period int) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	if period < 1 || n < period {
		return
	}
	hu := make([]float64, n)
	ll := make([]float64, n)
	for i := 0; i < n; i++ {
		s := high[i] + low[i]
		var r float64
		if s != 0 {
			r = (high[i] - low[i]) / s
		}
		hu[i] = high[i] * (1 + 4*r)
		ll[i] = low[i] * (1 - 4*r)
	}
	upper = SMAFn(hu, period)
	lower = SMAFn(ll, period)
	middle = SMAFn(close, period)
	return
}

// SSFFn — Ehlers Super Smoother Filter (2-pole).
//
//	a1 = exp(-π√2 / period)
//	b1 = 2 a1 cos(180·√2 / period)
//	c2 = b1, c3 = -a1², c1 = 1 - c2 - c3
//	out[i] = c1*(real[i]+real[i-1])/2 + c2*out[i-1] + c3*out[i-2]
func SSFFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < 2 {
		return out
	}
	a1 := math.Exp(-math.Pi * math.Sqrt2 / float64(period))
	b1 := 2 * a1 * math.Cos(math.Sqrt2*math.Pi/float64(period))
	c2 := b1
	c3 := -a1 * a1
	c1 := 1 - c2 - c3
	out[0] = real[0]
	out[1] = real[1]
	for i := 2; i < n; i++ {
		out[i] = c1*(real[i]+real[i-1])/2 + c2*out[i-1] + c3*out[i-2]
	}
	return out
}

// CFOFn — Chande Forecast Oscillator: 100 * (real - TSF(real, p)) / real.
func CFOFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	tsf := TSFFn(real, period)
	for i := 0; i < n; i++ {
		if real[i] != 0 {
			out[i] = 100 * (real[i] - tsf[i]) / real[i]
		}
	}
	return out
}

// VIDYAFn — Chande's Variable Index Dynamic Average. EMA whose smoothing
// is scaled by |CMO|/100 over `cmoPeriod`.
func VIDYAFn(real []float64, period, cmoPeriod int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || cmoPeriod < 1 || n <= cmoPeriod {
		return out
	}
	cmo := CMOFn(real, cmoPeriod)
	alpha := 2.0 / float64(period+1)
	out[cmoPeriod] = real[cmoPeriod]
	for i := cmoPeriod + 1; i < n; i++ {
		k := alpha * math.Abs(cmo[i]) / 100
		if k > 1 {
			k = 1
		}
		out[i] = k*real[i] + (1-k)*out[i-1]
	}
	return out
}

// CHAIKINVOLFn — Chaikin Volatility:
// 100 * (EMA(H-L, emaP)[i] - EMA(H-L, emaP)[i-rocP]) / EMA(H-L, emaP)[i-rocP].
func CHAIKINVOLFn(high, low []float64, emaPeriod, rocPeriod int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if emaPeriod < 1 || rocPeriod < 1 || n <= rocPeriod {
		return out
	}
	hl := make([]float64, n)
	for i := 0; i < n; i++ {
		hl[i] = high[i] - low[i]
	}
	e := EMAFn(hl, emaPeriod)
	for i := rocPeriod; i < n; i++ {
		if e[i-rocPeriod] != 0 {
			out[i] = 100 * (e[i] - e[i-rocPeriod]) / e[i-rocPeriod]
		}
	}
	return out
}

// HEIKINASHIFn — Heikin-Ashi candle series. Returns ha_open, ha_high, ha_low, ha_close.
//
//	ha_close = (o + h + l + c) / 4
//	ha_open  = (ha_open[i-1] + ha_close[i-1]) / 2  (seed: (o[0]+c[0])/2)
//	ha_high  = max(h, ha_open, ha_close)
//	ha_low   = min(l, ha_open, ha_close)
func HEIKINASHIFn(open, high, low, close []float64) (haOpen, haHigh, haLow, haClose []float64) {
	n := len(close)
	haOpen = make([]float64, n)
	haHigh = make([]float64, n)
	haLow = make([]float64, n)
	haClose = make([]float64, n)
	if n == 0 {
		return
	}
	haClose[0] = (open[0] + high[0] + low[0] + close[0]) / 4
	haOpen[0] = (open[0] + close[0]) / 2
	haHigh[0] = high[0]
	haLow[0] = low[0]
	for i := 1; i < n; i++ {
		haClose[i] = (open[i] + high[i] + low[i] + close[i]) / 4
		haOpen[i] = (haOpen[i-1] + haClose[i-1]) / 2
		hi := high[i]
		if haOpen[i] > hi {
			hi = haOpen[i]
		}
		if haClose[i] > hi {
			hi = haClose[i]
		}
		haHigh[i] = hi
		lo := low[i]
		if haOpen[i] < lo {
			lo = haOpen[i]
		}
		if haClose[i] < lo {
			lo = haClose[i]
		}
		haLow[i] = lo
	}
	return
}

// TRIXSIGNALFn — TRIX plus its SMA signal line.
func TRIXSIGNALFn(real []float64, period, signalPeriod int) (trix, signal []float64) {
	trix = TRIXFn(real, period)
	if signalPeriod > 0 {
		signal = SMAFn(trix, signalPeriod)
	} else {
		signal = make([]float64, len(real))
	}
	return
}

// CHANDELIEREXITFn — Chandelier Exit. Returns long-side and short-side trailing stops.
//
//	atr     = ATR(h,l,c, p) * mult
//	longEx  = HHV(high, p) - atr
//	shortEx = LLV(low,  p) + atr
func CHANDELIEREXITFn(high, low, close []float64, period int, mult float64) (longExit, shortExit []float64) {
	n := len(close)
	longExit = make([]float64, n)
	shortExit = make([]float64, n)
	if period < 1 || n <= period {
		return
	}
	atr := ATRFn(high, low, close, period)
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
		longExit[i] = hh - mult*atr[i]
		shortExit[i] = ll + mult*atr[i]
	}
	return
}

// OBVSMOOTHFn — EMA-smoothed On-Balance Volume.
func OBVSMOOTHFn(close, volume []float64, period int) []float64 {
	return EMAFn(OBVFn(close, volume), period)
}

// ZLHMAFn — Zero-Lag Hull MA: HMA computed on a de-lagged source series.
//
//	lag = (period - 1) / 2
//	src[i] = 2 * real[i] - real[i - lag]
//	out    = HMA(src, period)
func ZLHMAFn(real []float64, period int) []float64 {
	n := len(real)
	if period < 2 || n < period {
		return make([]float64, n)
	}
	lag := (period - 1) / 2
	src := make([]float64, n)
	for i := 0; i < n; i++ {
		if i < lag {
			src[i] = real[i]
		} else {
			src[i] = 2*real[i] - real[i-lag]
		}
	}
	return HMAFn(src, period)
}

// GDFn — Generalized DEMA: (1 + v) * EMA(real, p) - v * EMA(EMA(real, p), p).
func GDFn(real []float64, period int, v float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	e1 := EMAFn(real, period)
	e2 := EMAFn(e1, period)
	for i := 0; i < n; i++ {
		out[i] = (1+v)*e1[i] - v*e2[i]
	}
	return out
}

// PCTRETFn — Percent return: 100 * (real[i] - real[i-period]) / real[i-period].
func PCTRETFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	for i := period; i < n; i++ {
		if real[i-period] != 0 {
			out[i] = 100 * (real[i] - real[i-period]) / real[i-period]
		}
	}
	return out
}

// LOGRETFn — Log return: ln(real[i] / real[i-period]).
func LOGRETFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	for i := period; i < n; i++ {
		if real[i] > 0 && real[i-period] > 0 {
			out[i] = math.Log(real[i] / real[i-period])
		}
	}
	return out
}

// TTMTRENDFn — TTM Trend bar: +1 if close > SMA((H+L)/2, period), else -1.
func TTMTRENDFn(high, low, close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	hl2 := make([]float64, n)
	for i := 0; i < n; i++ {
		hl2[i] = (high[i] + low[i]) / 2
	}
	avg := SMAFn(hl2, period)
	for i := period - 1; i < n; i++ {
		if close[i] > avg[i] {
			out[i] = 1
		} else {
			out[i] = -1
		}
	}
	return out
}

// AOACCFn — Bill Williams' Acceleration/Deceleration: AO - SMA(AO, 5).
func AOACCFn(high, low []float64) []float64 {
	n := len(high)
	ao := AOFn(high, low)
	smaAO := SMAFn(ao, 5)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = ao[i] - smaAO[i]
	}
	return out
}

// ROCSFn — Smoothed ROC: SMA(ROC(real, p), smooth).
func ROCSFn(real []float64, period, smooth int) []float64 {
	return SMAFn(ROCFn(real, period), smooth)
}

// WTFn — LazyBear's Wave Trend Oscillator. Returns wt1 and wt2 (SMA(wt1, 4)).
//
//	ap   = HLC3
//	esa  = EMA(ap, n1)
//	d    = EMA(|ap - esa|, n1)
//	ci   = (ap - esa) / (0.015 * d)
//	wt1  = EMA(ci, n2)
//	wt2  = SMA(wt1, 4)
func WTFn(high, low, close []float64, n1, n2 int) (wt1, wt2 []float64) {
	n := len(close)
	wt1 = make([]float64, n)
	wt2 = make([]float64, n)
	if n1 < 1 || n2 < 1 || n < n1 {
		return
	}
	ap := make([]float64, n)
	for i := 0; i < n; i++ {
		ap[i] = (high[i] + low[i] + close[i]) / 3
	}
	esa := EMAFn(ap, n1)
	diff := make([]float64, n)
	for i := 0; i < n; i++ {
		diff[i] = math.Abs(ap[i] - esa[i])
	}
	d := EMAFn(diff, n1)
	ci := make([]float64, n)
	for i := 0; i < n; i++ {
		if d[i] != 0 {
			ci[i] = (ap[i] - esa[i]) / (0.015 * d[i])
		}
	}
	wt1 = EMAFn(ci, n2)
	wt2 = SMAFn(wt1, 4)
	return
}

// VOFn — Volume Oscillator: 100 * (EMA(vol, fast) - EMA(vol, slow)) / EMA(vol, slow).
func VOFn(volume []float64, fast, slow int) []float64 {
	n := len(volume)
	out := make([]float64, n)
	fastE := EMAFn(volume, fast)
	slowE := EMAFn(volume, slow)
	for i := 0; i < n; i++ {
		if slowE[i] != 0 {
			out[i] = 100 * (fastE[i] - slowE[i]) / slowE[i]
		}
	}
	return out
}

// PVTFn — Price Volume Trend (cumulative).
//
//	pvt[i] = pvt[i-1] + volume[i] * (close[i] - close[i-1]) / close[i-1]
func PVTFn(close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 1; i < n; i++ {
		if close[i-1] != 0 {
			out[i] = out[i-1] + volume[i]*(close[i]-close[i-1])/close[i-1]
		} else {
			out[i] = out[i-1]
		}
	}
	return out
}

// PVRFn — Price Volume Rank (cinar). 1..4 categorical ranks per bar:
//
//	close↑ vol↑ → 1
//	close↑ vol↓ → 2
//	close↓ vol↑ → 3
//	close↓ vol↓ → 4
func PVRFn(close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 1; i < n; i++ {
		cu := close[i] > close[i-1]
		vu := volume[i] > volume[i-1]
		switch {
		case cu && vu:
			out[i] = 1
		case cu && !vu:
			out[i] = 2
		case !cu && vu:
			out[i] = 3
		default:
			out[i] = 4
		}
	}
	return out
}

// MFVFn — Per-bar Money Flow Volume:
// ((close-low) - (high-close)) / (high-low) * volume.
func MFVFn(high, low, close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		rng := high[i] - low[i]
		if rng == 0 {
			continue
		}
		mfm := ((close[i] - low[i]) - (high[i] - close[i])) / rng
		out[i] = mfm * volume[i]
	}
	return out
}

// DEMFn — DeMarker indicator over `period` bars.
//
//	demax[i] = max(high[i] - high[i-1], 0)
//	demin[i] = max(low[i-1] - low[i], 0)
//	DEM      = SMA(demax,p) / (SMA(demax,p) + SMA(demin,p))
func DEMFn(high, low []float64, period int) []float64 {
	n := len(high)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	dmax := make([]float64, n)
	dmin := make([]float64, n)
	for i := 1; i < n; i++ {
		d1 := high[i] - high[i-1]
		if d1 > 0 {
			dmax[i] = d1
		}
		d2 := low[i-1] - low[i]
		if d2 > 0 {
			dmin[i] = d2
		}
	}
	smaMax := SMAFn(dmax, period)
	smaMin := SMAFn(dmin, period)
	for i := 0; i < n; i++ {
		denom := smaMax[i] + smaMin[i]
		if denom != 0 {
			out[i] = smaMax[i] / denom
		}
	}
	return out
}

// RSISMOOTHFn — EMA-smoothed RSI.
func RSISMOOTHFn(close []float64, rsiPeriod, smoothPeriod int) []float64 {
	return EMAFn(RSIFn(close, rsiPeriod), smoothPeriod)
}

// EMAENVFn — EMA-based envelope: EMA(close, p) * (1 ± pct/100).
func EMAENVFn(close []float64, period int, pct float64) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	mid := EMAFn(close, period)
	f := pct / 100.0
	for i := 0; i < n; i++ {
		middle[i] = mid[i]
		upper[i] = mid[i] * (1 + f)
		lower[i] = mid[i] * (1 - f)
	}
	return
}

// MACDZLFn — Zero-Lag MACD: ZLEMA(fast) - ZLEMA(slow), signal via ZLEMA.
func MACDZLFn(real []float64, fast, slow, signal int) (macd, sig, hist []float64) {
	n := len(real)
	fastE := ZLEMAFn(real, fast)
	slowE := ZLEMAFn(real, slow)
	macd = make([]float64, n)
	for i := 0; i < n; i++ {
		macd[i] = fastE[i] - slowE[i]
	}
	sig = ZLEMAFn(macd, signal)
	hist = make([]float64, n)
	for i := 0; i < n; i++ {
		hist[i] = macd[i] - sig[i]
	}
	return
}

// VOLRATIOFn — current volume / SMA(volume, period).
func VOLRATIOFn(volume []float64, period int) []float64 {
	n := len(volume)
	out := make([]float64, n)
	avg := SMAFn(volume, period)
	for i := 0; i < n; i++ {
		if avg[i] != 0 {
			out[i] = volume[i] / avg[i]
		}
	}
	return out
}

// KCBFn — Keltner %B: (close - lower_kc) / (upper_kc - lower_kc).
func KCBFn(high, low, close []float64, period int, mult float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	upper, _, lower := KCFn(high, low, close, period, mult)
	for i := 0; i < n; i++ {
		w := upper[i] - lower[i]
		if w != 0 {
			out[i] = (close[i] - lower[i]) / w
		}
	}
	return out
}

// TRENDFLEXFn — Ehlers TrendFlex. SSF(close, p) prefilters; trendflex is the
// average of past p one-step slopes normalized by its rolling RMS.
func TRENDFLEXFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	filt := SSFFn(real, period)
	ms := make([]float64, n)
	for i := period; i < n; i++ {
		var sum float64
		for j := 1; j <= period; j++ {
			sum += filt[i] - filt[i-j]
		}
		slope := sum / float64(period)
		ms[i] = 0.04*slope*slope + 0.96*ms[i-1]
		if ms[i] > 0 {
			out[i] = slope / math.Sqrt(ms[i])
		}
	}
	return out
}

// TMFFn — Twiggs Money Flow:
// EMA(volume * ((c - trueLow) - (trueHigh - c)) / (trueHigh - trueLow), p)
// divided by EMA(volume, p), where trueHigh/trueLow use the prior close.
func TMFFn(high, low, close, volume []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n < 2 {
		return out
	}
	adv := make([]float64, n)
	for i := 1; i < n; i++ {
		th := high[i]
		if close[i-1] > th {
			th = close[i-1]
		}
		tl := low[i]
		if close[i-1] < tl {
			tl = close[i-1]
		}
		rng := th - tl
		if rng == 0 {
			continue
		}
		adv[i] = volume[i] * ((close[i] - tl) - (th - close[i])) / rng
	}
	ea := EMAFn(adv, period)
	ev := EMAFn(volume, period)
	for i := 0; i < n; i++ {
		if ev[i] != 0 {
			out[i] = ea[i] / ev[i]
		}
	}
	return out
}

// TVIFn — Trade Volume Index. Direction +1 if Δclose > minTick, -1 if
// Δclose < -minTick, else 0; cumulative sum of direction*volume.
func TVIFn(close, volume []float64, minTick float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 1; i < n; i++ {
		d := close[i] - close[i-1]
		var dir float64
		switch {
		case d > minTick:
			dir = 1
		case d < -minTick:
			dir = -1
		}
		out[i] = out[i-1] + dir*volume[i]
	}
	return out
}

// CVDFn — Cumulative Volume Delta (price-tick approximation):
// signs each bar's volume by close vs prior close, then cumulates.
func CVDFn(close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 1; i < n; i++ {
		var d float64
		switch {
		case close[i] > close[i-1]:
			d = volume[i]
		case close[i] < close[i-1]:
			d = -volume[i]
		}
		out[i] = out[i-1] + d
	}
	return out
}

// REFLEXFn — Ehlers Reflex (companion to TrendFlex). Detects mean-reversion.
//
//	filt = SSF(close, period)
//	slope = (filt[i] - filt[i-period]) / period
//	sum  = avg over j=1..period of (filt[i] + j*slope - filt[i-j])
//	ms[i] = 0.04 * sum² + 0.96 * ms[i-1]
//	reflex[i] = sum / sqrt(ms[i]) if ms[i] > 0
func REFLEXFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	filt := SSFFn(real, period)
	ms := make([]float64, n)
	for i := period; i < n; i++ {
		slope := (filt[i] - filt[i-period]) / float64(period)
		var sum float64
		for j := 1; j <= period; j++ {
			sum += (filt[i] + float64(j)*slope) - filt[i-j]
		}
		sum /= float64(period)
		ms[i] = 0.04*sum*sum + 0.96*ms[i-1]
		if ms[i] > 0 {
			out[i] = sum / math.Sqrt(ms[i])
		}
	}
	return out
}

// FRACTALFn — Bill Williams Fractals. Marks bar i-2 with 1.0 when it is the
// centre of an Up or Down 5-bar fractal pattern.
func FRACTALFn(high, low []float64) (up, down []float64) {
	n := len(high)
	up = make([]float64, n)
	down = make([]float64, n)
	for i := 4; i < n; i++ {
		c := i - 2
		if high[c] > high[c-1] && high[c] > high[c-2] && high[c] > high[c+1] && high[c] > high[i] {
			up[c] = 1
		}
		if low[c] < low[c-1] && low[c] < low[c-2] && low[c] < low[c+1] && low[c] < low[i] {
			down[c] = 1
		}
	}
	return
}

// ALLIGATORFn — Bill Williams Alligator. Returns jaw (SMMA(13) lagged 8),
// teeth (SMMA(8) lagged 5), and lips (SMMA(5) lagged 3) of the median price.
// Lagged values for early bars are 0.
func ALLIGATORFn(high, low []float64) (jaw, teeth, lips []float64) {
	n := len(high)
	jaw = make([]float64, n)
	teeth = make([]float64, n)
	lips = make([]float64, n)
	med := make([]float64, n)
	for i := 0; i < n; i++ {
		med[i] = (high[i] + low[i]) / 2
	}
	smJaw := SMMAFn(med, 13)
	smTeeth := SMMAFn(med, 8)
	smLips := SMMAFn(med, 5)
	for i := 0; i < n; i++ {
		if i-8 >= 0 {
			jaw[i] = smJaw[i-8]
		}
		if i-5 >= 0 {
			teeth[i] = smTeeth[i-5]
		}
		if i-3 >= 0 {
			lips[i] = smLips[i-3]
		}
	}
	return
}

// GATORFn — Gator Oscillator: upper = |jaw - teeth|, lower = -|teeth - lips|.
func GATORFn(high, low []float64) (upper, lower []float64) {
	n := len(high)
	upper = make([]float64, n)
	lower = make([]float64, n)
	jaw, teeth, lips := ALLIGATORFn(high, low)
	for i := 0; i < n; i++ {
		upper[i] = math.Abs(jaw[i] - teeth[i])
		lower[i] = -math.Abs(teeth[i] - lips[i])
	}
	return
}

// SQUEEZEPROFn — TTM Squeeze Pro. Tracks compression at 3 KC widths
// (low=1.0, mid=1.5, high=2.0 ATR-multipliers) plus the LINEARREG-style
// momentum series. Each `squeeze_*` is 1 when Bollinger sits inside Keltner
// at that level, else 0.
func SQUEEZEPROFn(high, low, close []float64, bbLen int, bbMult float64, kcLen int, mLow, mMid, mHigh float64, momLen int) (sqLow, sqMid, sqHigh, momentum []float64) {
	n := len(close)
	sqLow = make([]float64, n)
	sqMid = make([]float64, n)
	sqHigh = make([]float64, n)
	momentum = make([]float64, n)
	if n < bbLen || n < kcLen || n < momLen {
		return
	}
	smaC := SMAFn(close, bbLen)
	std := STDDEVFn(close, bbLen, 1.0)
	smaCkc := SMAFn(close, kcLen)
	tr := trueRange(high, low, close)
	rngMA := SMAFn(tr, kcLen)
	for i := 0; i < n; i++ {
		bbUp := smaC[i] + bbMult*std[i]
		bbLo := smaC[i] - bbMult*std[i]
		check := func(mult float64) float64 {
			kcUp := smaCkc[i] + mult*rngMA[i]
			kcLo := smaCkc[i] - mult*rngMA[i]
			if bbLo > kcLo && bbUp < kcUp {
				return 1
			}
			return 0
		}
		sqLow[i] = check(mLow)
		sqMid[i] = check(mMid)
		sqHigh[i] = check(mHigh)
	}
	src := make([]float64, n)
	smaSrc := SMAFn(close, momLen)
	for i := momLen - 1; i < n; i++ {
		hh, ll := high[i-momLen+1], low[i-momLen+1]
		for j := i - momLen + 2; j <= i; j++ {
			if high[j] > hh {
				hh = high[j]
			}
			if low[j] < ll {
				ll = low[j]
			}
		}
		mid := (hh + ll) / 2
		src[i] = close[i] - (mid+smaSrc[i])/2
	}
	momentum = LINEARREGFn(src, momLen)
	return
}

// CKSPFn — Chande Kroll Stop:
//
//	hAtr = high - mult * ATR(p_atr); lAtr = low + mult * ATR(p_atr)
//	long_stop  = highest(hAtr, p_max)
//	short_stop = lowest(lAtr,  p_max)
func CKSPFn(high, low, close []float64, atrPeriod int, mult float64, maxPeriod int) (longStop, shortStop []float64) {
	n := len(close)
	longStop = make([]float64, n)
	shortStop = make([]float64, n)
	if atrPeriod < 1 || maxPeriod < 1 || n <= atrPeriod {
		return
	}
	atr := ATRFn(high, low, close, atrPeriod)
	hAtr := make([]float64, n)
	lAtr := make([]float64, n)
	for i := 0; i < n; i++ {
		hAtr[i] = high[i] - mult*atr[i]
		lAtr[i] = low[i] + mult*atr[i]
	}
	for i := maxPeriod - 1; i < n; i++ {
		hh, ll := hAtr[i-maxPeriod+1], lAtr[i-maxPeriod+1]
		for j := i - maxPeriod + 2; j <= i; j++ {
			if hAtr[j] > hh {
				hh = hAtr[j]
			}
			if lAtr[j] < ll {
				ll = lAtr[j]
			}
		}
		longStop[i] = hh
		shortStop[i] = ll
	}
	return
}

// SINWMAFn — Sine Weighted MA. Weights w[k] = sin((k+1)·π / (p+1)).
func SINWMAFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	w := make([]float64, period)
	var ws float64
	for k := 0; k < period; k++ {
		w[k] = math.Sin(float64(k+1) * math.Pi / float64(period+1))
		ws += w[k]
	}
	if ws == 0 {
		return out
	}
	for i := period - 1; i < n; i++ {
		var v float64
		for k := 0; k < period; k++ {
			v += real[i-period+1+k] * w[k]
		}
		out[i] = v / ws
	}
	return out
}

// FVEFn — Finite Volume Element (Markos Katsanos / Pandas TA).
//
//	hlc3 = (h+l+c)/3
//	mf   = close - hlc3 + (hlc3 - hlc3[-1])
//	intra = ln(high) - ln(low)
//	cutoff = STDDEV(intra, period) * factor
//	flag = +1 if mf > cutoff; -1 if mf < -cutoff; else 0
//	FVE  = SUM(flag*volume, period) / (period * SMA(volume, period)) * 100
func FVEFn(high, low, close, volume []float64, period int, factor float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n < period+1 {
		return out
	}
	intra := make([]float64, n)
	for i := 0; i < n; i++ {
		if high[i] > 0 && low[i] > 0 {
			intra[i] = math.Log(high[i]) - math.Log(low[i])
		}
	}
	cutoff := STDDEVFn(intra, period, 1.0)
	for i := 0; i < n; i++ {
		cutoff[i] *= factor
	}
	flag := make([]float64, n)
	for i := 1; i < n; i++ {
		hlc3 := (high[i] + low[i] + close[i]) / 3
		hlc3Prev := (high[i-1] + low[i-1] + close[i-1]) / 3
		mf := close[i] - hlc3 + (hlc3 - hlc3Prev)
		switch {
		case mf > cutoff[i]:
			flag[i] = 1
		case mf < -cutoff[i]:
			flag[i] = -1
		}
	}
	fv := make([]float64, n)
	for i := 0; i < n; i++ {
		fv[i] = flag[i] * volume[i]
	}
	smaVol := SMAFn(volume, period)
	pf := float64(period)
	var sumFV float64
	for i := 0; i < period; i++ {
		sumFV += fv[i]
	}
	if smaVol[period-1] != 0 {
		out[period-1] = sumFV / (pf * smaVol[period-1]) * 100
	}
	for i := period; i < n; i++ {
		sumFV += fv[i] - fv[i-period]
		if smaVol[i] != 0 {
			out[i] = sumFV / (pf * smaVol[i]) * 100
		}
	}
	return out
}

// MANSFIELDFn — Mansfield Relative Strength: 100 * (a/b - SMA(a/b, p)) / SMA(a/b, p).
func MANSFIELDFn(a, b []float64, period int) []float64 {
	n := len(a)
	out := make([]float64, n)
	if len(b) < n || period < 1 || n < period {
		return out
	}
	r := make([]float64, n)
	for i := 0; i < n; i++ {
		if b[i] != 0 {
			r[i] = a[i] / b[i]
		}
	}
	sma := SMAFn(r, period)
	for i := 0; i < n; i++ {
		if sma[i] != 0 {
			out[i] = 100 * (r[i] - sma[i]) / sma[i]
		}
	}
	return out
}

// ICHIMOKUFn — Ichimoku Cloud (Pandas TA convention: per-bar array element
// represents what's plotted/visible at that bar).
//
//	tenkan[i]   = (HHV(high, tenkanP) + LLV(low, tenkanP)) / 2
//	kijun[i]    = (HHV(high, kijunP)  + LLV(low, kijunP))  / 2
//	senkou_a[i] = (tenkan[i-displacement] + kijun[i-displacement]) / 2
//	senkou_b[i] = ((HHV(high, senkouP) + LLV(low, senkouP)) / 2)[i-displacement]
//	chikou[i]   = close[i + displacement]    (0 when out of range)
func ICHIMOKUFn(high, low, close []float64, tenkanP, kijunP, senkouP, displacement int) (tenkan, kijun, senkouA, senkouB, chikou []float64) {
	n := len(close)
	tenkan = make([]float64, n)
	kijun = make([]float64, n)
	senkouA = make([]float64, n)
	senkouB = make([]float64, n)
	chikou = make([]float64, n)
	if n == 0 {
		return
	}
	hhll := func(period int) []float64 {
		out := make([]float64, n)
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
			out[i] = (hh + ll) / 2
		}
		return out
	}
	tenkan = hhll(tenkanP)
	kijun = hhll(kijunP)
	senkouBSrc := hhll(senkouP)
	for i := 0; i < n; i++ {
		if i-displacement >= 0 {
			senkouA[i] = (tenkan[i-displacement] + kijun[i-displacement]) / 2
			senkouB[i] = senkouBSrc[i-displacement]
		}
		if i+displacement < n {
			chikou[i] = close[i+displacement]
		}
	}
	return
}

// CPRFn — Central Pivot Range (per-bar using prior bar's HLC).
//
//	pivot[i] = (h[i-1] + l[i-1] + c[i-1]) / 3
//	bc[i]    = (h[i-1] + l[i-1]) / 2
//	tc[i]    = 2*pivot[i] - bc[i]
func CPRFn(high, low, close []float64) (pivot, tc, bc []float64) {
	n := len(close)
	pivot = make([]float64, n)
	tc = make([]float64, n)
	bc = make([]float64, n)
	for i := 1; i < n; i++ {
		p := (high[i-1] + low[i-1] + close[i-1]) / 3
		b := (high[i-1] + low[i-1]) / 2
		pivot[i] = p
		bc[i] = b
		tc[i] = 2*p - b
	}
	return
}

// FCBFn — Fractal Chaos Bands. Tracks the most recent Bill Williams up- and
// down-fractal high/low and carries the price level forward.
func FCBFn(high, low []float64) (upper, lower []float64) {
	n := len(high)
	upper = make([]float64, n)
	lower = make([]float64, n)
	up, dn := FRACTALFn(high, low)
	var lastUp, lastDn float64
	for i := 0; i < n; i++ {
		if up[i] == 1 {
			lastUp = high[i]
		}
		if dn[i] == 1 {
			lastDn = low[i]
		}
		upper[i] = lastUp
		lower[i] = lastDn
	}
	return
}

// WMAENVFn — WMA-based envelope: WMA(close, p) * (1 ± pct/100).
func WMAENVFn(close []float64, period int, pct float64) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	mid := WMAFn(close, period)
	f := pct / 100.0
	for i := 0; i < n; i++ {
		middle[i] = mid[i]
		upper[i] = mid[i] * (1 + f)
		lower[i] = mid[i] * (1 - f)
	}
	return
}

// MACDHISTFn — Just the MACD histogram (macd - signal). Convenience wrapper.
func MACDHISTFn(real []float64, fast, slow, signal int) []float64 {
	_, _, hist := MACDFn(real, fast, slow, signal)
	return hist
}

// ADSMOOTHFn — EMA-smoothed Accumulation/Distribution line.
func ADSMOOTHFn(high, low, close, volume []float64, period int) []float64 {
	return EMAFn(ADFn(high, low, close, volume), period)
}

// EMADIFFFn — EMA(real, fast) - EMA(real, slow). The MACD line absent the signal.
func EMADIFFFn(real []float64, fast, slow int) []float64 {
	n := len(real)
	out := make([]float64, n)
	f := EMAFn(real, fast)
	s := EMAFn(real, slow)
	for i := 0; i < n; i++ {
		out[i] = f[i] - s[i]
	}
	return out
}

// BBSQUEEZEFn — Boolean (1/0) Bollinger Squeeze: 1 when BBW(close,p,dev,dev)
// equals the rolling minimum BBW over `lookback` bars.
func BBSQUEEZEFn(close []float64, period int, dev float64, lookback int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || lookback < 1 || n < period+lookback {
		return out
	}
	bbw := BBWFn(close, period, dev, dev)
	for i := period + lookback - 2; i < n; i++ {
		mn := bbw[i-lookback+1]
		for j := i - lookback + 2; j <= i; j++ {
			if bbw[j] < mn {
				mn = bbw[j]
			}
		}
		if bbw[i] <= mn {
			out[i] = 1
		}
	}
	return out
}

// DONCHIANPCTFn — Donchian %B: (close - lower) / (upper - lower).
func DONCHIANPCTFn(high, low, close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	upper, _, lower := DONCHIANFn(high, low, period)
	for i := 0; i < n; i++ {
		w := upper[i] - lower[i]
		if w != 0 {
			out[i] = (close[i] - lower[i]) / w
		}
	}
	return out
}

// TRENDSCOREFn — Sum over the last `period` bars of sign(Δclose):
// +1 for up bars, -1 for down bars, 0 for unchanged. Range [-period, +period].
func TRENDSCOREFn(close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	dir := make([]float64, n)
	for i := 1; i < n; i++ {
		switch {
		case close[i] > close[i-1]:
			dir[i] = 1
		case close[i] < close[i-1]:
			dir[i] = -1
		}
	}
	var s float64
	for i := 1; i <= period; i++ {
		s += dir[i]
	}
	out[period] = s
	for i := period + 1; i < n; i++ {
		s += dir[i] - dir[i-period]
		out[i] = s
	}
	return out
}

// RANGEPCTFn — Bar range as a percentage of close: 100 * (high - low) / close.
func RANGEPCTFn(high, low, close []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		if close[i] != 0 {
			out[i] = 100 * (high[i] - low[i]) / close[i]
		}
	}
	return out
}

// TRUERANGEPCTFn — True Range as a percentage of close.
func TRUERANGEPCTFn(high, low, close []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	tr := trueRange(high, low, close)
	for i := 0; i < n; i++ {
		if close[i] != 0 {
			out[i] = 100 * tr[i] / close[i]
		}
	}
	return out
}

// CAMARILLAFn — Camarilla Pivot Points using prior bar's HLC. Returns
// pp + r1..r4 + s1..s4 in eight series via the helper.
func CAMARILLAFn(high, low, close []float64) (pp, r1, r2, r3, r4, s1, s2, s3, s4 []float64) {
	n := len(close)
	pp = make([]float64, n)
	r1 = make([]float64, n)
	r2 = make([]float64, n)
	r3 = make([]float64, n)
	r4 = make([]float64, n)
	s1 = make([]float64, n)
	s2 = make([]float64, n)
	s3 = make([]float64, n)
	s4 = make([]float64, n)
	for i := 1; i < n; i++ {
		ph, pl, pc := high[i-1], low[i-1], close[i-1]
		rng := ph - pl
		pp[i] = (ph + pl + pc) / 3
		r1[i] = pc + rng*1.1/12
		r2[i] = pc + rng*1.1/6
		r3[i] = pc + rng*1.1/4
		r4[i] = pc + rng*1.1/2
		s1[i] = pc - rng*1.1/12
		s2[i] = pc - rng*1.1/6
		s3[i] = pc - rng*1.1/4
		s4[i] = pc - rng*1.1/2
	}
	return
}

// WOODIEFn — Woodie Pivot Points using prior bar's HLC.
//
//	pp = (h + l + 2c) / 4
//	r1 = 2*pp - l;  r2 = pp + (h - l)
//	s1 = 2*pp - h;  s2 = pp - (h - l)
func WOODIEFn(high, low, close []float64) (pp, r1, r2, s1, s2 []float64) {
	n := len(close)
	pp = make([]float64, n)
	r1 = make([]float64, n)
	r2 = make([]float64, n)
	s1 = make([]float64, n)
	s2 = make([]float64, n)
	for i := 1; i < n; i++ {
		ph, pl, pc := high[i-1], low[i-1], close[i-1]
		p := (ph + pl + 2*pc) / 4
		pp[i] = p
		r1[i] = 2*p - pl
		s1[i] = 2*p - ph
		r2[i] = p + (ph - pl)
		s2[i] = p - (ph - pl)
	}
	return
}

// FIBPIVOTSFn — Fibonacci Pivot Points using prior bar's HLC.
//
//	pp = (h + l + c) / 3
//	r1/s1 = pp ± 0.382 * range
//	r2/s2 = pp ± 0.618 * range
//	r3/s3 = pp ± 1.000 * range
func FIBPIVOTSFn(high, low, close []float64) (pp, r1, r2, r3, s1, s2, s3 []float64) {
	n := len(close)
	pp = make([]float64, n)
	r1 = make([]float64, n)
	r2 = make([]float64, n)
	r3 = make([]float64, n)
	s1 = make([]float64, n)
	s2 = make([]float64, n)
	s3 = make([]float64, n)
	for i := 1; i < n; i++ {
		ph, pl, pc := high[i-1], low[i-1], close[i-1]
		p := (ph + pl + pc) / 3
		rng := ph - pl
		pp[i] = p
		r1[i] = p + 0.382*rng
		r2[i] = p + 0.618*rng
		r3[i] = p + 1.000*rng
		s1[i] = p - 0.382*rng
		s2[i] = p - 0.618*rng
		s3[i] = p - 1.000*rng
	}
	return
}

// DEMARKFn — DeMark Pivot Points using prior bar's OHLC.
//
//	x  = if c < o: 2h + l + c
//	     if c > o: h + 2l + c     (note: DeMark swaps the doubled bar)
//	     else:     h + l + 2c
//	pp = x / 4
//	r1 = x/2 - l;  s1 = x/2 - h
func DEMARKFn(open, high, low, close []float64) (pp, r1, s1 []float64) {
	n := len(close)
	pp = make([]float64, n)
	r1 = make([]float64, n)
	s1 = make([]float64, n)
	for i := 1; i < n; i++ {
		po, ph, pl, pc := open[i-1], high[i-1], low[i-1], close[i-1]
		var x float64
		switch {
		case pc < po:
			x = ph + 2*pl + pc
		case pc > po:
			x = 2*ph + pl + pc
		default:
			x = ph + pl + 2*pc
		}
		pp[i] = x / 4
		r1[i] = x/2 - pl
		s1[i] = x/2 - ph
	}
	return
}

// ATRBANDSFn — ATR Bands around an SMA mid: SMA(close, p) ± mult * ATR(h,l,c, p).
func ATRBANDSFn(high, low, close []float64, period int, mult float64) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	if period < 1 || n < period {
		return
	}
	mid := SMAFn(close, period)
	atr := ATRFn(high, low, close, period)
	for i := 0; i < n; i++ {
		middle[i] = mid[i]
		upper[i] = mid[i] + mult*atr[i]
		lower[i] = mid[i] - mult*atr[i]
	}
	return
}

// VORTEXDIFFFn — VI+ minus VI− over `period` bars.
func VORTEXDIFFFn(high, low, close []float64, period int) []float64 {
	vip, vim := VORTEXFn(high, low, close, period)
	out := make([]float64, len(vip))
	for i := range vip {
		out[i] = vip[i] - vim[i]
	}
	return out
}

// MFISMOOTHFn — EMA-smoothed Money Flow Index.
func MFISMOOTHFn(high, low, close, volume []float64, mfiPeriod, smoothPeriod int) []float64 {
	return EMAFn(MFIFn(high, low, close, volume, mfiPeriod), smoothPeriod)
}

// ADPCTFn — Period-over-period percent change of the AD line:
// 100 * (AD[i] - AD[i-period]) / |AD[i-period]|.
func ADPCTFn(high, low, close, volume []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	ad := ADFn(high, low, close, volume)
	for i := period; i < n; i++ {
		den := math.Abs(ad[i-period])
		if den != 0 {
			out[i] = 100 * (ad[i] - ad[i-period]) / den
		}
	}
	return out
}

// RANGEATRPCTFn — TR / ATR(h,l,c, period). >1 means today's range is larger
// than the recent average.
func RANGEATRPCTFn(high, low, close []float64, period int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	tr := trueRange(high, low, close)
	atr := ATRFn(high, low, close, period)
	for i := 0; i < n; i++ {
		if atr[i] != 0 {
			out[i] = tr[i] / atr[i]
		}
	}
	return out
}

// MACDZLHISTFn — Just the histogram from the Zero-Lag MACD.
func MACDZLHISTFn(real []float64, fast, slow, signal int) []float64 {
	_, _, hist := MACDZLFn(real, fast, slow, signal)
	return hist
}

// STOCHDIFFFn — slow K − slow D from STOCH (default 14/3-SMA/3-SMA).
func STOCHDIFFFn(high, low, close []float64, fastK, slowK, slowD int) []float64 {
	k, d := STOCHFn(high, low, close, fastK, slowK, SMA, slowD, SMA)
	out := make([]float64, len(k))
	for i := range k {
		out[i] = k[i] - d[i]
	}
	return out
}

// GAUSSIANFn — Ehlers' 4-pole Gaussian Filter.
//
//	β  = (1 - cos(2π/period)) / (√2 - 1)
//	α  = -β + √(β² + 2β)
//	out[i] = α⁴·real[i]
//	       + 4(1-α)·out[i-1]
//	       − 6(1-α)²·out[i-2]
//	       + 4(1-α)³·out[i-3]
//	       − (1-α)⁴·out[i-4]
func GAUSSIANFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < 4 {
		return out
	}
	beta := (1 - math.Cos(2*math.Pi/float64(period))) / (math.Sqrt2 - 1)
	alpha := -beta + math.Sqrt(beta*beta+2*beta)
	a4 := alpha * alpha * alpha * alpha
	one := 1 - alpha
	c1 := 4 * one
	c2 := 6 * one * one
	c3 := 4 * one * one * one
	c4 := one * one * one * one
	out[0] = real[0]
	out[1] = real[1]
	out[2] = real[2]
	out[3] = real[3]
	for i := 4; i < n; i++ {
		out[i] = a4*real[i] + c1*out[i-1] - c2*out[i-2] + c3*out[i-3] - c4*out[i-4]
	}
	return out
}

// HWMAFn — Holt-Winters MA (Pandas TA hwma) with smoothing constants
// na (level), nb (trend), nc (acceleration).
func HWMAFn(real []float64, na, nb, nc float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	var lastF, lastV, lastA float64
	out[0] = real[0]
	lastF = real[0]
	for i := 1; i < n; i++ {
		f := (1-na)*(lastF+lastV+0.5*lastA) + na*real[i]
		v := (1-nb)*(lastV+lastA) + nb*(f-lastF)
		a := (1-nc)*lastA + nc*(v-lastV)
		out[i] = f + v + 0.5*a
		lastF, lastV, lastA = f, v, a
	}
	return out
}

// HWCFn — Holt-Winters Channel: HWMA mid + scalar*STDDEV(close-mid, channel_period).
func HWCFn(close []float64, na, nb, nc float64, channelPeriod int, scalar float64) (upper, middle, lower []float64) {
	n := len(close)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	if channelPeriod < 1 || n < channelPeriod {
		return
	}
	mid := HWMAFn(close, na, nb, nc)
	err := make([]float64, n)
	for i := 0; i < n; i++ {
		err[i] = close[i] - mid[i]
	}
	sd := STDDEVFn(err, channelPeriod, 1.0)
	for i := 0; i < n; i++ {
		middle[i] = mid[i]
		upper[i] = mid[i] + scalar*sd[i]
		lower[i] = mid[i] - scalar*sd[i]
	}
	return
}

// MACDVFn — MACD-V (ATR-normalized MACD): (EMAfast-EMAslow) / ATR(slow) * 100.
// Returns macdv, signal, hist.
func MACDVFn(high, low, close []float64, fast, slow, signal int) (macdv, sig, hist []float64) {
	n := len(close)
	macdv = make([]float64, n)
	sig = make([]float64, n)
	hist = make([]float64, n)
	ef := EMAFn(close, fast)
	es := EMAFn(close, slow)
	atr := ATRFn(high, low, close, slow)
	for i := 0; i < n; i++ {
		if atr[i] != 0 {
			macdv[i] = (ef[i] - es[i]) / atr[i] * 100
		}
	}
	sig = EMAFn(macdv, signal)
	for i := 0; i < n; i++ {
		hist[i] = macdv[i] - sig[i]
	}
	return
}

// TDIFn — Trader's Dynamic Index. Returns rsi (raw), mab (fast SMA of RSI),
// mbl (medium SMA of RSI), upper/middle/lower volatility bands of RSI.
func TDIFn(close []float64, rsiPeriod, mabPeriod, mblPeriod, bandPeriod int, devMult float64) (rsi, mab, mbl, upper, middle, lower []float64) {
	n := len(close)
	rsi = RSIFn(close, rsiPeriod)
	mab = SMAFn(rsi, mabPeriod)
	mbl = SMAFn(rsi, mblPeriod)
	bandMid := SMAFn(rsi, bandPeriod)
	bandStd := STDDEVFn(rsi, bandPeriod, 1.0)
	upper = make([]float64, n)
	middle = make([]float64, n)
	lower = make([]float64, n)
	for i := 0; i < n; i++ {
		middle[i] = bandMid[i]
		upper[i] = bandMid[i] + devMult*bandStd[i]
		lower[i] = bandMid[i] - devMult*bandStd[i]
	}
	return
}

// THERMOFn — Elder Thermometer:
//
//	thermo[i] = max(|high[i] - high[i-1]|, |low[i-1] - low[i]|)
//	ema       = EMA(thermo, period)
//
// Returns the raw thermo and its EMA-smoothed version.
func THERMOFn(high, low []float64, period int) (thermo, smoothed []float64) {
	n := len(high)
	thermo = make([]float64, n)
	for i := 1; i < n; i++ {
		dh := math.Abs(high[i] - high[i-1])
		dl := math.Abs(low[i-1] - low[i])
		if dh > dl {
			thermo[i] = dh
		} else {
			thermo[i] = dl
		}
	}
	smoothed = EMAFn(thermo, period)
	return
}

// QQEFn — Quantitative Qualitative Estimation (simplified Pandas TA flavor).
// Returns the EMA-smoothed RSI (rsi_ma) and its dynamic ATR-of-RSI band (dar).
//
//	rsi_ma  = EMA(RSI(close, p), smooth)
//	tr_rsi  = |rsi_ma - rsi_ma[-1]|
//	dar     = EMA(EMA(tr_rsi, 2*p-1), 2*p-1) * factor
func QQEFn(close []float64, period, smooth int, factor float64) (rsiMA, dar []float64) {
	n := len(close)
	rsiMA = make([]float64, n)
	dar = make([]float64, n)
	if n < period+smooth {
		return
	}
	rsiMA = EMAFn(RSIFn(close, period), smooth)
	tr := make([]float64, n)
	for i := 1; i < n; i++ {
		tr[i] = math.Abs(rsiMA[i] - rsiMA[i-1])
	}
	w := 2*period - 1
	smoothed := EMAFn(EMAFn(tr, w), w)
	for i := 0; i < n; i++ {
		dar[i] = smoothed[i] * factor
	}
	return
}

// AOBVFn — Accumulation OBV (cinar): EMA-smoothed OBV plus a long/short
// crossover line via fast/slow EMAs.
func AOBVFn(close, volume []float64, fast, slow int) (long, short []float64) {
	n := len(close)
	obv := OBVFn(close, volume)
	long = EMAFn(obv, fast)
	short = EMAFn(obv, slow)
	if n == 0 {
		return
	}
	return
}

// VWAPANCHFn — Anchored VWAP that resets every `anchorBars` bars.
//
//	tp = (h+l+c)/3
//	cumulative price-volume / cumulative volume within each anchor window.
func VWAPANCHFn(high, low, close, volume []float64, anchorBars int) []float64 {
	n := len(close)
	out := make([]float64, n)
	if anchorBars < 1 {
		anchorBars = 1
	}
	var pv, vv float64
	for i := 0; i < n; i++ {
		if i%anchorBars == 0 {
			pv, vv = 0, 0
		}
		tp := (high[i] + low[i] + close[i]) / 3
		pv += tp * volume[i]
		vv += volume[i]
		if vv != 0 {
			out[i] = pv / vv
		}
	}
	return out
}

// ROCSIGNALFn — ROC and an EMA-smoothed signal line.
func ROCSIGNALFn(real []float64, period, signalPeriod int) (roc, sig []float64) {
	roc = ROCFn(real, period)
	if signalPeriod > 0 {
		sig = EMAFn(roc, signalPeriod)
	} else {
		sig = make([]float64, len(real))
	}
	return
}

// LINREGRESIDFn — Residuals of a rolling linear regression: real - LINREG(real, p).
func LINREGRESIDFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	lr := LINEARREGFn(real, period)
	for i := 0; i < n; i++ {
		out[i] = real[i] - lr[i]
	}
	return out
}

// VWAPPCTFn — Percentage deviation of close from the running VWAP.
func VWAPPCTFn(high, low, close, volume []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	v := VWAPFn(high, low, close, volume)
	for i := 0; i < n; i++ {
		if v[i] != 0 {
			out[i] = 100 * (close[i] - v[i]) / v[i]
		}
	}
	return out
}

// WMADIFFFn — WMA(real, fast) - WMA(real, slow). WMA-based MACD line.
func WMADIFFFn(real []float64, fast, slow int) []float64 {
	n := len(real)
	out := make([]float64, n)
	f := WMAFn(real, fast)
	s := WMAFn(real, slow)
	for i := 0; i < n; i++ {
		out[i] = f[i] - s[i]
	}
	return out
}

// CYBERCYCLEFn — Ehlers Cyber Cycle.
//
//	smooth[i] = (c[i] + 2c[i-1] + 2c[i-2] + c[i-3]) / 6
//	cycle[i]  = (1 - α/2)² * (smooth[i] - 2 smooth[i-1] + smooth[i-2])
//	          + 2(1 - α) * cycle[i-1] - (1 - α)² * cycle[i-2]
func CYBERCYCLEFn(real []float64, alpha float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n < 4 {
		return out
	}
	smooth := make([]float64, n)
	for i := 3; i < n; i++ {
		smooth[i] = (real[i] + 2*real[i-1] + 2*real[i-2] + real[i-3]) / 6
	}
	a2 := (1 - alpha/2) * (1 - alpha/2)
	o := 1 - alpha
	for i := 4; i < n; i++ {
		out[i] = a2*(smooth[i]-2*smooth[i-1]+smooth[i-2]) + 2*o*out[i-1] - o*o*out[i-2]
	}
	return out
}

// DECYCLERFn — Ehlers High-Pass Decycler. Returns the low-pass series:
// decycler = real - hp where hp is a 2-pole Butterworth high-pass filter.
func DECYCLERFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < 3 {
		return out
	}
	w := 2 * math.Pi / float64(period)
	alpha := (math.Cos(w) + math.Sin(w) - 1) / math.Cos(w)
	a2 := (1 - alpha/2) * (1 - alpha/2)
	o := 1 - alpha
	hp := make([]float64, n)
	for i := 2; i < n; i++ {
		hp[i] = a2*(real[i]-2*real[i-1]+real[i-2]) + 2*o*hp[i-1] - o*o*hp[i-2]
	}
	for i := 0; i < n; i++ {
		out[i] = real[i] - hp[i]
	}
	return out
}

// SWINGINDEXFn — Simplified Wilder Swing Index.
//
//	R   = max(|h - c[-1]|, |l - c[-1]|)
//	SI  = 50 * ((c - c[-1]) + 0.5*(c - o) + 0.25*(c[-1] - o[-1])) / R
//
// Returns 0 for the first bar and bars where R is 0. Wilder's original
// formula scales by K/T (limit-move constant); omitted here so the result
// is comparable across instruments without a futures-style limit move.
func SWINGINDEXFn(open, high, low, close []float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	for i := 1; i < n; i++ {
		dh := math.Abs(high[i] - close[i-1])
		dl := math.Abs(low[i] - close[i-1])
		R := dh
		if dl > R {
			R = dl
		}
		if R == 0 {
			continue
		}
		out[i] = 50 * ((close[i] - close[i-1]) + 0.5*(close[i]-open[i]) + 0.25*(close[i-1]-open[i-1])) / R
	}
	return out
}

// CSIFn — Commodity Selection Index (simplified Wilder): ADXR * ATR over
// `period` bars, scaled by `scalar` (defaults to 1; original uses
// V/√M / (150 + COMM)).
func CSIFn(high, low, close []float64, period int, scalar float64) []float64 {
	n := len(close)
	out := make([]float64, n)
	adxr := ADXRFn(high, low, close, period)
	atr := ATRFn(high, low, close, period)
	for i := 0; i < n; i++ {
		out[i] = adxr[i] * atr[i] * scalar
	}
	return out
}

// MOMSIGNALFn — Momentum with an EMA signal line.
func MOMSIGNALFn(real []float64, period, signalPeriod int) (mom, sig []float64) {
	mom = MOMFn(real, period)
	sig = EMAFn(mom, signalPeriod)
	return
}

// MFISIGNALFn — MFI with an EMA signal line.
func MFISIGNALFn(high, low, close, volume []float64, period, signalPeriod int) (mfi, sig []float64) {
	mfi = MFIFn(high, low, close, volume, period)
	sig = EMAFn(mfi, signalPeriod)
	return
}

// WILLRSIGNALFn — Williams %R with an EMA signal line.
func WILLRSIGNALFn(high, low, close []float64, period, signalPeriod int) (wr, sig []float64) {
	wr = WILLRFn(high, low, close, period)
	sig = EMAFn(wr, signalPeriod)
	return
}

// CCISIGNALFn — CCI with an EMA signal line.
func CCISIGNALFn(high, low, close []float64, period, signalPeriod int) (cci, sig []float64) {
	cci = CCIFn(high, low, close, period)
	sig = EMAFn(cci, signalPeriod)
	return
}

// OBVSIGNALFn — OBV with an EMA signal line.
func OBVSIGNALFn(close, volume []float64, signalPeriod int) (obv, sig []float64) {
	obv = OBVFn(close, volume)
	sig = EMAFn(obv, signalPeriod)
	return
}

// CMFSIGNALFn — CMF with an EMA signal line.
func CMFSIGNALFn(high, low, close, volume []float64, period, signalPeriod int) (cmf, sig []float64) {
	cmf = CMFFn(high, low, close, volume, period)
	sig = EMAFn(cmf, signalPeriod)
	return
}

// SLOPEPCTFn — LINREG slope expressed as a percentage of the LINREG value.
func SLOPEPCTFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	slope := LINEARREGSLOPEFn(real, period)
	lr := LINEARREGFn(real, period)
	for i := 0; i < n; i++ {
		if lr[i] != 0 {
			out[i] = 100 * slope[i] / lr[i]
		}
	}
	return out
}

// TPSMOOTHFn — EMA-smoothed Typical Price (HLC3).
func TPSMOOTHFn(high, low, close []float64, period int) []float64 {
	return EMAFn(HLC3Fn(high, low, close), period)
}

// PWOFn — Percent Williams Oscillator (Pandas TA):
// 100 * (WMA(real, fast) - WMA(real, slow)) / WMA(real, slow).
func PWOFn(real []float64, fast, slow int) []float64 {
	n := len(real)
	out := make([]float64, n)
	f := WMAFn(real, fast)
	s := WMAFn(real, slow)
	for i := 0; i < n; i++ {
		if s[i] != 0 {
			out[i] = 100 * (f[i] - s[i]) / s[i]
		}
	}
	return out
}

// TRANGESMOOTHFn — EMA-smoothed True Range.
func TRANGESMOOTHFn(high, low, close []float64, period int) []float64 {
	return EMAFn(trueRange(high, low, close), period)
}

// WCPSMOOTHFn — EMA-smoothed Weighted Close Price (h+l+2c)/4.
func WCPSMOOTHFn(high, low, close []float64, period int) []float64 {
	n := len(close)
	wcp := make([]float64, n)
	for i := 0; i < n; i++ {
		wcp[i] = (high[i] + low[i] + 2*close[i]) / 4
	}
	return EMAFn(wcp, period)
}

// MEDPRICESMOOTHFn — EMA-smoothed Median Price (h+l)/2.
func MEDPRICESMOOTHFn(high, low []float64, period int) []float64 {
	n := len(high)
	med := make([]float64, n)
	for i := 0; i < n; i++ {
		med[i] = (high[i] + low[i]) / 2
	}
	return EMAFn(med, period)
}

// ADZSCOREFn — rolling Z-score of the AD line.
func ADZSCOREFn(high, low, close, volume []float64, period int) []float64 {
	return ZSCOREFn(ADFn(high, low, close, volume), period)
}

// OBVZSCOREFn — rolling Z-score of OBV.
func OBVZSCOREFn(close, volume []float64, period int) []float64 {
	return ZSCOREFn(OBVFn(close, volume), period)
}

// RETURNZSCOREFn — rolling Z-score of one-bar percent returns.
func RETURNZSCOREFn(real []float64, period int) []float64 {
	return ZSCOREFn(PCTRETFn(real, 1), period)
}

// KVOPCTFn — KVO expressed as a percentage of |EMA(volume_force, slow)|.
func KVOPCTFn(high, low, close, volume []float64, fast, slow, signal int) []float64 {
	n := len(close)
	kvo, _ := KVOFn(high, low, close, volume, fast, slow, signal)
	out := make([]float64, n)
	// Use rolling normalisation by the absolute KVO mean over the slow window.
	abs := make([]float64, n)
	for i := 0; i < n; i++ {
		abs[i] = math.Abs(kvo[i])
	}
	scale := SMAFn(abs, slow)
	for i := 0; i < n; i++ {
		if scale[i] != 0 {
			out[i] = 100 * kvo[i] / scale[i]
		}
	}
	return out
}

// MACDPCTFn — MACD divided by close, in percent: 100 * MACD / close.
func MACDPCTFn(real []float64, fast, slow, signal int) []float64 {
	n := len(real)
	out := make([]float64, n)
	macd, _, _ := MACDFn(real, fast, slow, signal)
	for i := 0; i < n; i++ {
		if real[i] != 0 {
			out[i] = 100 * macd[i] / real[i]
		}
	}
	return out
}

// PPOSIGNALFn — PPO with an EMA-smoothed signal line.
func PPOSIGNALFn(real []float64, fast, slow, signal int, t MaType) (ppo, sig []float64) {
	ppo = PPOFn(real, fast, slow, t)
	sig = EMAFn(ppo, signal)
	return
}

// BBPSIGNALFn — BBP with an EMA-smoothed signal line.
func BBPSIGNALFn(real []float64, period int, devUp, devDn float64, signal int) (bbp, sig []float64) {
	bbp = BBPFn(real, period, devUp, devDn)
	sig = EMAFn(bbp, signal)
	return
}

// ADXSIGNALFn — ADX with an EMA-smoothed signal line.
func ADXSIGNALFn(high, low, close []float64, period, signal int) (adx, sig []float64) {
	adx = ADXFn(high, low, close, period)
	sig = EMAFn(adx, signal)
	return
}

// VORTEXFn — Vortex Indicator VI+ and VI− over `period` bars.
//
//	VM+ = |H_t - L_{t-1}|,  VM- = |L_t - H_{t-1}|
//	VI+ = SUM(VM+, p) / SUM(TR, p)
//	VI- = SUM(VM-, p) / SUM(TR, p)
func VORTEXFn(high, low, close []float64, period int) (vip, vim []float64) {
	n := len(high)
	vip = make([]float64, n)
	vim = make([]float64, n)
	if period < 1 || n <= period {
		return
	}
	tr := trueRange(high, low, close)
	vmp := make([]float64, n)
	vmm := make([]float64, n)
	for i := 1; i < n; i++ {
		vmp[i] = math.Abs(high[i] - low[i-1])
		vmm[i] = math.Abs(low[i] - high[i-1])
	}
	var sumP, sumM, sumT float64
	for i := 1; i <= period; i++ {
		sumP += vmp[i]
		sumM += vmm[i]
		sumT += tr[i]
	}
	if sumT != 0 {
		vip[period] = sumP / sumT
		vim[period] = sumM / sumT
	}
	for i := period + 1; i < n; i++ {
		sumP += vmp[i] - vmp[i-period]
		sumM += vmm[i] - vmm[i-period]
		sumT += tr[i] - tr[i-period]
		if sumT != 0 {
			vip[i] = sumP / sumT
			vim[i] = sumM / sumT
		}
	}
	return
}
