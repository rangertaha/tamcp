// Package crsi registers Connors RSI (Pandas TA).
package crsi

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "crsi",
		Description: "Connors RSI = (RSI(close, rsi_period) + RSI(streak, streak_period) + percent_rank(ROC, rank_period)) / 3.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "rsi_period", Type: "int", Default: 3},
			{Name: "streak_period", Type: "int", Default: 2},
			{Name: "rank_period", Type: "int", Default: 100},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			out := talib.CRSIFn(v,
				talib.ArgInt(args, "rsi_period", 3),
				talib.ArgInt(args, "streak_period", 2),
				talib.ArgInt(args, "rank_period", 100),
			)
			return talib.One(out), talib.Tersum("crsi", out), nil
		},
	})
}
