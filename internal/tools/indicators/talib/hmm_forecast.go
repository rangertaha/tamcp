package talib

import (
	"fmt"
	"math"
	"sort"
)

// HMMForecastResult is the one-step-ahead forecast and trained parameters
// for a K-state Gaussian HMM fit to log-returns of a price series.
type HMMForecastResult struct {
	// Means / Variances of the K Gaussian emission components.
	Means     []float64 `json:"means"`
	Variances []float64 `json:"variances"`
	// Transition[i*K+j] = P(state_t=j | state_{t-1}=i).
	Transition []float64 `json:"transition"`
	// Posteriors over the state at the last observed bar.
	Posteriors []float64 `json:"posteriors"`
	// NextStateProb[j] = Σᵢ Posteriors[i] · Transition[i,j].
	NextStateProb []float64 `json:"next_state_prob"`
	// ExpectedLogReturn = Σⱼ NextStateProb[j] · Means[j].
	ExpectedLogReturn float64 `json:"expected_log_return"`
	// ExpectedReturnSD is the Gaussian-mixture standard deviation of the next return.
	ExpectedReturnSD float64 `json:"expected_return_sd"`
	// LastClose / NextClose: last observed close and the model's expected close.
	LastClose float64 `json:"last_close"`
	NextClose float64 `json:"next_close"`
	// LogLikelihood at the final EM iteration.
	LogLikelihood float64 `json:"log_likelihood"`
	// Iterations actually run before the EM converged.
	Iterations int `json:"iterations"`
}

// HMMForecast fits a K-state Gaussian-emission HMM to the log-returns of
// the input close series and returns a one-step-ahead forecast plus the
// trained model parameters. Defaults: K=2, maxIter=15, tol=1e-3.
func HMMForecast(close []float64, K, maxIter int) HMMForecastResult {
	if K < 1 {
		K = 2
	}
	if maxIter <= 0 {
		maxIter = 15
	}
	const tol = 1e-3

	out := HMMForecastResult{
		Means:         make([]float64, K),
		Variances:     make([]float64, K),
		Transition:    make([]float64, K*K),
		Posteriors:    make([]float64, K),
		NextStateProb: make([]float64, K),
	}
	if len(close) < 30 {
		return out
	}
	out.LastClose = close[len(close)-1]

	// observations = log-returns
	obs := make([]float64, len(close)-1)
	for i := 1; i < len(close); i++ {
		if close[i-1] <= 0 || close[i] <= 0 {
			obs[i-1] = 0
			continue
		}
		obs[i-1] = math.Log(close[i] / close[i-1])
	}
	n := len(obs)

	// Initialise μ from sorted quantile centroids; σ² shared as global variance.
	sob := append([]float64(nil), obs...)
	sort.Float64s(sob)
	mean := 0.0
	for _, v := range obs {
		mean += v
	}
	mean /= float64(n)
	v0 := 0.0
	for _, v := range obs {
		v0 += (v - mean) * (v - mean)
	}
	v0 /= float64(n)
	if v0 < 1e-12 {
		v0 = 1e-12
	}
	mu := make([]float64, K)
	sigma2 := make([]float64, K)
	for k := 0; k < K; k++ {
		lo := k * n / K
		hi := (k + 1) * n / K
		s := 0.0
		for j := lo; j < hi && j < n; j++ {
			s += sob[j]
		}
		if hi > lo {
			mu[k] = s / float64(hi-lo)
		} else {
			mu[k] = sob[lo]
		}
		sigma2[k] = v0
	}

	// Initial state and transition: uniform π, sticky A.
	pi := make([]float64, K)
	A := make([]float64, K*K)
	for i := 0; i < K; i++ {
		pi[i] = 1.0 / float64(K)
		for j := 0; j < K; j++ {
			if i == j {
				A[i*K+j] = 0.9
			} else {
				A[i*K+j] = 0.1 / float64(K-1)
			}
		}
	}

	logA := make([]float64, n*K)
	logB := make([]float64, n*K)
	gamma := make([]float64, n*K)
	xi := make([]float64, K*K)
	tmp := make([]float64, K)

	gaussLog := func(x, m, v float64) float64 {
		if v <= 1e-12 {
			v = 1e-12
		}
		d := x - m
		return -0.5*math.Log(2*math.Pi*v) - (d*d)/(2*v)
	}
	logSafe := func(x float64) float64 {
		if x <= 0 {
			return -1e300
		}
		return math.Log(x)
	}

	prevLL := 0.0
	hadLL := false
	var iter int
	for iter = 0; iter < maxIter; iter++ {
		// Forward.
		for k := 0; k < K; k++ {
			logA[k] = logSafe(pi[k]) + gaussLog(obs[0], mu[k], sigma2[k])
		}
		for t := 1; t < n; t++ {
			for j := 0; j < K; j++ {
				maxv := math.Inf(-1)
				for i := 0; i < K; i++ {
					tmp[i] = logA[(t-1)*K+i] + logSafe(A[i*K+j])
					if tmp[i] > maxv {
						maxv = tmp[i]
					}
				}
				s := 0.0
				for i := 0; i < K; i++ {
					s += math.Exp(tmp[i] - maxv)
				}
				logA[t*K+j] = maxv + math.Log(s) + gaussLog(obs[t], mu[j], sigma2[j])
			}
		}
		// Log-likelihood.
		maxLL := math.Inf(-1)
		for k := 0; k < K; k++ {
			if logA[(n-1)*K+k] > maxLL {
				maxLL = logA[(n-1)*K+k]
			}
		}
		s := 0.0
		for k := 0; k < K; k++ {
			s += math.Exp(logA[(n-1)*K+k] - maxLL)
		}
		ll := maxLL + math.Log(s)

		// Backward.
		for k := 0; k < K; k++ {
			logB[(n-1)*K+k] = 0
		}
		for t := n - 2; t >= 0; t-- {
			for i := 0; i < K; i++ {
				maxv := math.Inf(-1)
				for j := 0; j < K; j++ {
					tmp[j] = logSafe(A[i*K+j]) + gaussLog(obs[t+1], mu[j], sigma2[j]) + logB[(t+1)*K+j]
					if tmp[j] > maxv {
						maxv = tmp[j]
					}
				}
				s2 := 0.0
				for j := 0; j < K; j++ {
					s2 += math.Exp(tmp[j] - maxv)
				}
				logB[t*K+i] = maxv + math.Log(s2)
			}
		}

		// γ = exp(α + β − ll); ξ accumulated for transition.
		for t := 0; t < n; t++ {
			for k := 0; k < K; k++ {
				gamma[t*K+k] = math.Exp(logA[t*K+k] + logB[t*K+k] - ll)
			}
		}
		for i := range xi {
			xi[i] = 0
		}
		for t := 0; t < n-1; t++ {
			denom := 0.0
			for i := 0; i < K; i++ {
				for j := 0; j < K; j++ {
					term := math.Exp(
						logA[t*K+i] + logSafe(A[i*K+j]) +
							gaussLog(obs[t+1], mu[j], sigma2[j]) +
							logB[(t+1)*K+j] - ll)
					xi[i*K+j] += term
					denom += term
				}
			}
			_ = denom
		}

		// M-step.
		for k := 0; k < K; k++ {
			pi[k] = gamma[k]
		}
		normalize(pi)

		for i := 0; i < K; i++ {
			rowSum := 0.0
			for j := 0; j < K; j++ {
				rowSum += xi[i*K+j]
			}
			if rowSum > 0 {
				for j := 0; j < K; j++ {
					A[i*K+j] = xi[i*K+j] / rowSum
				}
			} else {
				for j := 0; j < K; j++ {
					A[i*K+j] = 1.0 / float64(K)
				}
			}
		}
		for k := 0; k < K; k++ {
			w, sm, sm2 := 0.0, 0.0, 0.0
			for t := 0; t < n; t++ {
				g := gamma[t*K+k]
				w += g
				sm += g * obs[t]
				sm2 += g * obs[t] * obs[t]
			}
			if w > 0 {
				m := sm / w
				v := sm2/w - m*m
				if v < 1e-12 {
					v = 1e-12
				}
				mu[k] = m
				sigma2[k] = v
			}
		}

		if hadLL && math.Abs(ll-prevLL) < tol {
			iter++
			prevLL = ll
			break
		}
		prevLL = ll
		hadLL = true
	}

	// Posteriors over state at t=n-1 from the final α.
	maxA := math.Inf(-1)
	for k := 0; k < K; k++ {
		if logA[(n-1)*K+k] > maxA {
			maxA = logA[(n-1)*K+k]
		}
	}
	post := make([]float64, K)
	totalP := 0.0
	for k := 0; k < K; k++ {
		post[k] = math.Exp(logA[(n-1)*K+k] - maxA)
		totalP += post[k]
	}
	if totalP > 0 {
		for k := 0; k < K; k++ {
			post[k] /= totalP
		}
	}
	// Next-state distribution.
	nextP := make([]float64, K)
	for j := 0; j < K; j++ {
		s := 0.0
		for i := 0; i < K; i++ {
			s += post[i] * A[i*K+j]
		}
		nextP[j] = s
	}
	// Mixture mean / variance for the next return.
	expRet := 0.0
	for k := 0; k < K; k++ {
		expRet += nextP[k] * mu[k]
	}
	expVar := 0.0
	for k := 0; k < K; k++ {
		expVar += nextP[k] * (sigma2[k] + (mu[k]-expRet)*(mu[k]-expRet))
	}

	out.Means = mu
	out.Variances = sigma2
	out.Transition = A
	out.Posteriors = post
	out.NextStateProb = nextP
	out.ExpectedLogReturn = expRet
	out.ExpectedReturnSD = math.Sqrt(math.Max(0, expVar))
	out.NextClose = out.LastClose * math.Exp(expRet)
	out.LogLikelihood = prevLL
	out.Iterations = iter
	return out
}

// HMMSummary builds a one-line text summary for an HMMForecastResult.
func HMMSummary(r HMMForecastResult) string {
	if len(r.Means) == 0 {
		return "hmm: insufficient data"
	}
	return fmt.Sprintf(
		"hmm: %d states, last=%.6f, next≈%.6f (E[r]=%.5f, σ=%.5f), %d iters, ll=%.3f",
		len(r.Means), r.LastClose, r.NextClose, r.ExpectedLogReturn, r.ExpectedReturnSD, r.Iterations, r.LogLikelihood,
	)
}

func normalize(v []float64) {
	s := 0.0
	for _, x := range v {
		s += x
	}
	if s == 0 {
		for i := range v {
			v[i] = 1.0 / float64(len(v))
		}
		return
	}
	for i := range v {
		v[i] /= s
	}
}
