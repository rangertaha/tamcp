// Package bsts registers the bsts indicator with the talib dispatcher.
//
// "bsts" fits a Bayesian Structural Time Series (local linear trend) model
// to the input series via Kalman-filter ML estimation of the variance
// components. Returns the filtered level/slope, per-bar series, and a
// one-step-ahead forecast.
package bsts

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "bsts",
		Description: "Bayesian Structural Time Series (local linear trend): fits level + slope state-space model via Kalman-filter ML and returns per-bar level/slope plus one-step forecast.",
		Group:       "statistic",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "input series (e.g. close prices)"},
			{Name: "max_iter", Type: "int", Default: 200, Desc: "max Nelder-Mead iterations (default 200)"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			it := talib.ArgInt(args, "max_iter", 200)
			r := talib.BSTS(v, it)
			return r, talib.BSTSSummary(r), nil
		},
	})
}
