package talib

import "math"

// VARFn — population variance over rolling period.
func VARFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		var sum, sum2 float64
		for j := i - period + 1; j <= i; j++ {
			sum += real[j]
			sum2 += real[j] * real[j]
		}
		mean := sum / float64(period)
		out[i] = sum2/float64(period) - mean*mean
	}
	return out
}

// STDDEVFn — standard deviation over rolling period, scaled by nbDev.
func STDDEVFn(real []float64, period int, nbDev float64) []float64 {
	v := VARFn(real, period)
	out := make([]float64, len(v))
	for i, x := range v {
		if x > 0 {
			out[i] = math.Sqrt(x) * nbDev
		}
	}
	return out
}

// linRegRaw computes the slope, intercept, and r at the last index of a window.
// Window indices are 0..period-1 with x = j+1 (1-based) following TA-Lib semantics.
func linRegStats(real []float64, end, period int) (slope, intercept float64) {
	pf := float64(period)
	sumX := pf * (pf + 1) / 2
	sumXSqr := pf * (pf + 1) * (2*pf + 1) / 6
	denom := pf*sumXSqr - sumX*sumX

	var sumY, sumXY float64
	for j := 0; j < period; j++ {
		x := float64(j + 1)
		y := real[end-period+1+j]
		sumY += y
		sumXY += x * y
	}
	slope = (pf*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / pf
	return
}

// LINEARREGFn — value of linear regression line at the most recent point.
func LINEARREGFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		slope, intercept := linRegStats(real, i, period)
		out[i] = intercept + slope*float64(period)
	}
	return out
}

// LINEARREGSLOPEFn — slope of the regression line.
func LINEARREGSLOPEFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		slope, _ := linRegStats(real, i, period)
		out[i] = slope
	}
	return out
}

// LINEARREGANGLEFn — slope expressed in degrees.
func LINEARREGANGLEFn(real []float64, period int) []float64 {
	s := LINEARREGSLOPEFn(real, period)
	out := make([]float64, len(s))
	for i, v := range s {
		out[i] = math.Atan(v) * 180 / math.Pi
	}
	return out
}

// LINEARREGINTERCEPTFn — y-intercept of the regression line.
func LINEARREGINTERCEPTFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		_, intercept := linRegStats(real, i, period)
		out[i] = intercept
	}
	return out
}

// TSFFn — Time Series Forecast: regression value extrapolated one bar ahead.
func TSFFn(real []float64, period int) []float64 {
	n := len(real)
	out := make([]float64, n)
	if period < 2 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		slope, intercept := linRegStats(real, i, period)
		out[i] = intercept + slope*float64(period+1)
	}
	return out
}

// CORRELFn — Pearson correlation of two series over rolling period.
func CORRELFn(real0, real1 []float64, period int) []float64 {
	n := len(real0)
	out := make([]float64, n)
	if period < 1 || n < period {
		return out
	}
	for i := period - 1; i < n; i++ {
		var sx, sy, sxx, syy, sxy float64
		for j := i - period + 1; j <= i; j++ {
			x, y := real0[j], real1[j]
			sx += x
			sy += y
			sxx += x * x
			syy += y * y
			sxy += x * y
		}
		pf := float64(period)
		cov := sxy/pf - (sx/pf)*(sy/pf)
		varx := sxx/pf - (sx/pf)*(sx/pf)
		vary := syy/pf - (sy/pf)*(sy/pf)
		denom := math.Sqrt(varx * vary)
		if denom > 0 {
			out[i] = cov / denom
		}
	}
	return out
}

// BETAFn — slope of regression of real0 on real1, computed from per-bar returns.
// Following TA-Lib: uses pct change r_t = x_t/x_{t-1} - 1; window size = period.
func BETAFn(real0, real1 []float64, period int) []float64 {
	n := len(real0)
	out := make([]float64, n)
	if period < 1 || n <= period {
		return out
	}
	// Pre-compute returns into temporary slices aligned to index 1..n-1.
	r0 := make([]float64, n)
	r1 := make([]float64, n)
	for i := 1; i < n; i++ {
		if real0[i-1] != 0 {
			r0[i] = real0[i]/real0[i-1] - 1
		}
		if real1[i-1] != 0 {
			r1[i] = real1[i]/real1[i-1] - 1
		}
	}
	for i := period; i < n; i++ {
		var sx, sy, sxx, sxy float64
		for j := i - period + 1; j <= i; j++ {
			x, y := r1[j], r0[j]
			sx += x
			sy += y
			sxx += x * x
			sxy += x * y
		}
		pf := float64(period)
		denom := pf*sxx - sx*sx
		if denom != 0 {
			out[i] = (pf*sxy - sx*sy) / denom
		}
	}
	return out
}
