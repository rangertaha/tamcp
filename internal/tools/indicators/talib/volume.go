package talib

// OBVFn — On Balance Volume. OBV[0]=volume[0]; if close[i]>close[i-1] add volume, < subtract, == carry.
func OBVFn(real, volume []float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	out[0] = volume[0]
	for i := 1; i < n; i++ {
		switch {
		case real[i] > real[i-1]:
			out[i] = out[i-1] + volume[i]
		case real[i] < real[i-1]:
			out[i] = out[i-1] - volume[i]
		default:
			out[i] = out[i-1]
		}
	}
	return out
}

// ADFn — Chaikin Accumulation/Distribution Line.
// MFM = ((close-low)-(high-close))/(high-low); MFV = MFM * volume; AD = cumsum(MFV).
func ADFn(high, low, close, volume []float64) []float64 {
	n := len(high)
	out := make([]float64, n)
	var ad float64
	for i := 0; i < n; i++ {
		hl := high[i] - low[i]
		if hl > 0 {
			mfm := ((close[i] - low[i]) - (high[i] - close[i])) / hl
			ad += mfm * volume[i]
		}
		out[i] = ad
	}
	return out
}

// ADOSCFn — Chaikin A/D Oscillator: EMA(AD,fast) - EMA(AD,slow).
func ADOSCFn(high, low, close, volume []float64, fastPeriod, slowPeriod int) []float64 {
	ad := ADFn(high, low, close, volume)
	fast := EMAFn(ad, fastPeriod)
	slow := EMAFn(ad, slowPeriod)
	n := len(ad)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		// Both fast and slow are zero in their warmup region; result valid where slow != 0.
		if slow[i] != 0 || fast[i] != 0 {
			out[i] = fast[i] - slow[i]
		}
	}
	return out
}
