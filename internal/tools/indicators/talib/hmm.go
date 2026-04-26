//go:build ignore

package talib

import (
	"math"
	"sort"
)

// HMMResult holds the trained parameters and Viterbi-decoded state sequence
// for a Gaussian Hidden Markov Model.
type HMMResult struct {
	// States is the most-likely hidden state for each observation (Viterbi).
	States []int `json:"states"`
	// Means[k] is the emission mean for state k.
	Means []float64 `json:"means"`
	// Variances[k] is the emission variance for state k.
	Variances []float64 `json:"variances"`
	// Initial[k] is the initial state probability.
	Initial []float64 `json:"initial"`
	// Transition[k*K+l] is P(state_t=l | state_{t-1}=k).
	Transition []float64 `json:"transition"`
	// LogLikelihood of the observations under the trained model.
	LogLikelihood float64 `json:"log_likelihood"`
	// Iterations actually run before convergence (or maxIter cap).
	Iterations int `json:"iterations"`
}

// HMM fits a K-state Gaussian-emission HMM to the observation series using the
// Baum–Welch (EM) algorithm, then decodes the most likely hidden state sequence
// via Viterbi.
//
// Parameters:
//   - obs:      observation series (e.g. log-returns for regime detection).
//   - numStates: number of hidden states K (>= 1).
//   - maxIter:  maximum EM iterations (default 50 if <= 0).
//   - tol:      log-likelihood improvement threshold for early stop (default 1e-4).
//
// Initialization sorts observations and splits into K equal quantile buckets to
// seed the means; variances are seeded as the global variance; transitions and
// initial probabilities are uniform.
func HMM(obs []float64, numStates, maxIter int, tol float64) HMMResult {
	n := len(obs)
	K := numStates
	if K < 1 {
		K = 1
	}
	if maxIter <= 0 {
		maxIter = 50
	}
	if tol <= 0 {
		tol = 1e-4
	}
	if n == 0 {
		return HMMResult{Means: make([]float64, K), Variances: make([]float64, K), Initial: uniform(K), Transition: uniformTransition(K)}
	}

	mu, sigma2 := hmmInit(obs, K)
	pi := uniform(K)
	A := uniformTransition(K)

	logAlpha := make([]float64, n*K)
	logBeta := make([]float64, n*K)
	gamma := make([]float64, n*K)
	xi := make([]float64, K*K) // accumulated over t for transition update

	var prevLL float64
	var iter int
	for iter = 0; iter < maxIter; iter++ {
		ll := hmmForward(obs, pi, A, mu, sigma2, logAlpha)
		hmmBackward(obs, A, mu, sigma2, logBeta)

		// gamma[t,k] = exp(logAlpha[t,k] + logBeta[t,k] - ll)
		for t := 0; t < n; t++ {
			for k := 0; k < K; k++ {
				gamma[t*K+k] = math.Exp(logAlpha[t*K+k] + logBeta[t*K+k] - ll)
			}
		}

		// xi[i*K+j] = sum_t exp(logAlpha[t,i] + log A[i,j] + log B(o_{t+1};j) + logBeta[t+1,j] - ll)
		for i := range xi {
			xi[i] = 0
		}
		for t := 0; t < n-1; t++ {
			for i := 0; i < K; i++ {
				for j := 0; j < K; j++ {
					logA := math.Log(safe(A[i*K+j]))
					logB := gaussLog(obs[t+1], mu[j], sigma2[j])
					xi[i*K+j] += math.Exp(logAlpha[t*K+i] + logA + logB + logBeta[(t+1)*K+j] - ll)
				}
			}
		}

		// M-step.
		// pi_k = gamma[0,k]
		for k := 0; k < K; k++ {
			pi[k] = gamma[k]
		}
		normalize(pi)

		// A[i,j] = xi[i,j] / sum_t<n-1 gamma[t,i]
		for i := 0; i < K; i++ {
			var denom float64
			for t := 0; t < n-1; t++ {
				denom += gamma[t*K+i]
			}
			if denom == 0 {
				for j := 0; j < K; j++ {
					A[i*K+j] = 1.0 / float64(K)
				}
				continue
			}
			for j := 0; j < K; j++ {
				A[i*K+j] = xi[i*K+j] / denom
			}
			normalizeRow(A, i, K)
		}

		// mu_k, sigma2_k from gamma-weighted observations.
		for k := 0; k < K; k++ {
			var w, sum, sum2 float64
			for t := 0; t < n; t++ {
				g := gamma[t*K+k]
				w += g
				sum += g * obs[t]
				sum2 += g * obs[t] * obs[t]
			}
			if w > 0 {
				m := sum / w
				v := sum2/w - m*m
				if v < 1e-12 {
					v = 1e-12
				}
				mu[k] = m
				sigma2[k] = v
			}
		}

		if iter > 0 && math.Abs(ll-prevLL) < tol {
			iter++
			prevLL = ll
			break
		}
		prevLL = ll
	}

	states := hmmViterbi(obs, pi, A, mu, sigma2)

	return HMMResult{
		States:        states,
		Means:         mu,
		Variances:     sigma2,
		Initial:       pi,
		Transition:    A,
		LogLikelihood: prevLL,
		Iterations:    iter,
	}
}

// hmmInit seeds means by quantile splits over a sorted copy of obs; variances
// are seeded with the global variance.
func hmmInit(obs []float64, K int) (mu, sigma2 []float64) {
	n := len(obs)
	mu = make([]float64, K)
	sigma2 = make([]float64, K)
	sorted := make([]float64, n)
	copy(sorted, obs)
	sort.Float64s(sorted)
	for k := 0; k < K; k++ {
		lo := k * n / K
		hi := (k + 1) * n / K
		if hi <= lo {
			mu[k] = sorted[lo]
			continue
		}
		var s float64
		for j := lo; j < hi; j++ {
			s += sorted[j]
		}
		mu[k] = s / float64(hi-lo)
	}
	// global variance
	var mean float64
	for _, v := range obs {
		mean += v
	}
	mean /= float64(n)
	var v float64
	for _, x := range obs {
		v += (x - mean) * (x - mean)
	}
	v /= float64(n)
	if v < 1e-12 {
		v = 1e-12
	}
	for k := 0; k < K; k++ {
		sigma2[k] = v
	}
	return
}

func uniform(K int) []float64 {
	out := make([]float64, K)
	for i := range out {
		out[i] = 1.0 / float64(K)
	}
	return out
}

func uniformTransition(K int) []float64 {
	out := make([]float64, K*K)
	for i := range out {
		out[i] = 1.0 / float64(K)
	}
	return out
}

func gaussLog(x, mu, sigma2 float64) float64 {
	if sigma2 <= 0 {
		sigma2 = 1e-12
	}
	d := x - mu
	return -0.5*math.Log(2*math.Pi*sigma2) - (d*d)/(2*sigma2)
}

// hmmForward runs the log-space forward pass; returns the log-likelihood of obs.
func hmmForward(obs, pi, A, mu, sigma2 []float64, logAlpha []float64) float64 {
	n := len(obs)
	K := len(pi)
	for k := 0; k < K; k++ {
		logAlpha[k] = math.Log(safe(pi[k])) + gaussLog(obs[0], mu[k], sigma2[k])
	}
	for t := 1; t < n; t++ {
		for j := 0; j < K; j++ {
			// logsumexp over i of logAlpha[t-1,i] + log A[i,j]
			max := math.Inf(-1)
			tmp := make([]float64, K)
			for i := 0; i < K; i++ {
				tmp[i] = logAlpha[(t-1)*K+i] + math.Log(safe(A[i*K+j]))
				if tmp[i] > max {
					max = tmp[i]
				}
			}
			var sum float64
			for _, v := range tmp {
				sum += math.Exp(v - max)
			}
			logAlpha[t*K+j] = max + math.Log(sum) + gaussLog(obs[t], mu[j], sigma2[j])
		}
	}
	// log P(obs) = logsumexp_k logAlpha[n-1,k]
	max := math.Inf(-1)
	for k := 0; k < K; k++ {
		if logAlpha[(n-1)*K+k] > max {
			max = logAlpha[(n-1)*K+k]
		}
	}
	var sum float64
	for k := 0; k < K; k++ {
		sum += math.Exp(logAlpha[(n-1)*K+k] - max)
	}
	return max + math.Log(sum)
}

// hmmBackward runs the log-space backward pass.
func hmmBackward(obs, A, mu, sigma2 []float64, logBeta []float64) {
	n := len(obs)
	K := len(mu)
	for k := 0; k < K; k++ {
		logBeta[(n-1)*K+k] = 0
	}
	for t := n - 2; t >= 0; t-- {
		for i := 0; i < K; i++ {
			max := math.Inf(-1)
			tmp := make([]float64, K)
			for j := 0; j < K; j++ {
				tmp[j] = math.Log(safe(A[i*K+j])) + gaussLog(obs[t+1], mu[j], sigma2[j]) + logBeta[(t+1)*K+j]
				if tmp[j] > max {
					max = tmp[j]
				}
			}
			var sum float64
			for _, v := range tmp {
				sum += math.Exp(v - max)
			}
			logBeta[t*K+i] = max + math.Log(sum)
		}
	}
}

// hmmViterbi decodes the most likely state sequence in log space.
func hmmViterbi(obs, pi, A, mu, sigma2 []float64) []int {
	n := len(obs)
	K := len(pi)
	delta := make([]float64, n*K)
	psi := make([]int, n*K)
	for k := 0; k < K; k++ {
		delta[k] = math.Log(safe(pi[k])) + gaussLog(obs[0], mu[k], sigma2[k])
	}
	for t := 1; t < n; t++ {
		for j := 0; j < K; j++ {
			best := math.Inf(-1)
			bestI := 0
			for i := 0; i < K; i++ {
				v := delta[(t-1)*K+i] + math.Log(safe(A[i*K+j]))
				if v > best {
					best = v
					bestI = i
				}
			}
			delta[t*K+j] = best + gaussLog(obs[t], mu[j], sigma2[j])
			psi[t*K+j] = bestI
		}
	}
	states := make([]int, n)
	best := math.Inf(-1)
	for k := 0; k < K; k++ {
		if delta[(n-1)*K+k] > best {
			best = delta[(n-1)*K+k]
			states[n-1] = k
		}
	}
	for t := n - 2; t >= 0; t-- {
		states[t] = psi[(t+1)*K+states[t+1]]
	}
	return states
}

func safe(x float64) float64 {
	if x <= 0 {
		return 1e-300
	}
	return x
}

func normalize(v []float64) {
	var s float64
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

func normalizeRow(A []float64, row, K int) {
	var s float64
	for j := 0; j < K; j++ {
		s += A[row*K+j]
	}
	if s == 0 {
		for j := 0; j < K; j++ {
			A[row*K+j] = 1.0 / float64(K)
		}
		return
	}
	for j := 0; j < K; j++ {
		A[row*K+j] /= s
	}
}
