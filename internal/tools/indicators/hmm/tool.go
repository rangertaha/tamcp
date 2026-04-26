// Package hmm registers the hmm indicator with the talib dispatcher.
//
// The "hmm" indicator fits a K-state Gaussian Hidden Markov Model to the
// log-returns of an input close-price series via Baum-Welch (EM), then
// returns a one-step-ahead forecast plus the trained model parameters.
package hmm

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "hmm",
		Description: "Gaussian HMM one-step forecast: fits a K-state HMM to log-returns and predicts the next bar's expected return, expected close, and Gaussian-mixture standard deviation.",
		Group:       "statistic",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "close-price series"},
			{Name: "num_states", Type: "int", Default: 2, Desc: "number of hidden states K (default 2)"},
			{Name: "max_iter", Type: "int", Default: 15, Desc: "maximum EM iterations (default 15)"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			K := talib.ArgInt(args, "num_states", 2)
			it := talib.ArgInt(args, "max_iter", 15)
			r := talib.HMMForecast(v, K, it)
			return r, talib.HMMSummary(r), nil
		},
	})
}
