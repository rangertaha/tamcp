// Package lstm registers the lstm indicator with the talib dispatcher.
//
// "lstm" trains a single-cell LSTM on the standardised log-returns of a
// close-price series and forecasts the next bar's close. Inputs match the
// shape of "hmm" / "garch" (values + a couple of training knobs).
package lstm

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "lstm",
		Description: "Single-cell LSTM trained via BPTT on log-returns; predicts the next bar's expected log-return and close.",
		Group:       "statistic",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "close-price series"},
			{Name: "hidden_size", Type: "int", Default: 8, Desc: "LSTM hidden width (default 8)"},
			{Name: "epochs", Type: "int", Default: 50, Desc: "training epochs (default 50)"},
			{Name: "learn_rate", Type: "number", Default: 0.05, Desc: "SGD learning rate (default 0.05)"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			h := talib.ArgInt(args, "hidden_size", 8)
			ep := talib.ArgInt(args, "epochs", 50)
			lr := talib.ArgFloat(args, "learn_rate", 0.05)
			r := talib.LSTM(v, h, ep, lr)
			return r, talib.LSTMSummary(r), nil
		},
	})
}
