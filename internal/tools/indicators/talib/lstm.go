package talib

import (
	"fmt"
	"math"
	"math/rand"
)

// LSTMResult holds a trained single-cell LSTM and its one-step forecast.
//
// Input: a close-price series. The model trains on standardised log-returns
// and predicts the next return; that prediction is mapped back to a forecast
// close.
type LSTMResult struct {
	HiddenSize int     `json:"hidden_size"`
	Epochs     int     `json:"epochs"`
	FinalMSE   float64 `json:"final_mse"`
	LearnRate  float64 `json:"learn_rate"`
	// Predicted next log-return and price.
	ExpectedLogReturn float64 `json:"expected_log_return"`
	LastClose         float64 `json:"last_close"`
	NextClose         float64 `json:"next_close"`
	// Trained weights' L2 norm (sanity-check that training didn't diverge).
	WeightNorm float64 `json:"weight_norm"`
}

// LSTM trains a single-cell LSTM with the given hidden size on the log-returns
// of the close-price series for one-step-ahead prediction, then forecasts the
// next bar. Defaults: hiddenSize = 8, epochs = 50, learnRate = 0.05.
func LSTM(close []float64, hiddenSize, epochs int, learnRate float64) LSTMResult {
	out := LSTMResult{}
	if len(close) < 30 {
		return out
	}
	if hiddenSize <= 0 {
		hiddenSize = 8
	}
	if epochs <= 0 {
		epochs = 50
	}
	if learnRate <= 0 {
		learnRate = 0.05
	}
	out.HiddenSize = hiddenSize
	out.Epochs = epochs
	out.LearnRate = learnRate
	out.LastClose = close[len(close)-1]

	// Build returns and standardise.
	r := make([]float64, len(close)-1)
	for i := 1; i < len(close); i++ {
		if close[i-1] > 0 && close[i] > 0 {
			r[i-1] = math.Log(close[i] / close[i-1])
		}
	}
	mean, sd := meanSD(r)
	if sd <= 0 {
		sd = 1
	}
	z := make([]float64, len(r))
	for i, v := range r {
		z[i] = (v - mean) / sd
	}
	// Inputs xₜ = z[t], targets are z[t+1] for t in [0, n-2].
	n := len(z) - 1
	if n < 10 {
		return out
	}

	rng := rand.New(rand.NewSource(20260425)) // deterministic init
	H := hiddenSize

	// Concatenated [x; h] has size H+1; gates get an (H × (H+1)) weight matrix.
	wf := initMat(rng, H, H+1)
	wi := initMat(rng, H, H+1)
	wo := initMat(rng, H, H+1)
	wg := initMat(rng, H, H+1)
	bf := make([]float64, H)
	bi := make([]float64, H)
	bo := make([]float64, H)
	bg := make([]float64, H)
	wy := make([]float64, H) // output is scalar
	by := 0.0
	for j := 0; j < H; j++ {
		wy[j] = (rng.Float64()*2 - 1) * 0.1
	}

	// Forward + BPTT. Store gate activations per t for backward.
	type cache struct {
		x       float64
		hPrev   []float64
		cPrev   []float64
		concat  []float64
		f, i, o []float64
		g       []float64
		c       []float64
		h       []float64
		tanhC   []float64
		yhat    float64
	}

	matVec := func(W [][]float64, v []float64, b []float64) []float64 {
		out := make([]float64, len(W))
		for i := range W {
			s := b[i]
			for j := range v {
				s += W[i][j] * v[j]
			}
			out[i] = s
		}
		return out
	}
	sig := func(x float64) float64 { return 1 / (1 + math.Exp(-x)) }
	finalMSE := 0.0
	for ep := 0; ep < epochs; ep++ {
		caches := make([]cache, n)
		hPrev := make([]float64, H)
		cPrev := make([]float64, H)
		mse := 0.0
		// ── forward ─────────────────────────────────
		for t := 0; t < n; t++ {
			x := z[t]
			concat := append([]float64{x}, hPrev...)
			pf := matVec(wf, concat, bf)
			pi := matVec(wi, concat, bi)
			po := matVec(wo, concat, bo)
			pg := matVec(wg, concat, bg)
			f := make([]float64, H)
			i_ := make([]float64, H)
			o := make([]float64, H)
			g := make([]float64, H)
			cT := make([]float64, H)
			hT := make([]float64, H)
			tanhC := make([]float64, H)
			for j := 0; j < H; j++ {
				f[j] = sig(pf[j])
				i_[j] = sig(pi[j])
				o[j] = sig(po[j])
				g[j] = math.Tanh(pg[j])
				cT[j] = f[j]*cPrev[j] + i_[j]*g[j]
				tanhC[j] = math.Tanh(cT[j])
				hT[j] = o[j] * tanhC[j]
			}
			yhat := by
			for j := 0; j < H; j++ {
				yhat += wy[j] * hT[j]
			}
			caches[t] = cache{
				x: x, hPrev: hPrev, cPrev: cPrev,
				concat: concat, f: f, i: i_, o: o, g: g,
				c: cT, h: hT, tanhC: tanhC, yhat: yhat,
			}
			err := yhat - z[t+1]
			mse += err * err
			hPrev, cPrev = hT, cT
		}
		mse /= float64(n)
		finalMSE = mse

		// ── backward (BPTT) ────────────────────────
		dwf := zeroMat(H, H+1)
		dwi := zeroMat(H, H+1)
		dwo := zeroMat(H, H+1)
		dwg := zeroMat(H, H+1)
		dbf := make([]float64, H)
		dbi := make([]float64, H)
		dbo := make([]float64, H)
		dbg := make([]float64, H)
		dwy := make([]float64, H)
		dby := 0.0

		dhNext := make([]float64, H)
		dcNext := make([]float64, H)

		for t := n - 1; t >= 0; t-- {
			c := caches[t]
			err := c.yhat - z[t+1]
			// dL/dy
			dby += err
			for j := 0; j < H; j++ {
				dwy[j] += err * c.h[j]
			}
			dh := make([]float64, H)
			for j := 0; j < H; j++ {
				dh[j] = err*wy[j] + dhNext[j]
			}
			// h = o · tanh(c)
			do := make([]float64, H)
			dc := make([]float64, H)
			for j := 0; j < H; j++ {
				do[j] = dh[j] * c.tanhC[j]
				dc[j] = dh[j]*c.o[j]*(1-c.tanhC[j]*c.tanhC[j]) + dcNext[j]
			}
			// c = f · c_prev + i · g
			df := make([]float64, H)
			di := make([]float64, H)
			dg := make([]float64, H)
			for j := 0; j < H; j++ {
				df[j] = dc[j] * c.cPrev[j]
				di[j] = dc[j] * c.g[j]
				dg[j] = dc[j] * c.i[j]
				dcNext[j] = dc[j] * c.f[j]
			}
			// Pre-activation gradients.
			dpf := make([]float64, H)
			dpi := make([]float64, H)
			dpo := make([]float64, H)
			dpg := make([]float64, H)
			for j := 0; j < H; j++ {
				dpf[j] = df[j] * c.f[j] * (1 - c.f[j])
				dpi[j] = di[j] * c.i[j] * (1 - c.i[j])
				dpo[j] = do[j] * c.o[j] * (1 - c.o[j])
				dpg[j] = dg[j] * (1 - c.g[j]*c.g[j])
			}
			// Accumulate weight grads (outer product with c.concat).
			for j := 0; j < H; j++ {
				dbf[j] += dpf[j]
				dbi[j] += dpi[j]
				dbo[j] += dpo[j]
				dbg[j] += dpg[j]
				for k := 0; k < H+1; k++ {
					dwf[j][k] += dpf[j] * c.concat[k]
					dwi[j][k] += dpi[j] * c.concat[k]
					dwo[j][k] += dpo[j] * c.concat[k]
					dwg[j][k] += dpg[j] * c.concat[k]
				}
			}
			// dh_prev: contributions through the H entries of concat (indices 1..H).
			dhPrev := make([]float64, H)
			for k := 0; k < H; k++ {
				for j := 0; j < H; j++ {
					dhPrev[k] += dpf[j]*wf[j][k+1] + dpi[j]*wi[j][k+1] + dpo[j]*wo[j][k+1] + dpg[j]*wg[j][k+1]
				}
			}
			dhNext = dhPrev
		}

		// SGD step (no momentum to keep this compact). Gradient already summed; divide by n.
		nf := float64(n)
		for j := 0; j < H; j++ {
			by -= learnRate * dby / nf
			break
		}
		for j := 0; j < H; j++ {
			wy[j] -= learnRate * dwy[j] / nf
			bf[j] -= learnRate * dbf[j] / nf
			bi[j] -= learnRate * dbi[j] / nf
			bo[j] -= learnRate * dbo[j] / nf
			bg[j] -= learnRate * dbg[j] / nf
			for k := 0; k < H+1; k++ {
				wf[j][k] -= learnRate * dwf[j][k] / nf
				wi[j][k] -= learnRate * dwi[j][k] / nf
				wo[j][k] -= learnRate * dwo[j][k] / nf
				wg[j][k] -= learnRate * dwg[j][k] / nf
			}
		}
	}

	// One-step forecast: feed the entire sequence forward and read the final
	// hidden state's projection. The model was trained to predict z[t+1]
	// from z[t]; the last training example yielded z[n] from z[n-1], so to
	// get the prediction for "the bar after the last one" feed z[n-1] through
	// the cell once more starting from the state after the last training step.
	hPrev := make([]float64, H)
	cPrev := make([]float64, H)
	var yhat float64
	for t := 0; t < len(z); t++ {
		x := z[t]
		concat := append([]float64{x}, hPrev...)
		pf := matVec(wf, concat, bf)
		pi := matVec(wi, concat, bi)
		po := matVec(wo, concat, bo)
		pg := matVec(wg, concat, bg)
		f := make([]float64, H)
		i_ := make([]float64, H)
		o := make([]float64, H)
		g := make([]float64, H)
		cT := make([]float64, H)
		hT := make([]float64, H)
		yhat = by
		for j := 0; j < H; j++ {
			f[j] = sig(pf[j])
			i_[j] = sig(pi[j])
			o[j] = sig(po[j])
			g[j] = math.Tanh(pg[j])
			cT[j] = f[j]*cPrev[j] + i_[j]*g[j]
			hT[j] = o[j] * math.Tanh(cT[j])
			yhat += wy[j] * hT[j]
		}
		hPrev, cPrev = hT, cT
	}
	// yhat is in standardised space; un-standardise back to a log-return.
	expRet := yhat*sd + mean

	wn := 0.0
	for _, row := range wf {
		for _, v := range row {
			wn += v * v
		}
	}
	for _, row := range wi {
		for _, v := range row {
			wn += v * v
		}
	}
	for _, row := range wo {
		for _, v := range row {
			wn += v * v
		}
	}
	for _, row := range wg {
		for _, v := range row {
			wn += v * v
		}
	}
	for _, v := range wy {
		wn += v * v
	}

	out.FinalMSE = finalMSE
	out.ExpectedLogReturn = expRet
	out.NextClose = out.LastClose * math.Exp(expRet)
	out.WeightNorm = math.Sqrt(wn)
	return out
}

// LSTMSummary builds a one-line summary for an LSTMResult.
func LSTMSummary(r LSTMResult) string {
	if r.HiddenSize == 0 {
		return "lstm: insufficient data"
	}
	return fmt.Sprintf(
		"lstm: H=%d epochs=%d MSE=%.6f next≈%.6f (E[r]=%.5f), |W|=%.3f",
		r.HiddenSize, r.Epochs, r.FinalMSE, r.NextClose, r.ExpectedLogReturn, r.WeightNorm,
	)
}

func initMat(rng *rand.Rand, rows, cols int) [][]float64 {
	scale := math.Sqrt(1.0 / float64(cols))
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			m[i][j] = (rng.Float64()*2 - 1) * scale
		}
	}
	return m
}

func zeroMat(rows, cols int) [][]float64 {
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
	}
	return m
}

func meanSD(x []float64) (float64, float64) {
	if len(x) == 0 {
		return 0, 0
	}
	var s float64
	for _, v := range x {
		s += v
	}
	m := s / float64(len(x))
	var v float64
	for _, x := range x {
		v += (x - m) * (x - m)
	}
	v /= float64(len(x))
	return m, math.Sqrt(v)
}
