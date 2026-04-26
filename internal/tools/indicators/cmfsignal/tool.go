// Package cmfsignal registers CMF with EMA signal (utility).
package cmfsignal

import "github.com/rangertaha/tamcp/internal/tools/indicators/talib"

func init() {
	talib.RegisterEntry(&talib.Entry{
		Name:        "cmf_signal",
		Description: "Chaikin Money Flow with an EMA-smoothed signal line.",
		Group:       "volume",
		Params: []talib.Param{
			{Name: "high", Type: "number[]", Required: true, Desc: "High prices"},
			{Name: "low", Type: "number[]", Required: true, Desc: "Low prices"},
			{Name: "close", Type: "number[]", Required: true, Desc: "Close prices"},
			{Name: "volume", Type: "number[]", Required: true, Desc: "Volume"},
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
			vol, err := talib.ArgFloats(args, "volume")
			if err != nil {
				return nil, "", err
			}
			cm, sig := talib.CMFSIGNALFn(h, l, c, vol,
				talib.ArgInt(args, "period", 20),
				talib.ArgInt(args, "signal_period", 9),
			)
			return talib.Two(cm, sig, [2]string{"cmf", "signal"}), talib.Tersum("cmf_signal", cm), nil
		},
	})
}
