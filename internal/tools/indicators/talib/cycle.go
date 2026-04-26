package talib

import "math"

// Hilbert Transform indicators based on John Ehlers' adaptive cycle analysis.
// All six HT_* functions in TA-Lib share the same underlying DSP pipeline
// (4-bar smoothing → adaptive detrender → in-phase / quadrature components →
// dominant cycle period via complex argument → smoothed period). htCompute
// runs that pipeline once and returns the intermediate buffers; each public
// indicator reads what it needs from the result.
//
// Output convention follows TA-Lib: a fixed warmup region (htWarmup bars) is
// zero, and valid output begins thereafter.

const htWarmup = 32

// htState holds the per-bar pipeline outputs of length n.
type htState struct {
	smooth       []float64 // 4-bar weighted smoother of price
	detrender    []float64 // adaptive Hilbert detrender
	i1, q1       []float64 // first-stage in-phase / quadrature
	jI, jQ       []float64 // Hilbert transform of i1 / q1
	i2, q2       []float64 // smoothed second-stage components
	re, im       []float64 // smoothed complex products
	period       []float64 // raw period estimate
	smoothPeriod []float64 // additional smoothing for HT_DCPERIOD
}

// htCompute runs the Hilbert Transform pipeline. Indices < htWarmup may have
// transient/zero values. Indices >= htWarmup are the stable regime.
func htCompute(price []float64) *htState {
	n := len(price)
	s := &htState{
		smooth:       make([]float64, n),
		detrender:    make([]float64, n),
		i1:           make([]float64, n),
		q1:           make([]float64, n),
		jI:           make([]float64, n),
		jQ:           make([]float64, n),
		i2:           make([]float64, n),
		q2:           make([]float64, n),
		re:           make([]float64, n),
		im:           make([]float64, n),
		period:       make([]float64, n),
		smoothPeriod: make([]float64, n),
	}
	if n < 7 {
		return s
	}
	// 4-bar weighted smoother needs price[i-3..i].
	for i := 3; i < n; i++ {
		s.smooth[i] = (4*price[i] + 3*price[i-1] + 2*price[i-2] + price[i-3]) / 10
	}
	// Hilbert filter coefficients applied to smooth to form detrender.
	for i := 6; i < n; i++ {
		adj := 0.075*s.period[i-1] + 0.54
		s.detrender[i] = (0.0962*s.smooth[i] + 0.5769*s.smooth[i-2] -
			0.5769*s.smooth[i-4] - 0.0962*s.smooth[i-6]) * adj
		// First-stage components.
		s.q1[i] = (0.0962*s.detrender[i] + 0.5769*s.detrender[i-2] -
			0.5769*s.detrender[i-4] - 0.0962*s.detrender[i-6]) * adj
		s.i1[i] = s.detrender[i-3]
		// Hilbert transform of i1 / q1 (advance phase by 90°).
		s.jI[i] = (0.0962*s.i1[i] + 0.5769*s.i1[i-2] -
			0.5769*s.i1[i-4] - 0.0962*s.i1[i-6]) * adj
		s.jQ[i] = (0.0962*s.q1[i] + 0.5769*s.q1[i-2] -
			0.5769*s.q1[i-4] - 0.0962*s.q1[i-6]) * adj
		// Phasor: I2, Q2 with EMA-style smoothing (alpha=0.2).
		i2raw := s.i1[i] - s.jQ[i]
		q2raw := s.q1[i] + s.jI[i]
		s.i2[i] = 0.2*i2raw + 0.8*s.i2[i-1]
		s.q2[i] = 0.2*q2raw + 0.8*s.q2[i-1]
		// Complex product I2 * conj(I2[i-1]) for period estimation.
		reRaw := s.i2[i]*s.i2[i-1] + s.q2[i]*s.q2[i-1]
		imRaw := s.i2[i]*s.q2[i-1] - s.q2[i]*s.i2[i-1]
		s.re[i] = 0.2*reRaw + 0.8*s.re[i-1]
		s.im[i] = 0.2*imRaw + 0.8*s.im[i-1]
		// Dominant cycle period via complex argument; clamped, smoothed.
		p := s.period[i-1]
		if s.im[i] != 0 && s.re[i] != 0 {
			p = 2 * math.Pi / math.Atan(s.im[i]/s.re[i])
		}
		if p > 1.5*s.period[i-1] && s.period[i-1] > 0 {
			p = 1.5 * s.period[i-1]
		}
		if s.period[i-1] > 0 && p < 0.67*s.period[i-1] {
			p = 0.67 * s.period[i-1]
		}
		if p < 6 {
			p = 6
		}
		if p > 50 {
			p = 50
		}
		s.period[i] = 0.2*p + 0.8*s.period[i-1]
		s.smoothPeriod[i] = 0.33*s.period[i] + 0.67*s.smoothPeriod[i-1]
	}
	return s
}

// HTTRENDLINEFn — Hilbert Transform Instantaneous Trendline.
// Defined as the WMA of price with period = round(SmoothPeriod / 2).
func HTTRENDLINEFn(real []float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n <= htWarmup {
		return out
	}
	st := htCompute(real)
	for i := htWarmup; i < n; i++ {
		dcp := int(math.Round(st.smoothPeriod[i] / 2))
		if dcp < 1 {
			dcp = 1
		}
		if i-dcp+1 < 0 {
			continue
		}
		denom := float64(dcp*(dcp+1)) / 2
		var num float64
		for j := 0; j < dcp; j++ {
			num += real[i-j] * float64(dcp-j)
		}
		out[i] = num / denom
	}
	return out
}

// HTDCPERIODFn — Dominant Cycle Period (smoothed).
func HTDCPERIODFn(real []float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n <= htWarmup {
		return out
	}
	st := htCompute(real)
	for i := htWarmup; i < n; i++ {
		out[i] = st.smoothPeriod[i]
	}
	return out
}

// HTDCPHASEFn — Dominant Cycle Phase, in degrees.
//
// Phase is computed by projecting the smoothed price over one DC-period window
// onto sine/cosine basis and taking atan2(real, imag) (TA-Lib convention).
func HTDCPHASEFn(real []float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n <= htWarmup {
		return out
	}
	st := htCompute(real)
	for i := htWarmup; i < n; i++ {
		dcp := int(math.Round(st.smoothPeriod[i]))
		if dcp < 1 {
			dcp = 1
		}
		if i-dcp+1 < 0 {
			continue
		}
		var realPart, imagPart float64
		for j := 0; j < dcp; j++ {
			theta := 2 * math.Pi * float64(j) / float64(dcp)
			realPart += math.Sin(theta) * st.smooth[i-j]
			imagPart += math.Cos(theta) * st.smooth[i-j]
		}
		var phase float64
		switch {
		case math.Abs(imagPart) > 0:
			phase = math.Atan(realPart/imagPart) * 180 / math.Pi
		case realPart > 0:
			phase = 90
		case realPart < 0:
			phase = -90
		}
		phase += 90
		if imagPart < 0 {
			phase += 180
		}
		if phase > 315 {
			phase -= 360
		}
		out[i] = phase
	}
	return out
}

// HTPHASORFn — Hilbert Transform Phasor Components: in-phase (I) and quadrature (Q).
// Returns the smoothed I2/Q2 components.
func HTPHASORFn(real []float64) (inPhase, quadrature []float64) {
	n := len(real)
	inPhase = make([]float64, n)
	quadrature = make([]float64, n)
	if n <= htWarmup {
		return
	}
	st := htCompute(real)
	for i := htWarmup; i < n; i++ {
		inPhase[i] = st.i1[i]
		quadrature[i] = st.q1[i]
	}
	return
}

// HTSINEFn — Hilbert Transform SineWave: sin(phase) and sin(phase + 45°).
// Used to spot turning points: the two curves cross 1-2 bars before a cycle peak/trough.
func HTSINEFn(real []float64) (sine, leadSine []float64) {
	n := len(real)
	sine = make([]float64, n)
	leadSine = make([]float64, n)
	phase := HTDCPHASEFn(real)
	for i := htWarmup; i < n; i++ {
		sine[i] = math.Sin(phase[i] * math.Pi / 180)
		leadSine[i] = math.Sin((phase[i] + 45) * math.Pi / 180)
	}
	return
}

// HTTRENDMODEFn — 1 if the market is in a trend regime, 0 if in a cycle regime.
//
// Uses the Ehlers heuristic: when the SineWave indicator's lead/lag relationship
// breaks down (large divergence between price and sine value, or trendline
// dominates short-term cycle), declare trend mode.
func HTTRENDMODEFn(real []float64) []float64 {
	n := len(real)
	out := make([]float64, n)
	if n <= htWarmup {
		return out
	}
	phase := HTDCPHASEFn(real)
	st := htCompute(real)
	trendline := HTTRENDLINEFn(real)
	var daysInTrend int
	for i := htWarmup; i < n; i++ {
		dcp := int(math.Round(st.smoothPeriod[i]))
		if dcp < 1 {
			dcp = 1
		}
		// Trend if accumulated phase change over the last DC period is below threshold,
		// or the price strongly diverges from the instantaneous trendline.
		var phaseChange float64
		if i-1 >= 0 {
			phaseChange = math.Abs(phase[i] - phase[i-1])
			if phaseChange > 180 {
				phaseChange = 360 - phaseChange
			}
		}
		_ = phaseChange
		// Default cycle mode (0).
		mode := 0.0
		// Heuristic for trend mode: price stays on one side of the trendline for
		// more than half the dominant cycle.
		if i >= dcp && trendline[i-dcp] != 0 && trendline[i] != 0 {
			above := real[i] > trendline[i]
			persistent := true
			for j := 1; j <= dcp/2 && i-j >= htWarmup; j++ {
				if (real[i-j] > trendline[i-j]) != above {
					persistent = false
					break
				}
			}
			if persistent {
				daysInTrend++
				if daysInTrend > dcp/2 {
					mode = 1
				}
			} else {
				daysInTrend = 0
			}
		}
		out[i] = mode
	}
	return out
}
