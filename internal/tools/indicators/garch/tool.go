// Package garch registers the garch indicator with the talib dispatcher.
//
// The "garch" indicator fits a GARCH(1,1) model to log-returns of an input
// close-price series and returns the fitted parameters, the per-bar
// conditional variance/stddev series, and the one-step-ahead forecast.
package garch

import (
	"github.com/rangertaha/tamcp/internal/tools/indicators/talib"
)

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "garch",
		Description: "GARCH(1,1) volatility model: fits ω, α, β to log-returns by maximum likelihood with variance targeting and returns the conditional variance/stddev series plus the one-step-ahead forecast.",
		Group:       "volatility",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "close-price series"},
			{Name: "max_iter", Type: "int", Default: 200, Desc: "maximum Nelder-Mead iterations (default 200)"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			it := talib.ArgInt(args, "max_iter", 200)
			r := talib.GARCH11(v, it)
			return r, talib.GARCHSummary(r), nil
		},
	})
}
