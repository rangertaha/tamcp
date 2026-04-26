// Package ccisignal registers CCI with EMA signal (utility).
package ccisignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cci_signal",
		Description: "Commodity Channel Index with an EMA-smoothed signal line.",
		Group:       "momentum",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "period", Type: "int", Default: 20},
			{Name: "signal_period", Type: "int", Default: 9},
		},
		Run: func(args map[string]any) (any, string, error) {
			h, err := talib.ArgFloats(args, "high")
			if err != nil {
				return nil, "", err
			}
			l, err := talib.ArgFloats(args, "low")
			if err != nil {
				return nil, "", err
			}
			c, err := talib.ArgFloats(args, "close")
			if err != nil {
				return nil, "", err
			}
			ci, sig := talib.CCISIGNALFn(h, l, c,
				talib.ArgInt(args, "period", 20),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(ci, sig, [2]string{"cci", "signal"}), talib.Tersum("cci_signal", ci), nil
		},
	})
}
