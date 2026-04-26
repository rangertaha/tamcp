package talib

import (
	"fmt"
	"math"
)

// BSTSResult holds the fitted Bayesian Structural Time Series (local linear
// trend) parameters and the one-step-ahead forecast.
//
// State-space model:
//
//	level_{t+1}  = level_t + slope_t + ε_l    ε_l ~ N(0, σ²_level)
//	slope_{t+1}  = slope_t + ε_s             ε_s ~ N(0, σ²_slope)
//	y_t          = level_t + ε_y             ε_y ~ N(0, σ²_obs)
//
// (BSTS in the strict sense uses MCMC for posterior sampling. This is the
// frequentist Kalman-filter / ML estimate of the variance components, which
// gives the same point forecasts as the posterior mean of a vague-prior BSTS.)
type BSTSResult struct {
	// Variance components (optimised by ML).
	SigmaObs   float64 `json:"sigma_obs"`
	SigmaLevel float64 `json:"sigma_level"`
	SigmaSlope float64 `json:"sigma_slope"`
	// Final filtered state mean: [level_T, slope_T].
	Level float64 `json:"level"`
	Slope float64 `json:"slope"`
	// One-step forecast.
	NextValue  float64 `json:"next_value"`
	NextStdDev float64 `json:"next_stddev"`
	// Smoothed level and slope series (length = len(values)).
	LevelSeries []float64 `json:"level_series"`
	SlopeSeries []float64 `json:"slope_series"`
	// Log-likelihood at the optimum.
	LogLikelihood float64 `json:"log_likelihood"`
	// Iterations Nelder-Mead steps actually taken.
	Iterations int `json:"iterations"`
}

// BSTS fits a local-linear-trend BSTS model to the input series via Kalman-
// filter ML estimation of the three variance components, then returns the
// filtered state, the per-bar level/slope series, and a one-step forecast.
// Defaults: maxIter = 200.
func BSTS(values []float64, maxIter int) BSTSResult {
	out := BSTSResult{}
	n := len(values)
	if n < 10 {
		return out
	}
	if maxIter <= 0 {
		maxIter = 200
	}

	// Sample variance for scaling the search.
	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(n)
	V := 0.0
	for _, v := range values {
		V += (v - mean) * (v - mean)
	}
	V /= float64(n)
	if V < 1e-12 {
		V = 1e-12
	}

	// Optimise log-variances unconstrained: σ² = exp(θ) keeps σ² > 0.
	negLL := func(theta [3]float64) float64 {
		s2obs := math.Exp(theta[0])
		s2lev := math.Exp(theta[1])
		s2slo := math.Exp(theta[2])
		_, _, ll, _, _ := kalmanLLT(values, s2obs, s2lev, s2slo)
		return -ll
	}

	// 3-D Nelder-Mead in θ = log(σ²).
	logV := math.Log(V)
	x := [4][3]float64{
		{logV, logV - 4, logV - 6},
		{logV + 1, logV - 4, logV - 6},
		{logV, logV - 3, logV - 6},
		{logV, logV - 4, logV - 5},
	}
	f := [4]float64{negLL(x[0]), negLL(x[1]), negLL(x[2]), negLL(x[3])}

	const tol = 1e-7
	var iters int
	for iters = 0; iters < maxIter; iters++ {
		// Sort.
		for i := 0; i < 4; i++ {
			for j := i + 1; j < 4; j++ {
				if f[j] < f[i] {
					f[i], f[j] = f[j], f[i]
					x[i], x[j] = x[j], x[i]
				}
			}
		}
		if math.Abs(f[3]-f[0]) < tol {
			break
		}
		// Centroid of best three.
		var c [3]float64
		for k := 0; k < 3; k++ {
			c[k] = (x[0][k] + x[1][k] + x[2][k]) / 3
		}
		// Reflect.
		var xr [3]float64
		for k := 0; k < 3; k++ {
			xr[k] = 2*c[k] - x[3][k]
		}
		fr := negLL(xr)
		switch {
		case fr < f[0]:
			var xe [3]float64
			for k := 0; k < 3; k++ {
				xe[k] = c[k] + 2*(c[k]-x[3][k])
			}
			fe := negLL(xe)
			if fe < fr {
				x[3], f[3] = xe, fe
			} else {
				x[3], f[3] = xr, fr
			}
		case fr < f[2]:
			x[3], f[3] = xr, fr
		default:
			var xc [3]float64
			for k := 0; k < 3; k++ {
				xc[k] = c[k] + 0.5*(x[3][k]-c[k])
			}
			fc := negLL(xc)
			if fc < f[3] {
				x[3], f[3] = xc, fc
			} else {
				for i := 1; i < 4; i++ {
					for k := 0; k < 3; k++ {
						x[i][k] = x[0][k] + 0.5*(x[i][k]-x[0][k])
					}
					f[i] = negLL(x[i])
				}
			}
		}
	}

	best := x[0]
	s2obs := math.Exp(best[0])
	s2lev := math.Exp(best[1])
	s2slo := math.Exp(best[2])

	level, slope, ll, levelSeries, slopeSeries := kalmanLLT(values, s2obs, s2lev, s2slo)
	nextLevel := level + slope
	// Predictive variance for the next observation: P_pred = F·P·F' + Q + σ²_obs.
	// Approximation: use s2obs + s2lev + s2slo as a rough scale.
	nextSD := math.Sqrt(s2obs + s2lev + s2slo)

	out.SigmaObs = math.Sqrt(s2obs)
	out.SigmaLevel = math.Sqrt(s2lev)
	out.SigmaSlope = math.Sqrt(s2slo)
	out.Level = level
	out.Slope = slope
	out.NextValue = nextLevel
	out.NextStdDev = nextSD
	out.LevelSeries = levelSeries
	out.SlopeSeries = slopeSeries
	out.LogLikelihood = ll
	out.Iterations = iters
	return out
}

// kalmanLLT runs the Kalman filter for the 2-state local-linear-trend model
// and returns the final filtered (level, slope), the log-likelihood, and the
// per-bar filtered level/slope series.
func kalmanLLT(y []float64, s2obs, s2lev, s2slo float64) (float64, float64, float64, []float64, []float64) {
	n := len(y)
	// State x = [level, slope]. Transition F = [[1,1],[0,1]]. Observation H = [1, 0].
	// P is 2×2 covariance.
	level := y[0]
	slope := 0.0
	// Diffuse prior.
	P := [2][2]float64{{1e6, 0}, {0, 1e6}}

	levelSeries := make([]float64, n)
	slopeSeries := make([]float64, n)
	ll := 0.0

	for t := 0; t < n; t++ {
		// Predict.
		predLevel := level + slope
		predSlope := slope
		// P' = F·P·F' + Q.
		// F·P:
		fp00 := P[0][0] + P[1][0]
		fp01 := P[0][1] + P[1][1]
		fp10 := P[1][0]
		fp11 := P[1][1]
		// (F·P)·F':
		ppP00 := fp00 + fp01 // F'_{0,*} = [1, 0; 1, 1]^T applied
		ppP01 := fp01
		ppP10 := fp10 + fp11
		ppP11 := fp11
		// + Q:
		ppP00 += s2lev
		ppP11 += s2slo

		// Innovation.
		v := y[t] - predLevel // y - H·x'
		s := ppP00 + s2obs    // H·P'·H' + R
		if s < 1e-12 {
			s = 1e-12
		}
		// Kalman gain K = P'·H' / s = [P'_{0,0}, P'_{1,0}] / s.
		k0 := ppP00 / s
		k1 := ppP10 / s

		level = predLevel + k0*v
		slope = predSlope + k1*v
		// P = (I - K·H)·P'.
		newP00 := (1-k0)*ppP00 - 0*ppP10 // since H = [1,0]
		newP01 := (1-k0)*ppP01 - 0*ppP11
		newP10 := -k1*ppP00 + ppP10
		newP11 := -k1*ppP01 + ppP11
		P[0][0], P[0][1], P[1][0], P[1][1] = newP00, newP01, newP10, newP11

		levelSeries[t] = level
		slopeSeries[t] = slope

		ll += -0.5 * (math.Log(2*math.Pi*s) + v*v/s)
	}
	return level, slope, ll, levelSeries, slopeSeries
}

// BSTSSummary is a one-line text summary for a BSTSResult.
func BSTSSummary(r BSTSResult) string {
	if len(r.LevelSeries) == 0 {
		return "bsts: insufficient data"
	}
	return fmt.Sprintf(
		"bsts: level=%.4f slope=%.5f σ_obs=%.4f σ_level=%.4f σ_slope=%.4f next=%.4f (±%.4f), %d iters, ll=%.3f",
		r.Level, r.Slope, r.SigmaObs, r.SigmaLevel, r.SigmaSlope, r.NextValue, r.NextStdDev, r.Iterations, r.LogLikelihood,
	)
}
