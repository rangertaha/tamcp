// Package qqe registers the Quantitative Qualitative Estimation (Pandas TA, simplified).
package qqe

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "qqe",
		Description: "Quantitative Qualitative Estimation (simplified). Returns rsi_ma (EMA of RSI) and dar (dynamic ATR-of-RSI band).",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 14, Desc: "RSI period"},
			{Name: "smooth", Type: "int", Default: 5, Desc: "EMA smoothing of RSI"},
			{Name: "factor", Type: "float", Default: 4.236, Desc: "DAR multiplier"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			rm, dar := talib.QQEFn(v,
				talib.ArgInt(args, "period", 14),
				talib.ArgInt(args, "smooth", 5),
				talib.ArgFloat(args, "factor", 4.236),
			)
			return talib.Two(rm, dar, [2]string{"rsi_ma", "dar"}), talib.Tersum("qqe", rm), nil
		},
	})
}
