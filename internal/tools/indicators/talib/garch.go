package talib

import (
	"fmt"
	"math"
)

// GARCHResult holds the fitted GARCH(1,1) parameters, the per-bar conditional
// variance/stddev series, and the one-step-ahead forecast.
//
// Model:
//
//	rₜ = μ + εₜ
//	εₜ = σₜ · zₜ          zₜ ~ N(0, 1)
//	σ²ₜ = ω + α · ε²ₜ₋₁ + β · σ²ₜ₋₁     (α ≥ 0, β ≥ 0, α + β < 1)
type GARCHResult struct {
	// Mu is the constant mean of returns (sample mean).
	Mu float64 `json:"mu"`
	// Omega is the variance constant.
	Omega float64 `json:"omega"`
	// Alpha is the ARCH coefficient on lagged squared residual.
	Alpha float64 `json:"alpha"`
	// Beta is the GARCH coefficient on lagged variance.
	Beta float64 `json:"beta"`
	// Persistence = α + β. Closer to 1 means shocks decay slowly.
	Persistence float64 `json:"persistence"`
	// LongRunVariance = ω / (1 − α − β). The unconditional variance.
	LongRunVariance float64 `json:"long_run_variance"`
	// LongRunStdDev = √LongRunVariance.
	LongRunStdDev float64 `json:"long_run_stddev"`
	// Variance is the conditional σ²ₜ series, length len(returns).
	Variance []float64 `json:"variance"`
	// StdDev is the conditional σₜ series, length len(returns).
	StdDev []float64 `json:"stddev"`
	// NextVariance is the one-step-ahead forecast for σ²_{T+1}.
	NextVariance float64 `json:"next_variance"`
	// NextStdDev is √NextVariance.
	NextStdDev float64 `json:"next_stddev"`
	// LogLikelihood at the optimised parameters.
	LogLikelihood float64 `json:"log_likelihood"`
	// Iterations Nelder-Mead steps actually taken.
	Iterations int `json:"iterations"`
}

// GARCH11 fits a GARCH(1,1) model to log-returns of the close-price series
// using variance targeting (μ = sample mean of returns, ω = V·(1 − α − β))
// and Nelder-Mead optimisation in (α, β). Defaults: maxIter = 200.
func GARCH11(close []float64, maxIter int) GARCHResult {
	out := GARCHResult{}
	if len(close) < 30 {
		return out
	}
	if maxIter <= 0 {
		maxIter = 200
	}
	// log-returns
	r := make([]float64, len(close)-1)
	for i := 1; i < len(close); i++ {
		if close[i-1] <= 0 || close[i] <= 0 {
			continue
		}
		r[i-1] = math.Log(close[i] / close[i-1])
	}
	n := len(r)

	// Mean and unconditional variance from the data.
	var mu float64
	for _, v := range r {
		mu += v
	}
	mu /= float64(n)

	var V float64
	for _, v := range r {
		d := v - mu
		V += d * d
	}
	V /= float64(n)
	if V < 1e-12 {
		V = 1e-12
	}

	// Negative log-likelihood with variance targeting.
	negLL := func(alpha, beta float64) float64 {
		// Reject infeasible region with a soft, monotonic penalty so the
		// optimiser can still descend toward feasibility from outside.
		if alpha < 0 || beta < 0 || alpha+beta >= 0.999 {
			pen := 0.0
			if alpha < 0 {
				pen += -alpha * 1e6
			}
			if beta < 0 {
				pen += -beta * 1e6
			}
			if alpha+beta >= 0.999 {
				pen += (alpha + beta - 0.999) * 1e6
			}
			return 1e9 + pen
		}
		omega := V * (1 - alpha - beta)
		// Seed conditional variance with the unconditional variance.
		sigma2 := V
		ll := 0.0
		eps := r[0] - mu
		for t := 0; t < n; t++ {
			if t > 0 {
				sigma2 = omega + alpha*eps*eps + beta*sigma2
				eps = r[t] - mu
				if sigma2 < 1e-12 {
					sigma2 = 1e-12
				}
			}
			ll += -0.5 * (math.Log(2*math.Pi*sigma2) + eps*eps/sigma2)
		}
		return -ll
	}

	// 2-D Nelder-Mead in (α, β).
	x0 := [3][2]float64{
		{0.05, 0.90},
		{0.10, 0.85},
		{0.05, 0.80},
	}
	f := [3]float64{
		negLL(x0[0][0], x0[0][1]),
		negLL(x0[1][0], x0[1][1]),
		negLL(x0[2][0], x0[2][1]),
	}

	const tol = 1e-7
	var iters int
	for iters = 0; iters < maxIter; iters++ {
		// Sort vertices: best (lowest f) first.
		for i := 0; i < 3; i++ {
			for j := i + 1; j < 3; j++ {
				if f[j] < f[i] {
					f[i], f[j] = f[j], f[i]
					x0[i], x0[j] = x0[j], x0[i]
				}
			}
		}
		if math.Abs(f[2]-f[0]) < tol {
			break
		}
		// Centroid of best two.
		cx := (x0[0][0] + x0[0][1] + x0[1][0] + x0[1][1]) // unused — just below
		_ = cx
		c := [2]float64{(x0[0][0] + x0[1][0]) / 2, (x0[0][1] + x0[1][1]) / 2}

		// Reflection.
		xr := [2]float64{2*c[0] - x0[2][0], 2*c[1] - x0[2][1]}
		fr := negLL(xr[0], xr[1])
		switch {
		case fr < f[0]:
			// Expand.
			xe := [2]float64{c[0] + 2*(c[0]-x0[2][0]), c[1] + 2*(c[1]-x0[2][1])}
			fe := negLL(xe[0], xe[1])
			if fe < fr {
				x0[2], f[2] = xe, fe
			} else {
				x0[2], f[2] = xr, fr
			}
		case fr < f[1]:
			x0[2], f[2] = xr, fr
		default:
			// Contraction.
			xc := [2]float64{c[0] + 0.5*(x0[2][0]-c[0]), c[1] + 0.5*(x0[2][1]-c[1])}
			fc := negLL(xc[0], xc[1])
			if fc < f[2] {
				x0[2], f[2] = xc, fc
			} else {
				// Shrink toward best.
				for i := 1; i < 3; i++ {
					x0[i][0] = x0[0][0] + 0.5*(x0[i][0]-x0[0][0])
					x0[i][1] = x0[0][1] + 0.5*(x0[i][1]-x0[0][1])
					f[i] = negLL(x0[i][0], x0[i][1])
				}
			}
		}
	}

	alpha := x0[0][0]
	beta := x0[0][1]
	if alpha < 0 {
		alpha = 0
	}
	if beta < 0 {
		beta = 0
	}
	if alpha+beta >= 0.999 {
		// Project back into the feasible region while preserving the ratio.
		s := alpha + beta
		alpha *= 0.998 / s
		beta *= 0.998 / s
	}
	omega := V * (1 - alpha - beta)
	if omega < 1e-12 {
		omega = 1e-12
	}

	// Replay to fill the conditional variance series and grab the final ε² for forecasting.
	variance := make([]float64, n)
	stddev := make([]float64, n)
	sigma2 := V
	eps := r[0] - mu
	variance[0] = sigma2
	stddev[0] = math.Sqrt(sigma2)
	ll := -0.5 * (math.Log(2*math.Pi*sigma2) + eps*eps/sigma2)
	for t := 1; t < n; t++ {
		sigma2 = omega + alpha*eps*eps + beta*sigma2
		if sigma2 < 1e-12 {
			sigma2 = 1e-12
		}
		eps = r[t] - mu
		variance[t] = sigma2
		stddev[t] = math.Sqrt(sigma2)
		ll += -0.5 * (math.Log(2*math.Pi*sigma2) + eps*eps/sigma2)
	}
	nextVar := omega + alpha*eps*eps + beta*sigma2
	if nextVar < 1e-12 {
		nextVar = 1e-12
	}

	out.Mu = mu
	out.Omega = omega
	out.Alpha = alpha
	out.Beta = beta
	out.Persistence = alpha + beta
	out.LongRunVariance = V
	out.LongRunStdDev = math.Sqrt(V)
	out.Variance = variance
	out.StdDev = stddev
	out.NextVariance = nextVar
	out.NextStdDev = math.Sqrt(nextVar)
	out.LogLikelihood = ll
	out.Iterations = iters
	return out
}

// GARCHSummary builds a one-line summary for a GARCHResult.
func GARCHSummary(r GARCHResult) string {
	if len(r.Variance) == 0 {
		return "garch: insufficient data"
	}
	return fmt.Sprintf(
		"garch(1,1): α=%.4f β=%.4f ω=%.2e persistence=%.4f σ_next=%.5f longrun_σ=%.5f, %d iters, ll=%.3f",
		r.Alpha, r.Beta, r.Omega, r.Persistence, r.NextStdDev, r.LongRunStdDev, r.Iterations, r.LogLikelihood,
	)
}
