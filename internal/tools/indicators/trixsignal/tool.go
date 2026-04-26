// Package trixsignal registers TRIX with an SMA signal line (Pandas TA).
package trixsignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "trix_signal",
		Description: "TRIX with an SMA signal line. Returns trix and signal.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "values", Type: "number[]", Required: true, Desc: "Price series"},
			{Name: "period", Type: "int", Default: 30, Desc: "TRIX EMA period"},
			{Name: "signal_period", Type: "int", Default: 9, Desc: "SMA period for signal line"},
		},
		Run: func(args map[string]any) (any, string, error) {
			v, err := talib.ArgFloats(args, "values")
			if err != nil {
				return nil, "", err
			}
			t, sig := talib.TRIXSIGNALFn(v,
				talib.ArgInt(args, "period", 30),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(t, sig, [2]string{"trix", "signal"}), talib.Tersum("trix_signal", t), nil
		},
	})
}
